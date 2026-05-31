package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestFormAllowsQInTaskTitle(t *testing.T) {
	form := newDefaultForm()
	model, cmd := form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}, Alt: false})
	updated := model.(Form)

	if updated.title.Value() != "q" {
		t.Fatalf("title = %q, want q", updated.title.Value())
	}
	if cmd == nil {
		return
	}
	if msg := cmd(); msg == tea.Quit() {
		t.Fatal("q should not quit while editing task title")
	}
}

func TestHabitFormAllowsQInHabitName(t *testing.T) {
	form := newHabitForm()
	model, cmd := form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}, Alt: false})
	updated := model.(HabitForm)

	if updated.name.Value() != "q" {
		t.Fatalf("name = %q, want q", updated.name.Value())
	}
	if cmd == nil {
		return
	}
	if msg := cmd(); msg == tea.Quit() {
		t.Fatal("q should not quit while editing habit name")
	}
}
