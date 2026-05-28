package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	_ "github.com/tursodatabase/go-libsql"
)

const dbFileName = "gk.db"

type Storage struct {
	db *sql.DB
}

func NewStorage() (*Storage, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(home, ".gk")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("libsql", "file:"+filepath.Join(dir, dbFileName))
	if err != nil {
		return nil, err
	}

	storage := &Storage{db: db}
	if err := storage.init(); err != nil {
		db.Close()
		return nil, err
	}

	return storage, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) init() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			status INTEGER NOT NULL,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			position INTEGER NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS habits (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			position INTEGER NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS habit_checks (
			habit_id INTEGER NOT NULL,
			week_start TEXT NOT NULL,
			mon INTEGER NOT NULL DEFAULT 0,
			tue INTEGER NOT NULL DEFAULT 0,
			wed INTEGER NOT NULL DEFAULT 0,
			thu INTEGER NOT NULL DEFAULT 0,
			fri INTEGER NOT NULL DEFAULT 0,
			sat INTEGER NOT NULL DEFAULT 0,
			sun INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (habit_id, week_start)
		)
	`)
	return err
}

func (s *Storage) Load() (map[status][]list.Item, error) {
	rows, err := s.db.Query(`
		SELECT status, title, description
		FROM tasks
		ORDER BY status, position, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := map[status][]list.Item{
		todo:       {},
		inProgress: {},
		done:       {},
	}

	for rows.Next() {
		var stat status
		var title string
		var description string
		if err := rows.Scan(&stat, &title, &description); err != nil {
			return nil, err
		}
		tasks[stat] = append(tasks[stat], Task{status: stat, title: title, description: description})
	}

	return tasks, rows.Err()
}

func (s *Storage) Save(cols []column) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM tasks`); err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO tasks (status, title, description, position)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, col := range cols {
		for position, item := range col.list.Items() {
			task, ok := item.(Task)
			if !ok {
				return fmt.Errorf("unexpected list item type %T", item)
			}
			if _, err := stmt.Exec(col.status, task.title, task.description, position); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (s *Storage) LoadHabits(weekStart time.Time) ([]Habit, error) {
	return s.LoadHabitsWithStreakDate(weekStart, time.Now())
}

func (s *Storage) LoadHabitsWithStreakDate(weekStart, streakDate time.Time) ([]Habit, error) {
	rows, err := s.db.Query(`
		SELECT h.id, h.name,
			COALESCE(c.mon, 0), COALESCE(c.tue, 0), COALESCE(c.wed, 0), COALESCE(c.thu, 0),
			COALESCE(c.fri, 0), COALESCE(c.sat, 0), COALESCE(c.sun, 0)
		FROM habits h
		LEFT JOIN habit_checks c ON c.habit_id = h.id AND c.week_start = ?
		ORDER BY h.position, h.id
	`, weekKey(weekStart))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	habits := []Habit{}
	for rows.Next() {
		var id int64
		var name string
		var checked [7]bool
		var values [7]int
		if err := rows.Scan(&id, &name, &values[0], &values[1], &values[2], &values[3], &values[4], &values[5], &values[6]); err != nil {
			return nil, err
		}
		for i, value := range values {
			checked[i] = value != 0
		}
		habit := NewHabitWithChecks(name, checked)
		habit.id = id
		habit.streak, err = s.LoadHabitStreak(id, streakDate)
		if err != nil {
			return nil, err
		}
		habits = append(habits, habit)
	}

	return habits, rows.Err()
}

func (s *Storage) LoadHabitStreak(habitID int64, weekStart time.Time) (int, error) {
	week := startOfWeek(weekStart)
	day := weekdayIndex(weekStart)
	streak := 0

	for {
		checked, err := s.loadHabitChecks(habitID, week)
		if err != nil {
			return 0, err
		}

		for i := day; i >= 0; i-- {
			if !checked[i] {
				return streak, nil
			}
			streak++
		}

		week = week.AddDate(0, 0, -7)
		day = 6
	}
}

func (s *Storage) loadHabitChecks(habitID int64, weekStart time.Time) ([7]bool, error) {
	var values [7]int
	err := s.db.QueryRow(`
		SELECT mon, tue, wed, thu, fri, sat, sun
		FROM habit_checks
		WHERE habit_id = ? AND week_start = ?
	`, habitID, weekKey(weekStart)).Scan(&values[0], &values[1], &values[2], &values[3], &values[4], &values[5], &values[6])
	if err == sql.ErrNoRows {
		return [7]bool{}, nil
	}
	if err != nil {
		return [7]bool{}, err
	}

	var checked [7]bool
	for i, value := range values {
		checked[i] = value != 0
	}
	return checked, nil
}

func (s *Storage) SaveHabits(habits []Habit, weekStart time.Time) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM habit_checks WHERE week_start = ?`, weekKey(weekStart)); err != nil {
		return err
	}
	ids := make([]int64, 0, len(habits))
	for _, habit := range habits {
		if habit.id != 0 {
			ids = append(ids, habit.id)
		}
	}
	if err := deleteRemovedHabits(tx, ids); err != nil {
		return err
	}

	habitStmt, err := tx.Prepare(`
		INSERT INTO habits (id, name, position)
		VALUES (?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET name = excluded.name, position = excluded.position
	`)
	if err != nil {
		return err
	}
	defer habitStmt.Close()

	checkStmt, err := tx.Prepare(`
		INSERT INTO habit_checks (habit_id, week_start, mon, tue, wed, thu, fri, sat, sun)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer checkStmt.Close()

	for position, habit := range habits {
		id := habit.id
		if id == 0 {
			res, err := tx.Exec(`
				INSERT INTO habits (name, position)
				VALUES (?, ?)
			`, habit.name, position)
			if err != nil {
				return err
			}
			id, err = res.LastInsertId()
			if err != nil {
				return err
			}
		} else {
			if _, err := habitStmt.Exec(id, habit.name, position); err != nil {
				return err
			}
		}
		if _, err := checkStmt.Exec(id, weekKey(weekStart), boolInt(habit.checked[0]), boolInt(habit.checked[1]), boolInt(habit.checked[2]), boolInt(habit.checked[3]), boolInt(habit.checked[4]), boolInt(habit.checked[5]), boolInt(habit.checked[6])); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func isDuplicateColumnError(err error) bool {
	return strings.Contains(err.Error(), "duplicate column name")
}

func weekKey(t time.Time) string {
	return startOfWeek(t).Format("2006-01-02")
}

func deleteRemovedHabits(tx *sql.Tx, ids []int64) error {
	if len(ids) == 0 {
		if _, err := tx.Exec(`DELETE FROM habit_checks`); err != nil {
			return err
		}
		_, err := tx.Exec(`DELETE FROM habits`)
		return err
	}

	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	query := fmt.Sprintf(`DELETE FROM habit_checks WHERE habit_id NOT IN (%s)`, strings.Join(placeholders, ","))
	if _, err := tx.Exec(query, args...); err != nil {
		return err
	}
	query = fmt.Sprintf(`DELETE FROM habits WHERE id NOT IN (%s)`, strings.Join(placeholders, ","))
	_, err := tx.Exec(query, args...)
	return err
}
