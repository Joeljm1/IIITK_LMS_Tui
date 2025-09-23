// TODO: top bar small prob
package content

import (
	"github.com/charmbracelet/lipgloss"
)

func unselectedLeftBorder() lipgloss.Border {
	b := lipgloss.NormalBorder()
	b.BottomRight = "┘"
	b.TopRight = "┬"
	return b
}

func unselectedMidBorder() lipgloss.Border {
	b := lipgloss.NormalBorder()
	b.BottomRight = "┴"
	b.BottomLeft = "┴"
	b.TopRight = "┬"
	b.TopLeft = "┬"
	return b
}

func selectedLeftBorder() lipgloss.Border {
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
	// last val selected
	if selected == dashView {
		b.BottomLeft = "└"
	} else {
		b.BottomLeft = "┴"
	}
	s := lipgloss.NewStyle().Width(width).Border(b, true).Render("")
	return s
}

var (
	leftSelectedBorder   = selectedLeftBorder()
	leftUnselectedBorder = unselectedLeftBorder()
	unSelectedLeftStyle  = lipgloss.NewStyle().Border(leftUnselectedBorder, true).Width(10).Align(lipgloss.Center)
	selectedLeftStyle    = unSelectedLeftStyle.Border(leftSelectedBorder, true)
	selectedRightStyle   = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true, false, false, false).Width(10)
	unSelectedRightStyle = selectedRightStyle.Border(lipgloss.NormalBorder(), true, false, true, false)
	selectedMidStyle     = lipgloss.NewStyle().Border(leftSelectedBorder, true, true, true, false).Width(10)
	unSelectedMidStyle   = selectedRightStyle.Border(leftUnselectedBorder, true, true, true, false)
	selected1Bar         = selectedLeftStyle.Render("1)Today")
	unSelected1Bar       = unSelectedLeftStyle.Render("1)Today")
	selected2Bar         = selectedMidStyle.Render("2)Details")
	unSelected2Bar       = unSelectedMidStyle.Render("2)Details")
	selected3Bar         = selectedRightStyle.Render("3)Dash")
	unSelected3Bar       = unSelectedRightStyle.Render("3)Dash")
	boxLen               = 12
)

func (m Model) topBar() string {
	var topBar string
	switch m.tabView {
	case todayView:
		topBar = lipgloss.JoinHorizontal(lipgloss.Top, selected1Bar, unSelected2Bar, unSelected3Bar, remainingBarStyle(m.tabView, m.width-boxLen*3))
	case tableView:
		topBar = lipgloss.JoinHorizontal(lipgloss.Top, unSelected1Bar, selected2Bar, unSelected3Bar, remainingBarStyle(m.tabView, m.width-boxLen*3))
	case dashView:
		topBar = lipgloss.JoinHorizontal(lipgloss.Top, unSelected1Bar, unSelected2Bar, selected3Bar, remainingBarStyle(m.tabView, m.width-boxLen*3))
	default:
		panic("Unexpected tab no")
	}
	return topBar
}
