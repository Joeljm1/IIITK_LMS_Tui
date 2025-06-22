package mainModel

import (
	"fmt"
	"time"

	"github.com/Joeljm1/IIITKlmsTui/internal/client"
	"github.com/Joeljm1/IIITKlmsTui/tui/models/content"
	"github.com/Joeljm1/IIITKlmsTui/tui/models/courses"
	"github.com/Joeljm1/IIITKlmsTui/tui/models/login"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
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
	contentModel  content.Model
	err           error
	sentWidth     bool // to give width to table
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
		sentWidth:   false,
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
		if !m.sentWidth {
			m.contentModel = content.InitialModel(m.width/3, m.height)
			m.sentWidth = true
		}
		if m.contentModel.Attendance.Attendance != nil {
			contentModel, cmd := m.contentModel.Update(msg)
			m.contentModel = contentModel.(content.Model)
			return m, cmd
		} else if !m.courseModel.RecievedCourses {
			teaModel, cmd := m.login.Update(msg)
			m.login = teaModel.(login.Model)
			return m, cmd
		} else if m.courseModel.RecievedCourses {
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
		m.client.Choices = msg // make pointer check if error later
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

		m.isLoading = false
	case client.AllAttendance:
		m.isLoading = false
		m.contentModel.Attendance.Attendance = msg
		var l []list.Item
		var rows []table.Row
		overall := fmt.Sprintf("Total: %v\nPoints: %v\nPercentage:%v", msg[0].Overall.Total, msg[0].Overall.Points, msg[0].Overall.Percentage)
		for _, attend := range msg[0].Attendances {
			row := table.Row{
				attend.Date,
				attend.Time,
				attend.Status,
			}
			rows = append(rows, row)
		}
		m.contentModel.Attendance.OverAll = overall
		m.contentModel.Attendance.DetailedTable.SetRows(rows)
		for _, attendDet := range msg {
			l = append(l, content.CourseNameItem(attendDet.CourseName))
		}
		m.contentModel.Attendance.List = list.New(l, list.NewDefaultDelegate(), m.width/3, m.height)

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
		if m.contentModel.Attendance.Attendance != nil {
			teaModel, cmd := m.contentModel.Update(msg)
			m.contentModel, _ = teaModel.(content.Model)
			return m, cmd

		}
		if m.courseModel.RecievedCourses {
			teaModel, cmd := m.courseModel.Update(msg)
			m.courseModel, _ = teaModel.(courses.Model)
			return m, cmd
		}

	default: // any other msg from other models passes here
		cmds := make([]tea.Cmd, 2)
		var teaModel tea.Model
		if m.login.Err != nil {
			teaModel, cmds[0] = m.login.Update(msg)
			m.login, _ = teaModel.(login.Model)
		}
		if m.courseModel.RecievedCourses {
			teaModel, cmds[1] = m.courseModel.Update(msg)
			m.courseModel, _ = teaModel.(courses.Model)
		}
		return m, tea.Batch(cmds...)
	}
	return m, nil
}

func (m model) View() string {
	errStyle := lipgloss.NewStyle().Width(m.width).Height(m.height).Align(lipgloss.Center).Padding(m.height/2, 0).Foreground(lipgloss.Color("#f70000"))
	if m.err != nil {
		return errStyle.Render(m.err.Error())
	}
	if m.isLoading {
		return m.sp.View()
	}
	if m.login.Err != nil {
		return m.login.View()
	}
	if m.contentModel.Attendance.Attendance != nil {
		return m.contentModel.View()
	}
	if m.courseModel.RecievedCourses {
		return m.courseModel.View()
	}
	return errStyle.Render("An Error has occured.\nPlease press ctrl+d to logout and reset")
}
