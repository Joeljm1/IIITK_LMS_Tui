package mainModel

import tea "github.com/charmbracelet/bubbletea"

func (m model) GetDash() tea.Msg {
	dash, err := m.client.GetDashBoard()
	if err != nil {
		return err
	}
	return dash
}
