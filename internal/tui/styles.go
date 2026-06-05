// internal/tui/styles.go
package tui

import "github.com/charmbracelet/lipgloss"

var (
	colorPrimary = lipgloss.Color("69")  // indigo
	colorMuted   = lipgloss.Color("240") // gray
	colorSuccess = lipgloss.Color("76")  // green
	colorError   = lipgloss.Color("196") // red
	colorWarning = lipgloss.Color("214") // orange
	colorRunning = lipgloss.Color("33")  // blue

	StyleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary)

	StyleSubHeader = lipgloss.NewStyle().
			Foreground(colorMuted)

	StyleSelected = lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("255"))

	StyleMuted = lipgloss.NewStyle().
			Foreground(colorMuted)

	StyleHelp = lipgloss.NewStyle().
			Foreground(colorMuted).
			MarginTop(1)

	StyleSuccess = lipgloss.NewStyle().Foreground(colorSuccess)
	StyleError   = lipgloss.NewStyle().Foreground(colorError)
	StyleWarning = lipgloss.NewStyle().Foreground(colorWarning)
	StyleRunning = lipgloss.NewStyle().Foreground(colorRunning)
)

// StatusIcon returns a colored icon string for a pipeline/job status.
func StatusIcon(status string) string {
	switch status {
	case "success":
		return StyleSuccess.Render("✓")
	case "failed":
		return StyleError.Render("✗")
	case "running":
		return StyleRunning.Render("●")
	case "pending":
		return StyleWarning.Render("○")
	case "canceled":
		return StyleMuted.Render("⊘")
	default:
		return StyleMuted.Render("?")
	}
}
