package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

const habitDayWidth = 10

type Habit struct {
	name string
}

func NewHabit(name string) Habit {
	return Habit{name: name}
}

func habitTrackerView(now time.Time, width int, focused bool, habits []Habit) string {
	contentWidth := width - 6
	if contentWidth < 0 {
		contentWidth = 0
	}
	nameWidth := contentWidth - habitDayWidth*7
	if nameWidth < 12 {
		nameWidth = 12
	}
	start := startOfWeek(now)
	days := make([]string, 0, 8)
	days = append(days, lipgloss.NewStyle().Width(nameWidth).Render(""))

	for i := 0; i < 7; i++ {
		day := start.AddDate(0, 0, i)
		label := fmt.Sprintf("%s %02d.%02d", weekdayLabel(day), day.Day(), day.Month())
		style := lipgloss.NewStyle().Width(habitDayWidth).Align(lipgloss.Center)
		if sameDay(day, now) {
			style = style.Foreground(lipgloss.Color("62")).Bold(true)
		}
		days = append(days, style.Render(label))
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render("Habit tracker"),
		lipgloss.JoinHorizontal(lipgloss.Top, days...),
		habitsView(habits, nameWidth),
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

func habitsView(habits []Habit, nameWidth int) string {
	if len(habits) == 0 {
		return lipgloss.NewStyle().Faint(true).Render("No habits yet. Press n to add one.")
	}

	rows := make([]string, 0, len(habits))
	for _, habit := range habits {
		cells := make([]string, 0, 8)
		cells = append(cells, lipgloss.NewStyle().Width(nameWidth).Render(habit.name))
		for i := 0; i < 7; i++ {
			cells = append(cells, lipgloss.NewStyle().Width(habitDayWidth).Align(lipgloss.Center).Render("□"))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
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
