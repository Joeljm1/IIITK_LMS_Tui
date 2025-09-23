package content

import (
	"fmt"

	"github.com/Joeljm1/IIITKlmsTui/internal/client"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CourseAttendance struct {
	Attendance    client.AllAttendance
	List          list.Model // not properly initialized till selected courseList is sent in AllAttendance
	Pos           int
	focusTable    bool
	DetailedTable table.Model
	OverAll       string
}

func (ca CourseAttendance) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, ca.List.View(), lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true).Render(ca.DetailedTable.View()), ca.OverAll)
}

func InitialCouseAttendance(width, height int) CourseAttendance {
	newTable := table.New(table.WithColumns(
		[]table.Column{
			{Title: "Date", Width: width / 9},
			{Title: "Time", Width: width / 9},
			{Title: "Status", Width: width / 9},
		}))
	newTable.SetHeight(height - height/8 - 3)
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
	return CourseAttendance{
		Attendance:    nil,
		DetailedTable: newTable,
		Pos:           0,
	}
}
func (ca CourseAttendance) Update(msg tea.Msg) (CourseAttendance, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ca.List.SetSize(msg.Width/3, msg.Height-3)
		ca.DetailedTable.SetColumns([]table.Column{
			{Title: "Date", Width: msg.Width / 9},
			{Title: "Time", Width: msg.Width / 9},
			{Title: "Status", Width: msg.Width / 9},
		})
		ca.DetailedTable.SetHeight(msg.Height - msg.Height/8 - 3) // i dont remember why -msg.height/8 but needed
		ca.List.SetHeight(msg.Height - 3)

	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			ca.focusTable = false
			ca.DetailedTable.Blur()
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
			ca.DetailedTable.SetStyles(s)
		case "right", "enter", "l":
			// if ca.tabView == tableView {
			ca.focusTable = true
			ca.DetailedTable.Focus()
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
			ca.DetailedTable.SetStyles(s)
			// }

		case "down", "j":
			if !ca.focusTable {
				ca.DetailedTable.GotoTop()
				if ca.Pos != len(ca.Attendance)-1 {
					ca.Pos = ca.Pos + 1
					var rows []table.Row
					for _, attend := range ca.Attendance[ca.Pos].Attendances {
						row := table.Row{
							attend.Date,
							attend.Time,
							attend.Status,
						}
						rows = append(rows, row)
					}
					ca.DetailedTable.SetRows(rows)
					total := ca.Attendance[ca.Pos].Overall
					overall := fmt.Sprintf("Total: %v\nPoints: %v\nPercentage:%v", total.Total, total.Points, total.Percentage)
					ca.OverAll = overall
				}
			}
		case "up", "k":
			if !ca.focusTable {
				ca.DetailedTable.GotoTop()
				if ca.Pos != 0 {
					ca.Pos = ca.Pos - 1
					var rows []table.Row
					for _, attend := range ca.Attendance[ca.Pos].Attendances {
						row := table.Row{
							attend.Date,
							attend.Time,
							attend.Status,
						}
						rows = append(rows, row)
					}
					ca.DetailedTable.SetRows(rows)

					total := ca.Attendance[ca.Pos].Overall
					overall := fmt.Sprintf("Total: %v\nPoints: %v\nPercentage:%v", total.Total, total.Points, total.Percentage)
					ca.OverAll = overall
				}
			}
		case "?":
			var cmd tea.Cmd
			list, cmd := ca.List.Update(msg)
			ca.List = list
			return ca, cmd
		}

	}

	var cmd1, cmd2 tea.Cmd
	if !ca.focusTable {
		ca.List, cmd1 = ca.List.Update(msg)
	} else {
		ca.DetailedTable, cmd2 = ca.DetailedTable.Update(msg)
	}

	return ca, tea.Batch(cmd1, cmd2)
}
