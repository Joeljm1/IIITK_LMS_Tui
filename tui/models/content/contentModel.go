package content

import (
	"fmt"
	"log"

	"github.com/Joeljm1/IIITKlmsTui/internal/client"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	tabNo      int
	Attendance CourseAttendance
	Today      TodayAttendance
	// TODO: TODAY table
}

type CourseAttendance struct {
	Attendance    client.AllAttendance
	List          list.Model
	Pos           int
	focusTable    bool
	DetailedTable table.Model
	OverAll       string
}

type TodayAttendance struct {
	Table table.Model
}

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
	newTable.SetHeight(height - height/8)
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
		Background(lipgloss.Color("57")).
		Bold(false)
	newTable.SetStyles(s)
	todayTable.SetStyles(s)
	todayTable.Focus()
	// newTable.Blur()
	return Model{
		tabNo: 1,
		Attendance: CourseAttendance{
			Attendance:    nil,
			DetailedTable: newTable,
			Pos:           0,
		},
		Today: TodayAttendance{
			Table: todayTable,
		},
	}
}

// TODO: Split update of today and details
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Attendance.List.SetSize(msg.Width/3, msg.Height)

		m.Attendance.DetailedTable.SetColumns([]table.Column{
			{Title: "Date", Width: msg.Width / 9},
			{Title: "Time", Width: msg.Width / 9},
			{Title: "Status", Width: msg.Width / 9},
		})
		m.Attendance.DetailedTable.SetHeight(msg.Height - msg.Height/8)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			m.Attendance.focusTable = false
			m.Attendance.DetailedTable.Blur()
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
			m.Attendance.DetailedTable.SetStyles(s)
		case "right", "enter", "l":
			if m.tabNo == 2 {
				m.Attendance.focusTable = true
				m.Attendance.DetailedTable.Focus()
				s := table.DefaultStyles()
				s.Header = s.Header.
					BorderStyle(lipgloss.NormalBorder()).
					BorderForeground(lipgloss.Color("240")).
					BorderBottom(true).
					Bold(false)
				s.Selected = s.Selected.
					Foreground(lipgloss.Color("229")).
					Background(lipgloss.Color("57")).
					Bold(false)
				m.Attendance.DetailedTable.SetStyles(s)
			}
		case "down", "j":
			if !m.Attendance.focusTable {
				if m.Attendance.Pos != len(m.Attendance.Attendance)-1 {
					m.Attendance.Pos = m.Attendance.Pos + 1
					var rows []table.Row
					for _, attend := range m.Attendance.Attendance[m.Attendance.Pos].Attendances {
						row := table.Row{
							attend.Date,
							attend.Time,
							attend.Status,
						}
						rows = append(rows, row)
					}
					m.Attendance.DetailedTable.SetRows(rows)
					total := m.Attendance.Attendance[m.Attendance.Pos].Overall
					overall := fmt.Sprintf("Total: %v\nPoints: %v\nPercentage:%v", total.Total, total.Points, total.Percentage)
					m.Attendance.OverAll = overall
				}
			}
		case "up", "k":
			if !m.Attendance.focusTable {
				if m.Attendance.Pos != 0 {
					m.Attendance.Pos = m.Attendance.Pos - 1
					var rows []table.Row
					for _, attend := range m.Attendance.Attendance[m.Attendance.Pos].Attendances {
						row := table.Row{
							attend.Date,
							attend.Time,
							attend.Status,
						}
						rows = append(rows, row)
					}
					m.Attendance.DetailedTable.SetRows(rows)

					total := m.Attendance.Attendance[m.Attendance.Pos].Overall
					overall := fmt.Sprintf("Total: %v\nPoints: %v\nPercentage:%v", total.Total, total.Points, total.Percentage)
					m.Attendance.OverAll = overall
				}
			}
		case "1":
			m.tabNo = 1
		case "2":
			m.tabNo = 2
		}
	}
	var cmd1, cmd2, cmd3 tea.Cmd
	if m.tabNo == 2 {
		if !m.Attendance.focusTable {
			m.Attendance.List, cmd1 = m.Attendance.List.Update(msg)
		} else {
			m.Attendance.DetailedTable, cmd2 = m.Attendance.DetailedTable.Update(msg)
		}
	} else if m.tabNo == 1 {
		log.Println("Accept")
		m.Today.Table, cmd3 = m.Today.Table.Update(msg)
	}
	return m, tea.Batch(cmd1, cmd2, cmd3)
}

func (ca CourseAttendance) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, ca.List.View(), lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true).Render(ca.DetailedTable.View()), ca.OverAll)
}

func (ta TodayAttendance) View() string {
	return ta.Table.View()
}

func (m Model) View() string {
	switch m.tabNo {
	case 1:
		return m.Today.View()
	case 2:
		return m.Attendance.View()
	default:
		panic("Unexpected tab no")
	}
}
