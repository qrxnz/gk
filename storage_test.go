package main

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/tursodatabase/go-libsql"
)

func newTestStorage(t *testing.T) *Storage {
	t.Helper()

	db, err := sql.Open("libsql", "file:"+t.TempDir()+"/gk.db")
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}

	storage := &Storage{db: db}
	if err := storage.init(); err != nil {
		db.Close()
		t.Fatalf("storage.init() error = %v", err)
	}

	t.Cleanup(func() {
		if err := storage.Close(); err != nil {
			t.Fatalf("storage.Close() error = %v", err)
		}
	})

	return storage
}

func TestStorageLoadReturnsEmptyStatuses(t *testing.T) {
	storage := newTestStorage(t)

	tasks, err := storage.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	for _, stat := range []status{todo, inProgress, done} {
		if _, ok := tasks[stat]; !ok {
			t.Fatalf("Load() missing status %v", stat)
		}
		if len(tasks[stat]) != 0 {
			t.Fatalf("len(tasks[%v]) = %d, want 0", stat, len(tasks[stat]))
		}
	}
}

func TestStorageSaveAndLoadRoundTrip(t *testing.T) {
	storage := newTestStorage(t)
	cols := []column{newColumn(todo), newColumn(inProgress), newColumn(done)}
	runCmd(t, cols[todo].Set(APPEND, NewTask(todo, "todo 1", "todo description 1")))
	runCmd(t, cols[todo].Set(APPEND, NewTask(todo, "todo 2", "todo description 2")))
	runCmd(t, cols[inProgress].Set(APPEND, NewTask(inProgress, "progress", "progress description")))
	runCmd(t, cols[done].Set(APPEND, NewTask(done, "done", "done description")))

	if err := storage.Save(cols); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	tasks, err := storage.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	assertLoadedTask(t, tasks[todo][0], todo, "todo 2", "todo description 2")
	assertLoadedTask(t, tasks[todo][1], todo, "todo 1", "todo description 1")
	assertLoadedTask(t, tasks[inProgress][0], inProgress, "progress", "progress description")
	assertLoadedTask(t, tasks[done][0], done, "done", "done description")
}

func TestStorageSaveReplacesExistingTasks(t *testing.T) {
	storage := newTestStorage(t)
	cols := []column{newColumn(todo), newColumn(inProgress), newColumn(done)}
	runCmd(t, cols[todo].Set(APPEND, NewTask(todo, "old", "old description")))
	if err := storage.Save(cols); err != nil {
		t.Fatalf("first Save() error = %v", err)
	}

	cols = []column{newColumn(todo), newColumn(inProgress), newColumn(done)}
	runCmd(t, cols[done].Set(APPEND, NewTask(done, "new", "new description")))
	if err := storage.Save(cols); err != nil {
		t.Fatalf("second Save() error = %v", err)
	}

	tasks, err := storage.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(tasks[todo]) != 0 {
		t.Fatalf("len(tasks[todo]) = %d, want 0", len(tasks[todo]))
	}
	if len(tasks[done]) != 1 {
		t.Fatalf("len(tasks[done]) = %d, want 1", len(tasks[done]))
	}
	assertLoadedTask(t, tasks[done][0], done, "new", "new description")
}

func TestStorageSaveAndLoadHabitsRoundTrip(t *testing.T) {
	storage := newTestStorage(t)
	checked := [7]bool{true, false, true, false, false, true, false}
	week := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)

	if err := storage.SaveHabits([]Habit{NewHabitWithChecks("Read", checked), NewHabit("Run")}, week); err != nil {
		t.Fatalf("SaveHabits() error = %v", err)
	}

	habits, err := storage.LoadHabits(week)
	if err != nil {
		t.Fatalf("LoadHabits() error = %v", err)
	}
	if len(habits) != 2 {
		t.Fatalf("len(habits) = %d, want 2", len(habits))
	}
	if habits[0].name != "Read" || habits[1].name != "Run" {
		t.Fatalf("habits = %#v, want Read and Run", habits)
	}
	if habits[0].checked != checked {
		t.Fatalf("checked = %#v, want %#v", habits[0].checked, checked)
	}
}

func TestStorageHabitStreakUsesStableHabitIDAcrossWeeks(t *testing.T) {
	storage := newTestStorage(t)
	week := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	prevWeek := week.AddDate(0, 0, -7)
	checked := [7]bool{true, true, true, true, true, true, true}

	if err := storage.SaveHabits([]Habit{NewHabitWithChecks("Read", checked)}, prevWeek); err != nil {
		t.Fatalf("SaveHabits(prevWeek) error = %v", err)
	}
	habits, err := storage.LoadHabits(prevWeek)
	if err != nil {
		t.Fatalf("LoadHabits(prevWeek) error = %v", err)
	}
	if err := storage.SaveHabits([]Habit{habits[0].withChecks(checked)}, week); err != nil {
		t.Fatalf("SaveHabits(week) error = %v", err)
	}
	habits, err = storage.LoadHabitsWithStreakDate(week, time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("LoadHabits(week) error = %v", err)
	}
	if habits[0].streak != 10 {
		t.Fatalf("streak = %d, want 10", habits[0].streak)
	}
}

func TestStorageSaveHabitsDeletesRemovedHabits(t *testing.T) {
	storage := newTestStorage(t)
	week := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)

	if err := storage.SaveHabits([]Habit{NewHabit("Read"), NewHabit("Run")}, week); err != nil {
		t.Fatalf("SaveHabits() error = %v", err)
	}
	habits, err := storage.LoadHabits(week)
	if err != nil {
		t.Fatalf("LoadHabits() error = %v", err)
	}
	if err := storage.SaveHabits([]Habit{habits[1]}, week); err != nil {
		t.Fatalf("SaveHabits(delete) error = %v", err)
	}
	habits, err = storage.LoadHabits(week)
	if err != nil {
		t.Fatalf("LoadHabits() error = %v", err)
	}
	if len(habits) != 1 || habits[0].name != "Run" {
		t.Fatalf("habits = %#v, want only Run", habits)
	}
}

func assertLoadedTask(t *testing.T, item interface{}, stat status, title, description string) {
	t.Helper()

	task, ok := item.(Task)
	if !ok {
		t.Fatalf("item type = %T, want Task", item)
	}
	if task.status != stat {
		t.Fatalf("status = %v, want %v", task.status, stat)
	}
	if task.title != title {
		t.Fatalf("title = %q, want %q", task.title, title)
	}
	if task.description != description {
		t.Fatalf("description = %q, want %q", task.description, description)
	}
}
