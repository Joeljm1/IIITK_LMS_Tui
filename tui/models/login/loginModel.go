package login

import (
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
	width         int
	height        int
	moreHelp      bool
}

type (
	// Sent when details not in keyring
	LoginErr struct {
		e error
	}
	// Sent when details given by user are wrong
	LoginValidErr struct {
		e error
	}
	UName         string
	Psswd         string
	LoginComplete string
	Valid         bool
	Load          struct{}
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#c20000"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	blurredButton  = blurredStyle.Render(" [ Submit] ")
	focussedButton = focusedStyle.Foreground(lipgloss.Color("#17bd05")).Render(" [ Submit] ")

	errMsgStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	loginBoxBorder = lipgloss.RoundedBorder()
	topBox         = lipgloss.NewStyle().Align(lipgloss.Center).Background(lipgloss.Color("#03005e"))
	loginBox       = lipgloss.NewStyle().Border(loginBoxBorder, true).
			BorderForeground(lipgloss.Color("#05d7e6")).Align(lipgloss.Center).Background(lipgloss.Color("#000000"))
	helpBox = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true, false, false, false).
		Background(lipgloss.Color("#03005e"))
	/*Foreground(lipgloss.Color("#B1B1B1")).*/
)

func InitialModel() Model {
	t1 := textinput.New()
	t1.CharLimit = 11
	t1.Placeholder = "Username"
	t1.PromptStyle = focusedStyle
	t1.Cursor.SetMode(0)
	t1.Focus()
	t2 := textinput.New()
	t2.CharLimit = 8
	t2.Placeholder = "Password"
	t1.Cursor.SetMode(0)
	t2.Blur()
	return Model{
		Username:      "",
		Psswd:         "",
		Err:           nil,
		unameInp:      t1,
		psswdInp:      t2,
		focus:         0,
		validationErr: nil,
		moreHelp:      false,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.sendUsername, m.sendPasswd, textinput.Blink)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height, m.width = msg.Height, msg.Width
	case UName:
		m.Username = string(msg)
		if m.Psswd != "" {
			return m, tea.Batch(m.LoginComplete, m.loginWithCLient)
		}
	case Psswd:
		m.Psswd = string(msg)
		if m.Username != "" {
			return m, tea.Batch(m.LoginComplete, m.loginWithCLient)
		}
	case LoginErr: // no details in keyring
		m.Err = msg.e
	case LoginValidErr:
		m.validationErr = msg.e
	case Valid: // got valid details
		return m, tea.Batch(m.sendUsername, m.sendPasswd)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "?":
			m.moreHelp = !m.moreHelp
		case "up", "ctrl+shift+tab", "down", "tab", "enter":
			val := msg.String()
			if m.focus == 2 && val == "enter" {
				// TODO validate uname and psswd
				uname := m.unameInp.Value()
				psswd := m.psswdInp.Value()
				m.validationErr = nil // to remove error
				return m, tea.Batch(m.validateDetails(uname, psswd), m.load)
			}
			if val == "up" || val == "ctrl+shift+tab" {
				m.focus = max(m.focus-1, 0)
			} else {
				m.focus = min(m.focus+1, 2)
			}
			if m.focus == 0 {
				cmd = m.unameInp.Focus()
				m.psswdInp.Blur()
				m.unameInp.PromptStyle = focusedStyle
				m.unameInp.TextStyle = noStyle
				m.psswdInp.PromptStyle = noStyle
				m.psswdInp.TextStyle = noStyle
			} else if m.focus == 1 {
				m.unameInp.Blur()
				cmd = m.psswdInp.Focus()
				m.unameInp.PromptStyle = noStyle
				m.unameInp.TextStyle = noStyle
				m.psswdInp.PromptStyle = focusedStyle
				m.psswdInp.TextStyle = noStyle
			} else {
				m.unameInp.Blur()
				m.psswdInp.Blur()
				m.unameInp.PromptStyle = noStyle
				m.unameInp.TextStyle = noStyle
				m.psswdInp.PromptStyle = noStyle
				m.psswdInp.TextStyle = noStyle
			}
			return m, cmd
		}
		var cmd1 tea.Cmd
		var cmd2 tea.Cmd
		if msg.String() != "?" {
			m.unameInp, cmd1 = m.unameInp.Update(msg)
			m.psswdInp, cmd2 = m.psswdInp.Update(msg)
		}
		return m, tea.Batch(cmd1, cmd2)
	}
	return m, nil
}

func (m Model) View() string {
	if m.Err != nil {
		sb := strings.Builder{}
		sb.WriteString("Username: ")
		sb.WriteString(m.unameInp.View())
		sb.WriteRune('\n')
		sb.WriteString("Password: ")
		sb.WriteString(m.psswdInp.View())
		sb.WriteRune('\n')
		button := blurredButton
		// username=0, password=1 then submit =2
		if m.focus == 2 {
			button = focussedButton
		}
		sb.WriteString(button)
		if m.validationErr != nil {
			sb.WriteRune('\n')
			sb.WriteRune('\n')
			sb.WriteString(errMsgStyle.Render("Your username or password is invalid"))
		}
		var helpView string
		if m.moreHelp {
			helpView = " ↑      move up      ←     move left\n ↓      move down    →     move right\n ctrl+c quit         enter submit (on submit button)"
		} else {
			helpView = "↑ move up      ↓ move down      ctrl+c quit      ? help"
		}
		var red int
		if m.moreHelp {
			red = 3
		} else {
			red = 1
		}
		topBoxStyle := topBox.Width(m.width).Height(m.height-red).Padding(m.height/4, 0, m.height/4-2-red, 0)
		loginBoxStyle := loginBox.Width(m.width/2).Height(m.height/2).Padding(m.height/8, 0)
		helpBoxStyle := helpBox.Width(m.width)
		return topBoxStyle.Render(loginBoxStyle.Render(sb.String())) + helpBoxStyle.Render(helpView)
	}
	return ""
}
