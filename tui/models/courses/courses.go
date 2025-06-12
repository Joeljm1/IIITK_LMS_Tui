package courses

import (
	"encoding/json"
	"os"

	"github.com/Joeljm1/IIITKlmsTui/internal/client"
	tea "github.com/charmbracelet/bubbletea"
)

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
