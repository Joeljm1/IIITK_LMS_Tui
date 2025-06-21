package courses

import (
	"errors"

	"github.com/Joeljm1/IIITKlmsTui/internal/client"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	RecievedCourses bool
	Chosen          []client.Course
	List            list.Model
}

func InitialModel() Model {
	return Model{
		RecievedCourses: false,
		Chosen:          nil,
		// FIX: If nill pointer err may be cause List not initialized here but only when course is recieved so some fn it calling it before it is required
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

type Load struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "space":
			selectedItem := m.List.SelectedItem()
			if selectedItem == nil {
				return m, func() tea.Msg { return errors.New("selected Item is nil") }
			}
			selectedCourse := selectedItem.(CourseItem)
			if _, ok := Selected[selectedCourse.Id]; ok {
				delete(Selected, selectedCourse.Id)
			} else {
				Selected[selectedCourse.Id] = struct{}{}
			}
		case "ctrl+s", "ctrl+v":
			chosen := []client.Course{}
			for _, item := range m.List.Items() {
				course := item.(CourseItem)
				if _, ok := Selected[course.Id]; ok {
					chosen = append(chosen, client.Course(course))
				}
			}
			m.Chosen = chosen
			return m, load
		}
	case tea.WindowSizeMsg:
		h, v := DocStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	m.List.Title = "Select which courses you want to follow"
	return DocStyle.Render(m.List.View())
}
