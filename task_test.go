package main

import "testing"

func TestNewTask(t *testing.T) {
	task := NewTask(inProgress, "title", "description")

	if task.status != inProgress {
		t.Fatalf("status = %v, want %v", task.status, inProgress)
	}
	if task.Title() != "title" {
		t.Fatalf("Title() = %q, want %q", task.Title(), "title")
	}
	if task.Description() != "description" {
		t.Fatalf("Description() = %q, want %q", task.Description(), "description")
	}
	if task.FilterValue() != "title" {
		t.Fatalf("FilterValue() = %q, want %q", task.FilterValue(), "title")
	}
}

func TestTaskNextWrapsThroughStatuses(t *testing.T) {
	task := NewTask(todo, "title", "description")

	task.Next()
	if task.status != inProgress {
		t.Fatalf("after first Next status = %v, want %v", task.status, inProgress)
	}

	task.Next()
	if task.status != done {
		t.Fatalf("after second Next status = %v, want %v", task.status, done)
	}

	task.Next()
	if task.status != todo {
		t.Fatalf("after third Next status = %v, want %v", task.status, todo)
	}
}
