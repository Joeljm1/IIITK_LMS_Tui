package login

import (
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Username      string
	Psswd         string
	Err           error // when this err not nil username and psswd is asked from user
	unameInp      textinput.Model
	psswdInp      textinput.Model
	focus         int
	validationErr error
	W             io.Writer
}

type (
	// Sent when details not in keyring
	LoginErr error
	// Sent when details given by user are wrong
	LoginValidErr error
	UName         string
	Psswd         string
	LoginComplete string
	Valid         bool
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	blurredButton  = blurredStyle.Render(" [ Submit] ")
	focussedButton = focusedStyle.Render(" [ Submit] ")

	errMsgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
)

func InitialModel() Model {
	t1 := textinput.New()
	t1.CharLimit = 11
	t1.Placeholder = "Username"
	t1.PromptStyle = focusedStyle
	t1.TextStyle = focusedStyle
	t1.Focus()
	t2 := textinput.New()
	t1.CharLimit = 8
	t1.Placeholder = "Password"
	t2.PromptStyle = focusedStyle
	t2.TextStyle = focusedStyle
	t2.Blur()
	f, err := os.OpenFile("./debug.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	return Model{
		Username:      "",
		Psswd:         "",
		Err:           nil,
		unameInp:      textinput.New(),
		psswdInp:      textinput.New(),
		focus:         0,
		validationErr: nil,
		W:             f,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.sendUsername, m.sendPasswd, textinput.Blink)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case UName:
		m.Username = string(msg)
		if m.Psswd != "" {
			m.Err = nil
			cmd = m.LoginComplete
		}
		return m, cmd
	case Psswd:
		m.Psswd = string(msg)
		if m.Username != "" {
			m.Err = nil
			cmd = m.LoginComplete
		}
		return m, cmd
	case LoginErr: // no details in keyring
		m.Err = msg
	case LoginValidErr:
		m.validationErr = msg
	case Valid: // got valid details
		return m, tea.Batch(m.sendUsername, m.sendPasswd)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "up", "ctrl+shift+tab", "down", "shift+tab", "enter":
			val := msg.String()
			if m.focus == 2 && val == "enter" {
				// TODO validate uname and psswd
				uname := m.unameInp.Value()
				psswd := m.psswdInp.Value()

				return m, m.validateDetails(uname, psswd)
			}
			if val == "up" || val == "ctrl+shift+tab" || val == "enter" {
				m.focus = max(m.focus-1, 0)
			} else {
				m.focus = min(m.focus+1, 2)
			}
			if m.focus == 0 {
				cmd = m.unameInp.Focus()
				m.psswdInp.Blur()
			} else {
				m.unameInp.Blur()
				cmd = m.psswdInp.Focus()
			}
			return m, cmd
		}
		var cmd1 tea.Cmd
		var cmd2 tea.Cmd
		m.unameInp, cmd1 = m.unameInp.Update(msg)
		m.psswdInp, cmd2 = m.psswdInp.Update(msg)
		return m, tea.Batch(cmd1, cmd2)
	}
	return m, nil
}

func (m Model) View() string {
	if m.Err != nil {
		m.W.Write([]byte("Test"))
		sb := strings.Builder{}
		sb.WriteString(m.unameInp.View())
		sb.WriteRune('\n')
		sb.WriteString(m.psswdInp.View())
		sb.WriteRune('\n')
		button := blurredButton
		// username=0, password=1 then submit =2
		if m.focus == 2 {
			button = focussedButton
		}
		sb.WriteString(button)
		if m.validationErr != nil {
			sb.WriteString(errMsgStyle.Render("Your username or password is invalid"))
		}
		return sb.String()
	}
	return ""
}
