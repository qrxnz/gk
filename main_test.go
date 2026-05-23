package main

import "testing"

func TestStatusGetNextWrapsFromDoneToTodo(t *testing.T) {
	tests := []struct {
		name string
		stat status
		want status
	}{
		{name: "todo", stat: todo, want: inProgress},
		{name: "in progress", stat: inProgress, want: done},
		{name: "done", stat: done, want: todo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.stat.getNext(); got != tt.want {
				t.Fatalf("getNext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusGetPrevWrapsFromTodoToDone(t *testing.T) {
	tests := []struct {
		name string
		stat status
		want status
	}{
		{name: "todo", stat: todo, want: done},
		{name: "in progress", stat: inProgress, want: todo},
		{name: "done", stat: done, want: inProgress},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.stat.getPrev(); got != tt.want {
				t.Fatalf("getPrev() = %v, want %v", got, tt.want)
			}
		})
	}
}
