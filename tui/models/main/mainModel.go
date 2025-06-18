package mainModel

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Joeljm1/IIITKlmsTui/internal/client"
	"github.com/Joeljm1/IIITKlmsTui/tui/models/courses"
	"github.com/Joeljm1/IIITKlmsTui/tui/models/login"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	height, width int
	sp            spinner.Model
	login         login.Model
	isLoading     bool
	client        *client.LMSCLient
	courseModel   courses.Model
	err           error
	attendance    [][]client.Attendance // TODO: Put it seperate model for viewing later with course name
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
		sp:          sp,
		login:       login.InitialModel(),
		courseModel: courses.InitialModel(),
		isLoading:   true,
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.login.Init(), m.sp.Tick, m.courseModel.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		m.err = msg
	case login.UName, login.Psswd, login.Valid:
		teaModel, cmd := m.login.Update(msg)
		m.login, _ = teaModel.(login.Model)
		return m, cmd
	case login.LoginValidErr:
		m.isLoading = false
		teaModel, cmd := m.login.Update(msg)
		m.login, _ = teaModel.(login.Model)
		return m, cmd
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		if !m.courseModel.RecievedCourses {
			teaModel, cmd := m.login.Update(msg)
			m.login = teaModel.(login.Model)
			return m, cmd
		} else {
			courseModel, cmd := m.courseModel.Update(msg)
			m.courseModel = courseModel.(courses.Model)
			return m, cmd
		}
	case login.LoginErr:
		teaModel, cmd := m.login.Update(msg)
		m.login, _ = teaModel.(login.Model)
		m.isLoading = false
		return m, cmd
	case login.LoginComplete:
		m.login.Err = nil
	case login.Load:
		m.isLoading = true
	case courses.Load:
		m.isLoading = true
		return m, courses.GetAttendanceList(m.courseModel.Chosen, m.client)
	case *client.LMSCLient:
		m.client = msg
		return m, courses.CheckChoiceFile(m.client)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.sp, cmd = m.sp.Update(msg)
		return m, cmd

	case courses.ChoiceFromLMSNeeded:
		return m, courses.GetCoursesFromLms(m.client)
	case *client.Choices:
		m.isLoading = false
		m.client.Choices = msg // make pointer check for error later
		m.client.RecivedChoices = true
		m.isLoading = true //?????????
		return m, courses.GetAllAttendance(m.client)
	case client.CourseList:
		m.courseModel.RecievedCourses = true
		var l []list.Item
		for i := range msg {
			it := courses.CourseItem(msg[i])
			l = append(l, it)
		}
		h, v := courses.DocStyle.GetFrameSize()
		m.courseModel.List = list.New(l, list.NewDefaultDelegate(), m.width-h, m.height-v)
		m.courseModel.List.AdditionalShortHelpKeys = func() []key.Binding {
			return []key.Binding{courses.SubmitKey, courses.SelectKey, courses.LogoutKey}
		}
		m.courseModel.List.AdditionalFullHelpKeys = func() []key.Binding {
			return []key.Binding{courses.SubmitKey, courses.SelectKey, courses.LogoutKey}
		}

		// m.courseModel.List.SetFilteringEnabled(false) // dont know why by filter does not work so disabled it
		m.isLoading = false
	case [][]client.Attendance:
		m.isLoading = false
		m.attendance = msg

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+d":
			m.login.DeleteLoginDet() // Handle err??
			client.DeleteCoursesFile()
			return m, tea.Quit
		}

		if m.login.Err != nil {
			teaModel, cmd := m.login.Update(msg)
			m.login, _ = teaModel.(login.Model)
			return m, cmd
		}
		if m.courseModel.RecievedCourses {
			teaModel, cmd := m.courseModel.Update(msg)
			m.courseModel, _ = teaModel.(courses.Model)
			return m, cmd
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	if m.isLoading {
		return m.sp.View()
	}
	if m.login.Err != nil {
		return m.login.View()
	}
	if m.attendance != nil {
		b, err := json.MarshalIndent(m.attendance, "", "\t")
		if err != nil {
			panic(err)
		}
		return string(b)
	}
	if m.courseModel.RecievedCourses {
		return m.courseModel.View()
	}
	return fmt.Sprintf("username: %v\n password: %v", m.login.Username, m.login.Psswd)
}
