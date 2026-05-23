package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

const habitDayWidth = 10

type Habit struct {
	id      int64
	name    string
	checked [7]bool
	streak  int
}

func NewHabit(name string) Habit {
	return Habit{name: name}
}

func NewHabitWithChecks(name string, checked [7]bool) Habit {
	return Habit{name: name, checked: checked}
}

func (h Habit) withChecks(checked [7]bool) Habit {
	h.checked = checked
	return h
}

func habitTrackerView(now time.Time, width, height int, focused bool, habits []Habit, selectedHabit, selectedDay int) string {
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
		if focused && i == selectedDay {
			style = style.Foreground(lipgloss.Color("62")).Bold(true)
		}
		days = append(days, style.Render(label))
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render("Habit tracker"),
		lipgloss.JoinHorizontal(lipgloss.Top, days...),
		habitsView(now, habits, nameWidth, visibleHabitRows(height), focused, selectedHabit, selectedDay),
	)

	style := lipgloss.NewStyle().
		MarginTop(1).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		Height(height).
		Width(width).
		BorderForeground(lipgloss.Color("62"))
	if !focused {
		style = style.Border(lipgloss.HiddenBorder())
	}

	return style.
		Render(content)
}

func habitsView(now time.Time, habits []Habit, nameWidth, visibleRows int, focused bool, selectedHabit, selectedDay int) string {
	if len(habits) == 0 {
		return lipgloss.NewStyle().Faint(true).Render("No habits yet. Press n to add one.")
	}

	start, end := habitPageBounds(len(habits), visibleRows, selectedHabit)
	rows := make([]string, 0, visibleRows+1)
	for i, habit := range habits[start:end] {
		index := start + i
		selectedRow := focused && index == selectedHabit
		nameStyle := lipgloss.NewStyle().Width(nameWidth)
		if selectedRow {
			nameStyle = nameStyle.Foreground(lipgloss.Color("62")).Bold(true)
		}
		cells := make([]string, 0, 8)
		cells = append(cells, nameStyle.Render(fmt.Sprintf("%s  %dd", habit.name, habit.streak)))
		for day := 0; day < 7; day++ {
			cellStyle := lipgloss.NewStyle().Width(habitDayWidth).Align(lipgloss.Center)
			if selectedRow && day == selectedDay {
				cellStyle = cellStyle.Foreground(lipgloss.Color("62")).Bold(true)
			}
			mark := "□"
			if habit.checked[day] {
				mark = "■"
			}
			cells = append(cells, cellStyle.Render(mark))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}
	for len(rows) < visibleRows {
		rows = append(rows, "")
	}
	if len(habits) > visibleRows {
		page := selectedHabit/visibleRows + 1
		pages := (len(habits) + visibleRows - 1) / visibleRows
		rows = append(rows, pageDots(page, pages))
	}
	return strings.Join(rows, "\n")
}

func pageDots(page, pages int) string {
	dots := make([]string, 0, pages)
	for i := 1; i <= pages; i++ {
		if i == page {
			dots = append(dots, lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#847A85", Dark: "#979797"}).Render("•"))
		} else {
			dots = append(dots, lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}).Render("•"))
		}
	}
	return lipgloss.NewStyle().PaddingLeft(2).Render(strings.Join(dots, ""))
}

func visibleHabitRows(height int) int {
	rows := height - 6
	if rows < 1 {
		return 1
	}
	return rows
}

func habitPageBounds(total, visibleRows, selected int) (int, int) {
	if visibleRows < 1 {
		visibleRows = 1
	}
	if selected < 0 {
		selected = 0
	}
	if selected >= total {
		selected = total - 1
	}
	start := selected / visibleRows * visibleRows
	end := start + visibleRows
	if end > total {
		end = total
	}
	return start, end
}

func habitStreak(now time.Time, habit Habit) int {
	day := weekdayIndex(now)
	streak := 0
	for i := day; i >= 0; i-- {
		if !habit.checked[i] {
			break
		}
		streak++
	}
	return streak
}

func weekdayIndex(t time.Time) int {
	weekday := int(t.Weekday())
	if weekday == 0 {
		return 6
	}
	return weekday - 1
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
