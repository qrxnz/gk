package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

type status int

func (s status) getNext() status {
	if s == done {
		return todo
	}
	return s + 1
}

func (s status) getPrev() status {
	if s == todo {
		return done
	}
	return s - 1
}

const margin = 4

var board *Board

const (
	todo status = iota
	inProgress
	done
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	appDir := filepath.Join(homeDir, ".gk")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	f, err := tea.LogToFile(filepath.Join(appDir, "debug.log"), "debug")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	storage, err := NewStorage()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer storage.Close()

	board = NewBoard(storage)
	board.initLists()
	p := tea.NewProgram(board, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
