package mainModel

import (
	"fmt"
	"time"

	"github.com/Joeljm1/IIITKlmsTui/tui/models/login"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	sp        spinner.Model
	login     login.Model
	isLoading bool
}

func InitialModel() model {
	sp := spinner.New()
	mySpinner := spinner.Spinner{
		Frames: []string{"L", "Lo", "Loa", "Load", "Loadi", "Loadin", "Loading", "Loading😎"},
		FPS:    time.Second / 5,
	}
	sp.Spinner = mySpinner
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		sp:        sp,
		login:     login.InitialModel(),
		isLoading: true,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.login.Init(), m.sp.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case login.LoginErr, login.UName, login.Psswd, login.LoginValidErr, login.Valid, tea.WindowSizeMsg:
		// m.login,cmd=m.Update(msg) do it with type inference
		if _, ok := msg.(login.Valid); ok {
			m.isLoading = false
		}
		teaModel, cmd := m.login.Update(msg)
		m.login, _ = teaModel.(login.Model)
		m.isLoading = false
		return m, cmd
	case login.LoginComplete:
		m.isLoading = false
		m.login.Err = nil
		return m, nil
	case login.Load:
		m.isLoading = true
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.sp, cmd = m.sp.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		if m.login.Err != nil {
			teaModel, cmd := m.login.Update(msg)
			m.login, _ = teaModel.(login.Model)
			return m, cmd
		}
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.isLoading {
		return m.sp.View()
	}
	if m.login.Err != nil {
		return m.login.View()
	}
	return fmt.Sprintf("username: %v\n password: %v", m.login.Username, m.login.Psswd)
}
