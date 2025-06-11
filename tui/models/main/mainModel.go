package mainModel

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Joeljm1/IIITKlmsTui/internal/client"
	"github.com/Joeljm1/IIITKlmsTui/tui/models/courses"
	"github.com/Joeljm1/IIITKlmsTui/tui/models/login"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	sp          spinner.Model
	login       login.Model
	isLoading   bool
	client      *client.LMSCLient
	courseModel courses.Model
	err         error
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
	case login.UName, login.Psswd, login.LoginValidErr, login.Valid, tea.WindowSizeMsg:
		// m.login,cmd=m.Update(msg) do it with type inference
		teaModel, cmd := m.login.Update(msg)
		m.login, _ = teaModel.(login.Model)
		return m, cmd
	case login.LoginErr:
		teaModel, cmd := m.login.Update(msg)
		m.login, _ = teaModel.(login.Model)
		m.isLoading = false
		return m, cmd
	case login.LoginComplete:
		m.login.Err = nil
	case login.Load, courses.Load:
		m.isLoading = true
		return m, courses.GetAttendanceList(m.courseModel.Chosen, m.client)
	case *client.LMSCLient:
		log.Println("Got client")
		m.client = msg
		return m, courses.CheckChoiceFile(m.client)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.sp, cmd = m.sp.Update(msg)
		return m, cmd

	case courses.ChoiceFromLMSNeeded:
		return m, courses.GetCoursesFromLms(m.client)
	case client.Choices:
		m.isLoading = false
		m.client.Choices = msg
		m.client.RecivedChoices = true
		m.isLoading = false
	// update view and check attendance
	case client.CourseList:
		m.courseModel.AllCourses = msg
		m.isLoading = false
		// update view to ask

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
		if m.login.Err != nil {
			teaModel, cmd := m.login.Update(msg)
			m.login, _ = teaModel.(login.Model)
			return m, cmd
		}
		if m.courseModel.AllCourses != nil {
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
	if m.client.RecivedChoices {
		b, err := json.Marshal(m.client.Choices)
		if err != nil {
			panic(err)
		}
		return string(b)
	}
	if m.courseModel.AllCourses != nil {
		return m.courseModel.View()
	}
	return fmt.Sprintf("username: %v\n password: %v", m.login.Username, m.login.Psswd)
}
