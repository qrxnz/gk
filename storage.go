package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

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

func (s *Storage) LoadHabits() ([]Habit, error) {
	rows, err := s.db.Query(`
		SELECT name
		FROM habits
		ORDER BY position, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	habits := []Habit{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		habits = append(habits, NewHabit(name))
	}

	return habits, rows.Err()
}

func (s *Storage) SaveHabits(habits []Habit) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM habits`); err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO habits (name, position)
		VALUES (?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for position, habit := range habits {
		if _, err := stmt.Exec(habit.name, position); err != nil {
			return err
		}
	}

	return tx.Commit()
}
