package courses

import (
	"fmt"
	"strings"

	"github.com/Joeljm1/IIITKlmsTui/internal/client"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	AllCourses client.CourseList //[]Course
	Selected   map[int]struct{}
	focus      int
	Chosen     []client.Course
}

func InitialModel() Model {
	return Model{
		AllCourses: nil,
		Selected:   make(map[int]struct{}),
		focus:      0,
		Chosen:     nil,
	}
}

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#c20000"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	blurredButton  = blurredStyle.Render(" [ Submit] ")
	focussedButton = focusedStyle.Foreground(lipgloss.Color("#17bd05")).Render(" [ Submit] ")
)

type ChoiceFromLMSNeeded struct{}

func (m Model) Init() tea.Cmd {
	return nil
}

type Load struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.focus = max(m.focus-1, 0)
		case "down", "j":
			m.focus = min(m.focus+1, len(m.AllCourses))
		case "enter", "space":
			if m.focus == len(m.AllCourses) {
				// TODO:  Handle selected values
				chosen := []client.Course{}
				for i := range m.Selected {
					chosen = append(chosen, m.AllCourses[i])
				}
				m.Chosen = chosen
				return m, load // TODO: also do writing choice,cheching choice from file and fetching attendance
			} else {
				if _, ok := m.Selected[m.focus]; ok {
					delete(m.Selected, m.focus)
				} else {
					m.Selected[m.focus] = struct{}{}
				}
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	sb := strings.Builder{}
	sb.WriteString("Select which courses you want to follow:\n")
	for i := range m.AllCourses {
		var courseName string
		if m.focus == i {
			courseName = focusedStyle.Render(m.AllCourses[i].Fullname)
		} else {
			courseName = m.AllCourses[i].Fullname
		}
		selected := " "
		if _, ok := m.Selected[i]; ok {
			selected = "✅"
		}
		sb.WriteString(fmt.Sprintf("[%v] %v) %v\n", selected, i+1, courseName))
	}
	button := blurredButton
	if m.focus == len(m.AllCourses) {
		button = focussedButton
	}
	sb.WriteString("\n" + button)
	return sb.String()
}
