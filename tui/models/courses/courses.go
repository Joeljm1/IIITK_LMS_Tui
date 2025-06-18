package courses

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Joeljm1/IIITKlmsTui/internal/client"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	SubmitKey = key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "Submit Courses"))
	SelectKey = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Select Course"))
	LogoutKey = key.NewBinding(key.WithKeys("ctrl+d"), key.WithHelp("ctrl+d", "Logout"))
)

var (
	DocStyle = lipgloss.NewStyle().Margin(1, 2)
	Selected = make(map[int]struct{})

	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#c20000"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	blurredButton  = blurredStyle.Render(" [ Submit] ")
	focussedButton = focusedStyle.Foreground(lipgloss.Color("#17bd05")).Render(" [ Submit] ")
)

type ChoiceFromLMSNeeded struct{}

func CheckChoiceFile(lms *client.LMSCLient) tea.Cmd {
	return func() tea.Msg {
		choice, err := lms.ChoicesInFile()
		if err != nil {
			return ChoiceFromLMSNeeded{}
		}
		return choice
	}
}

func GetCoursesFromLms(lms *client.LMSCLient) tea.Cmd {
	return func() tea.Msg {
		courses, err := lms.GetCoursesFromLMS()
		if err != nil {
			return err
		}
		return courses
	}
}

func load() tea.Msg {
	return Load(struct{}{})
}

func GetAttendanceList(courseList []client.Course, lms *client.LMSCLient) tea.Cmd {
	return func() tea.Msg {
		choice, err := lms.AttendanceForAllSelection(courseList) // CourseList is []client.Course
		if err != nil {
			return err
		}
		b, err := json.Marshal(choice)
		if err != nil {
			return err
		}
		os.WriteFile(client.ChoiceFile, b, 0644)
		return choice
	}
}

func GetAllAttendance(lms *client.LMSCLient) tea.Cmd { // TODO: test if putting attendance detail outside the return matters
	return func() tea.Msg {
		attend, err := lms.AllAttendance()
		if err != nil {
			return err
		}
		return attend
	}
}

type CourseItem client.Course

func (item CourseItem) Title() string {
	sn := strings.Trim(item.Shortname, " ")
	isSelected := " "
	if _, ok := Selected[item.Id]; ok {
		isSelected = "✅"
	}
	return fmt.Sprintf("[ %v ] %v", isSelected, sn)
}
func (item CourseItem) Description() string { return item.Fullname }
func (item CourseItem) FilterValue() string {
	if strings.Trim(item.Fullname, " ") == "" {
		panic("WTF")
	}
	return strings.ToLower(item.Fullname)
}
