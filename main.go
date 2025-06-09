package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return ""
}

func main() {
	f, err := tea.LogToFile("./log.txt", "Log: ")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer f.Close()
	p := tea.NewProgram(model{})
	if _, err = p.Run(); err != nil {
		f.WriteString(err.Error())
		os.Exit(1)
	}
}
