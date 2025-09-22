package content

import (
	"github.com/Joeljm1/IIITKlmsTui/internal/client"
	"github.com/charmbracelet/bubbles/table"
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

type DashBoard struct {
	DashBoard *client.Dashboard
	//TODO: for the view
}

type views int

const (
	todayView views = iota + 1
	tableView
	dashView
)

func unselectedTodayBorder() lipgloss.Border {
	b := lipgloss.NormalBorder()
	b.BottomRight = "┘"
	b.TopRight = "┬"
	return b
}

func selectedTodayBorder() lipgloss.Border {
	b := lipgloss.NormalBorder()
	b.Bottom = " "
	b.BottomLeft = "╵"
	b.BottomRight = "└"
	b.TopRight = "┬"
	return b
}

func remainingBarStyle(selected views, width int) string {
	b := lipgloss.NormalBorder()
	b.TopLeft = "┬"
	switch selected {
	case todayView:
		b.BottomLeft = "┴"
	case tableView:
		b.BottomLeft = "└"
	}
	s := lipgloss.NewStyle().Width(width).Border(b, true).Render("")
	return s
}

var (
	todaySelectedBorder   = selectedTodayBorder()
	todayUnselectedBorder = unselectedTodayBorder()
	unSelectedTodayStyle  = lipgloss.NewStyle().Border(todayUnselectedBorder, true).Width(10).Align(lipgloss.Center)
	selectedTodayStyle    = unSelectedTodayStyle.Border(todaySelectedBorder, true)
	selectedTodayBar      = selectedTodayStyle.Render("1)Today")
	unSelectedTodayBar    = unSelectedTodayStyle.Render("1)Today")
	selectedDescStyle     = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true, false, false, false).Width(10)
	unSelectedDescStyle   = selectedDescStyle.Border(lipgloss.NormalBorder(), true, false, true, false)
	selectedDescBar       = selectedDescStyle.Render("2)Details")
	unSelectedDescBar     = unSelectedDescStyle.Render("2)Details")
)

func (m Model) Init() tea.Cmd {
	return nil
}

func InitialModel(width, height int) Model {
	newTable := table.New(table.WithColumns(
		[]table.Column{
			{Title: "Date", Width: width / 9},
			{Title: "Time", Width: width / 9},
			{Title: "Status", Width: width / 9},
		}))
	newTable.SetHeight(height - height/8 - 3)
	todayTable := table.New(table.WithColumns([]table.Column{
		{
			Title: "Course",
			Width: width / 3,
		},
		{
			Title: "Time",
			Width: width / 3,
		},
		{
			Title: "Status",
			Width: width / 3,
		},
	}))
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("0")).
		Bold(false)
	newTable.SetStyles(s)
	s.Selected.Background(lipgloss.Color("57"))
	todayTable.SetStyles(s)
	todayTable.Focus()
	// newTable.Blur()
	return Model{
		tabView: 1,
		Attendance: CourseAttendance{
			Attendance:    nil,
			DetailedTable: newTable,
			Pos:           0,
		},
		Today: TodayAttendance{
			Table: todayTable,
		},
		height: height,
		width:  width,
	}
}

// TODO: Split update of today and details
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			m.tabView = todayView
		case "2":
			m.tabView = tableView
		case "tab":
			m.tabView = m.tabView + 1
			if m.tabView > 2 {
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
	switch m.tabView {
	case todayView:
		topBar := lipgloss.JoinHorizontal(lipgloss.Top, selectedTodayBar, unSelectedDescBar, remainingBarStyle(m.tabView, m.width-24))
		return lipgloss.JoinVertical(lipgloss.Top, topBar, m.Today.View())
	case tableView:
		topBar := lipgloss.JoinHorizontal(lipgloss.Top, unSelectedTodayBar, selectedDescBar, remainingBarStyle(m.tabView, m.width-24))
		return lipgloss.JoinVertical(lipgloss.Top, topBar, m.Attendance.View())
		// return topBar + m.Attendance.View()
		// return "\n\n" + topBar + m.Attendance.View()
	default:
		panic("Unexpected tab no")
	}
}
