package main

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Board struct {
	help         help.Model
	loaded       bool
	focused      status
	habitFocused bool
	habits       []Habit
	cols         []column
	storage      *Storage
	err          error
	quitting     bool
}

func NewBoard(storage *Storage) *Board {
	help := help.New()
	help.ShowAll = true
	return &Board{help: help, focused: todo, storage: storage}
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
		return m, m.saveHabits()
	case moveMsg:
		return m, tea.Sequence(m.cols[m.focused.getNext()].Set(APPEND, msg.Task), m.save())
	case saveMsg:
		return m, m.save()
	case error:
		m.err = msg
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, keys.New):
			if m.habitFocused {
				return newHabitForm().Update(nil)
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
		case key.Matches(msg, keys.Left):
			if m.habitFocused {
				return m, nil
			}
			m.cols[m.focused].Blur()
			m.focused = m.focused.getPrev()
			m.cols[m.focused].Focus()
		case key.Matches(msg, keys.Right):
			if m.habitFocused {
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
		if err := m.storage.SaveHabits(m.habits); err != nil {
			return err
		}
		return nil
	}
}

// Changing to pointer receiver to get back to this model after adding a new task via the form... Otherwise I would need to pass this model along to the form and it becomes highly coupled to the other models.
func (m *Board) View() string {
	if m.quitting {
		return ""
	}
	if m.err != nil {
		return m.err.Error()
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
	return lipgloss.JoinVertical(lipgloss.Left, board, habitTrackerView(time.Now(), lipgloss.Width(board), m.habitFocused, m.habits), m.help.View(keys))
}
