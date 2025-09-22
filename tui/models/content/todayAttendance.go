package content

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: Handle resizing
type TodayAttendance struct {
	Table table.Model
}

func InitialTodayAttendance(width, height int) TodayAttendance {
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
	return TodayAttendance{
		Table: todayTable,
	}
}

func (ta TodayAttendance) View() string {
	return ta.Table.View()
}

// TODO: checj if height is updated to
func (ta TodayAttendance) Update(msg tea.Msg) (TodayAttendance, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cols := []table.Column{
			{
				Title: "Course",
				Width: msg.Width / 3,
			},
			{
				Title: "Time",
				Width: msg.Width / 3,
			},
			{
				Title: "Status",
				Width: msg.Width / 3,
			},
		}
		ta.Table.SetColumns(cols)
		ta.Table.SetHeight(msg.Height - msg.Height/8 - 3)
	}
	ta.Table, cmd = ta.Table.Update(msg)
	return ta, cmd
}
