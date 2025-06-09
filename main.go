package main

import (
	"fmt"
	"os"

	mainModel "github.com/Joeljm1/IIITKlmsTui/tui/models/main"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	f, err := tea.LogToFile("./log.txt", "Log: ")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer f.Close()
	p := tea.NewProgram(mainModel.InitialModel())
	if _, err = p.Run(); err != nil {
		f.WriteString(err.Error())
		os.Exit(1)
	}
}
