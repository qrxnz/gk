package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type Habit struct {
	name string
}

func NewHabit(name string) Habit {
	return Habit{name: name}
}

func habitTrackerView(now time.Time, width int, focused bool, habits []Habit) string {
	start := startOfWeek(now)
	days := make([]string, 0, 7)

	for i := 0; i < 7; i++ {
		day := start.AddDate(0, 0, i)
		label := fmt.Sprintf("%s %02d.%02d", weekdayLabel(day), day.Day(), day.Month())
		style := lipgloss.NewStyle().Padding(0, 1)
		if sameDay(day, now) {
			style = style.Foreground(lipgloss.Color("62")).Bold(true)
		}
		days = append(days, style.Render(label))
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render("Habit tracker"),
		strings.Join(days, " "),
		habitsView(habits),
	)

	style := lipgloss.NewStyle().
		MarginTop(1).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		Width(width).
		BorderForeground(lipgloss.Color("62"))
	if !focused {
		style = style.Border(lipgloss.HiddenBorder())
	}

	return style.
		Render(content)
}

func habitsView(habits []Habit) string {
	if len(habits) == 0 {
		return lipgloss.NewStyle().Faint(true).Render("No habits yet. Press n to add one.")
	}

	rows := make([]string, 0, len(habits))
	for _, habit := range habits {
		rows = append(rows, fmt.Sprintf("%s  %s", habit.name, strings.Repeat("□ ", 7)))
	}
	return strings.Join(rows, "\n")
}

func startOfWeek(t time.Time) time.Time {
	date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return date.AddDate(0, 0, 1-weekday)
}

func weekdayLabel(t time.Time) string {
	labels := [...]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	return labels[t.Weekday()]
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}
