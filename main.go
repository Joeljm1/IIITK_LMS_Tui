package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	mainModel "github.com/Joeljm1/IIITKlmsTui/tui/models/main"
	tea "github.com/charmbracelet/bubbletea"
)

func init() {
	go func() {
		log.Println("pprof: http://localhost:6060/debug/pprof/")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}

func main() {
	f, err := tea.LogToFile("./log.txt", "Log: ")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer f.Close()
	log.SetOutput(f)
	p := tea.NewProgram(mainModel.InitialModel(), tea.WithAltScreen())
	if _, err = p.Run(); err != nil {
		f.WriteString(err.Error())
		os.Exit(1)
	}
}
