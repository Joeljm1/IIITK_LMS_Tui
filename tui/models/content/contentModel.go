package content

import (
	"fmt"

	"github.com/Joeljm1/IIITKlmsTui/internal/client"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	tabNo      int
	Attendance CourseAttendance
	// TODO: TODAY table
}

type CourseAttendance struct {
	Attendance    client.AllAttendance
	List          list.Model
	Pos           int
	focusTable    bool
	DetailedTable table.Model
	OverAll       string // MAYBENOT
}

func (m Model) Init() tea.Cmd {
	return nil
}

func InitialModel(width, height int) Model {
	newTable := table.New(table.WithColumns(
		[]table.Column{
			{Title: "Date", Width: width / 3},
			{Title: "Time", Width: width / 3},
			{Title: "Status", Width: width / 3},
		}))
	// newTable.Blur()
	return Model{
		tabNo: 0,
		Attendance: CourseAttendance{
			Attendance:    nil,
			DetailedTable: newTable,
			Pos:           0,
		},
	}
}

func (m Model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, m.Attendance.List.View(), m.Attendance.DetailedTable.View(), m.Attendance.OverAll)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Attendance.List.SetSize(msg.Width/3, msg.Height)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			m.Attendance.focusTable = false
			m.Attendance.DetailedTable.Blur()
		case "right", "enter":
			m.Attendance.focusTable = true
			m.Attendance.DetailedTable.Focus()
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
		}
	}
	var cmd1, cmd2 tea.Cmd
	if !m.Attendance.focusTable {
		m.Attendance.List, cmd1 = m.Attendance.List.Update(msg)
	} else {
		m.Attendance.DetailedTable, cmd2 = m.Attendance.DetailedTable.Update(msg)
	}
	return m, tea.Batch(cmd1, cmd2)
}
