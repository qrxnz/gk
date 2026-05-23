package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewColumnFocusesTodoOnly(t *testing.T) {
	tests := []struct {
		name    string
		stat    status
		focused bool
	}{
		{name: "todo", stat: todo, focused: true},
		{name: "in progress", stat: inProgress, focused: false},
		{name: "done", stat: done, focused: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col := newColumn(tt.stat)
			if col.status != tt.stat {
				t.Fatalf("status = %v, want %v", col.status, tt.stat)
			}
			if col.Focused() != tt.focused {
				t.Fatalf("Focused() = %v, want %v", col.Focused(), tt.focused)
			}
		})
	}
}

func TestColumnFocusAndBlur(t *testing.T) {
	col := newColumn(done)

	col.Focus()
	if !col.Focused() {
		t.Fatal("Focused() = false, want true")
	}

	col.Blur()
	if col.Focused() {
		t.Fatal("Focused() = true, want false")
	}
}

func TestColumnSetAppendsAndReplacesItems(t *testing.T) {
	col := newColumn(todo)
	first := NewTask(todo, "first", "first description")
	second := NewTask(todo, "second", "second description")

	runCmd(t, col.Set(APPEND, first))

	items := col.list.Items()
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if got := items[0].(Task); got.title != "first" {
		t.Fatalf("appended task title = %q, want %q", got.title, "first")
	}

	runCmd(t, col.Set(0, second))
	items = col.list.Items()
	if len(items) != 1 {
		t.Fatalf("len(items) after replace = %d, want 1", len(items))
	}
	if got := items[0].(Task); got.title != "second" {
		t.Fatalf("replaced task title = %q, want %q", got.title, "second")
	}
}

func TestColumnDeleteCurrent(t *testing.T) {
	col := newColumn(todo)
	runCmd(t, col.Set(APPEND, NewTask(todo, "first", "description")))
	runCmd(t, col.Set(APPEND, NewTask(todo, "second", "description")))

	runCmd(t, col.DeleteCurrent())

	items := col.list.Items()
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if got := items[0].(Task); got.title != "first" {
		t.Fatalf("remaining task title = %q, want %q", got.title, "first")
	}
}

func TestColumnMoveToNextRemovesTaskAndReturnsMoveMessage(t *testing.T) {
	col := newColumn(inProgress)
	runCmd(t, col.Set(APPEND, NewTask(inProgress, "task", "description")))

	msg := runCmd(t, col.MoveToNext())

	if len(col.list.Items()) != 0 {
		t.Fatalf("len(items) = %d, want 0", len(col.list.Items()))
	}
	move, ok := msg.(moveMsg)
	if !ok {
		t.Fatalf("MoveToNext() message type = %T, want moveMsg", msg)
	}
	if move.Task.status != done {
		t.Fatalf("moved task status = %v, want %v", move.Task.status, done)
	}
	if move.Task.title != "task" {
		t.Fatalf("moved task title = %q, want %q", move.Task.title, "task")
	}
}

func TestColumnMoveToNextWithoutSelectedTaskReturnsNil(t *testing.T) {
	col := newColumn(todo)

	if cmd := col.MoveToNext(); cmd != nil {
		t.Fatalf("MoveToNext() = %v, want nil", cmd)
	}
}

func runCmd(t *testing.T, cmd tea.Cmd) tea.Msg {
	t.Helper()
	if cmd == nil {
		return nil
	}
	return cmd()
}
