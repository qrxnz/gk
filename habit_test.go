package main

import (
	"strings"
	"testing"
	"time"
)

func TestStartOfWeekReturnsMonday(t *testing.T) {
	loc := time.FixedZone("test", 0)
	tests := []struct {
		name string
		date time.Time
		want time.Time
	}{
		{name: "monday", date: time.Date(2026, 5, 18, 15, 30, 0, 0, loc), want: time.Date(2026, 5, 18, 0, 0, 0, 0, loc)},
		{name: "sunday", date: time.Date(2026, 5, 24, 15, 30, 0, 0, loc), want: time.Date(2026, 5, 18, 0, 0, 0, 0, loc)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := startOfWeek(tt.date); !got.Equal(tt.want) {
				t.Fatalf("startOfWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHabitTrackerViewShowsCurrentWeek(t *testing.T) {
	checked := [7]bool{true, false, false, false, true, true, false}
	view := habitTrackerView(time.Date(2026, 5, 23, 12, 0, 0, 0, time.UTC), 120, true, []Habit{NewHabitWithChecks("Read", checked)}, 0, 0)

	for _, want := range []string{"Habit tracker", "Mon 18.05", "Tue 19.05", "Wed 20.05", "Thu 21.05", "Fri 22.05", "Sat 23.05", "Sun 24.05"} {
		if !strings.Contains(view, want) {
			t.Fatalf("habitTrackerView() missing %q in %q", want, view)
		}
	}
	if !strings.Contains(view, "Read") {
		t.Fatalf("habitTrackerView() missing habit name in %q", view)
	}
	if !strings.Contains(view, "2d") {
		t.Fatalf("habitTrackerView() missing streak in %q", view)
	}
}

func TestHabitStreakCountsBackFromToday(t *testing.T) {
	checked := [7]bool{true, false, true, true, true, false, true}
	habit := NewHabitWithChecks("Read", checked)

	if got := habitStreak(time.Date(2026, 5, 22, 12, 0, 0, 0, time.UTC), habit); got != 3 {
		t.Fatalf("habitStreak() = %d, want 3", got)
	}
	if got := habitStreak(time.Date(2026, 5, 23, 12, 0, 0, 0, time.UTC), habit); got != 0 {
		t.Fatalf("habitStreak() = %d, want 0", got)
	}
}
