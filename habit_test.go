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
	view := habitTrackerView(time.Date(2026, 5, 23, 12, 0, 0, 0, time.UTC), 120, true)

	for _, want := range []string{"Habit tracker", "Mon 18.05", "Tue 19.05", "Wed 20.05", "Thu 21.05", "Fri 22.05", "Sat 23.05", "Sun 24.05"} {
		if !strings.Contains(view, want) {
			t.Fatalf("habitTrackerView() missing %q in %q", want, view)
		}
	}
}
