package content

import (
	"strings"
	"time"

	"github.com/Joeljm1/IIITKlmsTui/internal/client"
	"github.com/charmbracelet/lipgloss"
)

type DashBoard struct {
	DashBoard *client.Dashboard
	// TODO: for the view
}

// View generates a beautiful dashboard display using lipgloss
func (d DashBoard) View() string {
	if d.DashBoard == nil {
		noEventsStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(2, 4).
			Foreground(lipgloss.Color("244")).
			Align(lipgloss.Center)
		return noEventsStyle.Render("No events available")
	}

	if d.DashBoard.Error {
		errorStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("196")).
			Padding(1, 2).
			Foreground(lipgloss.Color("196")).
			Bold(true)
		return errorStyle.Render("Error loading dashboard")
	}

	if len(d.DashBoard.Data.Events) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(2, 4).
			Foreground(lipgloss.Color("244")).
			Align(lipgloss.Center)
		return emptyStyle.Render("No upcoming events")
	}

	// Header style
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		Padding(0, 1).
		MarginBottom(1)

	// Event card styles
	eventCardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(1, 2).
		MarginBottom(1).
		Width(60)

	overdueCardStyle := eventCardStyle.Copy().
		BorderForeground(lipgloss.Color("196")).
		Foreground(lipgloss.Color("196"))

	// Event content styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("229"))

	courseStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("117")).
		Italic(true)

	purposeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("244"))

	timeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214"))

	urlStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("75")).
		Underline(true)

	overdueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	var eventCards []string
	header := headerStyle.Render("📅 Upcoming Events")

	for _, event := range d.DashBoard.Data.Events {
		// Format timestamp
		timeStr := time.Unix(int64(event.Timestart), 0).Format("Jan 02, 2006 15:04")

		// Build event content
		var lines []string
		lines = append(lines, titleStyle.Render(event.Name))
		lines = append(lines, courseStyle.Render("📚 "+event.Course.Fullname))

		if event.Purpose != "" {
			lines = append(lines, purposeStyle.Render("Purpose: "+event.Purpose))
		}

		timeLineContent := timeStyle.Render("🕐 " + timeStr)
		if event.Overdue {
			timeLineContent += "  " + overdueStyle.Render("OVERDUE")
		}
		lines = append(lines, timeLineContent)

		if event.URL != "" {
			lines = append(lines, urlStyle.Render("🔗 "+event.URL))
		}

		eventContent := strings.Join(lines, "\n")

		// Apply appropriate card style
		var card string
		if event.Overdue {
			card = overdueCardStyle.Render(eventContent)
		} else {
			card = eventCardStyle.Render(eventContent)
		}

		eventCards = append(eventCards, card)
	}

	// Join all cards vertically
	allContent := strings.Join(eventCards, "\n")

	return lipgloss.JoinVertical(lipgloss.Left, header, allContent)
}
