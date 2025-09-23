package content

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	tabView    views
	Attendance CourseAttendance
	Today      TodayAttendance
	DashBoard  DashBoard
	height     int
	width      int
}

type views int

const (
	todayView views = iota + 1
	tableView
	dashView
)

func (m Model) Init() tea.Cmd {
	return nil
}

func InitialModel(width, height int) Model {
	// newTable.Blur()
	return Model{
		tabView:    1,
		Attendance: InitialCouseAttendance(width, height),
		Today:      InitialTodayAttendance(width, height),
		height:     height,
		width:      width,
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			m.tabView = todayView
		case "2":
			m.tabView = tableView
		case "3":
			m.tabView = dashView
		case "tab":
			m.tabView = m.tabView + 1
			if m.tabView > 3 {
				m.tabView = todayView
			}
		}
	}
	var cmd1, cmd2 tea.Cmd
	switch m.tabView {
	case tableView:
		m.Attendance, cmd1 = m.Attendance.Update(msg)
	case todayView:
		m.Today, cmd2 = m.Today.Update(msg)
	}
	return m, tea.Batch(cmd1, cmd2)
}

func (m Model) View() string {
	// switch m.tabView {
	// case todayView:
	// 	topBar := lipgloss.JoinHorizontal(lipgloss.Top, selected1Bar, unSelected2Bar, remainingBarStyle(m.tabView, m.width-24))
	// 	return lipgloss.JoinVertical(lipgloss.Top, topBar, m.Today.View())
	// case tableView:
	// 	topBar := lipgloss.JoinHorizontal(lipgloss.Top, unSelected1Bar, selected2Bar, remainingBarStyle(m.tabView, m.width-24))
	// 	return lipgloss.JoinVertical(lipgloss.Top, topBar, m.Attendance.View())
	// 	// return topBar + m.Attendance.View()
	// 	// return "\n\n" + topBar + m.Attendance.View()
	// case dashView:
	// 	topBar := lipgloss.JoinHorizontal(lipgloss.Top, unSelected1Bar, selected2Bar, remainingBarStyle(m.tabView, m.width-24))
	// 	return lipgloss.JoinVertical(lipgloss.Top, topBar, m.DashBoard.View())
	// default:
	// 	panic("Unexpected tab no")
	// }
	topBar := m.topBar()
	switch m.tabView {
	case todayView:
		return lipgloss.JoinVertical(lipgloss.Top, topBar, m.Today.View())
	case tableView:
		return lipgloss.JoinVertical(lipgloss.Top, topBar, m.Attendance.View())
	case dashView:
		return lipgloss.JoinVertical(lipgloss.Top, topBar, m.DashBoard.View())
	default:
		return "Error unexpected tab no"
	}
}
