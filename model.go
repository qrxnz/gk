package main

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Board struct {
	help               help.Model
	loaded             bool
	height             int
	focused            status
	habitFocused       bool
	habitIndex         int
	habitDay           int
	habitWeek          time.Time
	habitToday         []Habit
	habits             []Habit
	confirmDeleteHabit bool
	cols               []column
	storage            *Storage
	err                error
	quitting           bool
}

func NewBoard(storage *Storage) *Board {
	help := help.New()
	help.ShowAll = true
	now := time.Now()
	return &Board{help: help, focused: todo, habitWeek: startOfWeek(now), habitDay: weekdayIndex(now), storage: storage}
}

func (m *Board) Init() tea.Cmd {
	return nil
}

func (m *Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		var cmds []tea.Cmd
		m.help.Width = msg.Width - margin
		m.height = msg.Height
		for i := 0; i < len(m.cols); i++ {
			var res tea.Model
			res, cmd = m.cols[i].Update(msg)
			m.cols[i] = res.(column)
			cmds = append(cmds, cmd)
		}
		m.loaded = true
		return m, tea.Batch(cmds...)
	case Form:
		if m.habitFocused {
			return m, nil
		}
		return m, tea.Sequence(m.cols[m.focused].Set(msg.index, msg.CreateTask()), m.save())
	case HabitForm:
		m.habits = append(m.habits, msg.CreateHabit())
		m.habitIndex = len(m.habits) - 1
		return m, tea.Sequence(m.saveHabits(), m.loadHabitChecks())
	case moveMsg:
		return m, tea.Sequence(m.cols[m.focused.getNext()].Set(APPEND, msg.Task), m.save())
	case saveMsg:
		return m, m.save()
	case habitsSavedMsg:
		return m, nil
	case habitsLoadedMsg:
		m.habits = msg.habits
		m.habitToday = msg.today
		if m.habitIndex >= len(m.habits) {
			m.habitIndex = len(m.habits) - 1
		}
		if m.habitIndex < 0 {
			m.habitIndex = 0
		}
		return m, nil
	case error:
		m.err = msg
	case tea.KeyMsg:
		if m.confirmDeleteHabit {
			switch {
			case key.Matches(msg, keys.Quit):
				m.quitting = true
				return m, tea.Quit
			case key.Matches(msg, keys.Back), key.Matches(msg, keys.No):
				m.confirmDeleteHabit = false
				return m, nil
			case key.Matches(msg, keys.Enter), key.Matches(msg, keys.Yes):
				m.confirmDeleteHabit = false
				return m, m.deleteHabit()
			}
			return m, nil
		}
		switch {
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, keys.New):
			if m.habitFocused {
				return newHabitForm().Update(nil)
			}
		case key.Matches(msg, keys.Delete):
			if m.habitFocused {
				if len(m.habits) > 0 {
					m.confirmDeleteHabit = true
				}
				return m, nil
			}
		case key.Matches(msg, keys.Tab):
			if !m.habitFocused {
				m.cols[m.focused].Blur()
				m.habitFocused = true
				return m, nil
			}
			m.habitFocused = false
			m.cols[m.focused].Focus()
			return m, nil
		case key.Matches(msg, keys.Up):
			if m.habitFocused {
				m.selectPrevHabit()
				return m, nil
			}
		case key.Matches(msg, keys.Down):
			if m.habitFocused {
				m.selectNextHabit()
				return m, nil
			}
		case key.Matches(msg, keys.Enter):
			if m.habitFocused {
				return m, m.toggleHabitCheck()
			}
		case key.Matches(msg, keys.Left):
			if m.habitFocused {
				if m.habitDay == 0 {
					m.habitWeek = m.habitWeek.AddDate(0, 0, -7)
					m.habitDay = 6
					return m, m.loadHabitChecks()
				}
				m.habitDay--
				return m, nil
			}
			m.cols[m.focused].Blur()
			m.focused = m.focused.getPrev()
			m.cols[m.focused].Focus()
		case key.Matches(msg, keys.Right):
			if m.habitFocused {
				if m.habitDay == 6 {
					if !m.habitWeek.Before(startOfWeek(time.Now())) {
						return m, nil
					}
					m.habitWeek = m.habitWeek.AddDate(0, 0, 7)
					m.habitDay = 0
					return m, m.loadHabitChecks()
				}
				m.habitDay++
				return m, nil
			}
			m.cols[m.focused].Blur()
			m.focused = m.focused.getNext()
			m.cols[m.focused].Focus()
		}
	}
	if m.habitFocused {
		return m, nil
	}
	res, cmd := m.cols[m.focused].Update(msg)
	if _, ok := res.(column); ok {
		m.cols[m.focused] = res.(column)
	} else {
		return res, cmd
	}
	return m, cmd
}

