package content

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type TodayAttendance struct {
	Table table.Model
}

func (ta TodayAttendance) View() string {
	return ta.Table.View()
}

func (ta TodayAttendance) Update(msg tea.Msg) (TodayAttendance, tea.Cmd) {
	var cmd tea.Cmd
	ta.Table, cmd = ta.Table.Update(msg)
	return ta, cmd
}