func (m *Board) save() tea.Cmd {
	return func() tea.Msg {
		if err := m.storage.Save(m.cols); err != nil {
			return err
		}
		return nil
	}
}

func (m *Board) saveHabits() tea.Cmd {
	return func() tea.Msg {
		if err := m.storage.SaveHabits(m.habits, m.habitWeek); err != nil {
			return err
		}
		return nil
	}
}

func (m *Board) loadHabitChecks() tea.Cmd {
	return func() tea.Msg {
		habits, err := m.storage.LoadHabits(m.habitWeek)
		if err != nil {
			return err
		}
		today, err := m.storage.LoadHabits(startOfWeek(time.Now()))
		if err != nil {
			return err
		}
		return habitsLoadedMsg{habits: habits, today: today}
	}
}

func (m *Board) selectPrevHabit() {
	if len(m.habits) == 0 {
		m.habitIndex = 0
		return
	}
	m.habitIndex--
	if m.habitIndex < 0 {
		m.habitIndex = len(m.habits) - 1
	}
}

func (m *Board) selectNextHabit() {
	if len(m.habits) == 0 {
		m.habitIndex = 0
		return
	}
	m.habitIndex++
	if m.habitIndex >= len(m.habits) {
		m.habitIndex = 0
	}
}

func (m *Board) toggleHabitCheck() tea.Cmd {
	if len(m.habits) == 0 {
		return nil
	}
	m.habits[m.habitIndex].checked[m.habitDay] = !m.habits[m.habitIndex].checked[m.habitDay]
	if sameDay(m.habitWeek.AddDate(0, 0, m.habitDay), time.Now()) {
		m.habitToday = m.habits
	}
	return tea.Sequence(m.saveHabits(), m.loadHabitChecks())
}

func (m *Board) deleteHabit() tea.Cmd {
	if len(m.habits) == 0 {
		return nil
	}
	m.habits = append(m.habits[:m.habitIndex], m.habits[m.habitIndex+1:]...)
	if m.habitIndex >= len(m.habits) {
		m.habitIndex = len(m.habits) - 1
	}
	if m.habitIndex < 0 {
		m.habitIndex = 0
	}
	return tea.Sequence(m.saveHabits(), m.loadHabitChecks())
}

type habitsSavedMsg struct{}

type habitsLoadedMsg struct {
	habits []Habit
	today  []Habit
}

// Changing to pointer receiver to get back to this model after adding a new task via the form... Otherwise I would need to pass this model along to the form and it becomes highly coupled to the other models.
func (m *Board) View() string {
	if m.quitting {
		return ""
	}
	if m.err != nil {
		return m.err.Error()
	}
	if m.confirmDeleteHabit {
		return m.deleteHabitConfirmationView()
	}
	if !m.loaded {
		return "loading..."
	}
	board := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.cols[todo].View(),
		m.cols[inProgress].View(),
		m.cols[done].View(),
	)
	now := time.Now()
	habitWeek := m.habitWeek
	if habitWeek.IsZero() {
		habitWeek = startOfWeek(now)
	}
	return lipgloss.JoinVertical(lipgloss.Left, board, habitTrackerView(habitWeek, lipgloss.Width(board), m.habitHeight(), m.habitFocused, m.habits, m.habitIndex, m.habitDay, remainingHabitsToday(time.Now(), m.habitToday)), m.help.View(keys))
}

func (m *Board) habitHeight() int {
	if m.height <= 0 {
		return 8
	}
	height := m.height / 4
	if height < 8 {
		return 8
	}
	return height
}

func (m *Board) deleteHabitConfirmationView() string {
	name := "habit"
	if len(m.habits) > 0 {
		name = m.habits[m.habitIndex].name
	}
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().Bold(true).Render("Delete habit?"),
		"Are you sure you want to delete \""+name+"\"?",
		"y/enter confirm  n/esc cancel",
	)
	modal := lipgloss.NewStyle().
		Padding(1, 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(focusColor)).
		Render(content)
	return lipgloss.Place(m.help.Width+margin, m.height, lipgloss.Center, lipgloss.Center, modal)
}
