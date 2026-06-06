// internal/tui/styles.go
package tui

import "github.com/charmbracelet/lipgloss"

var (
	colorPrimary = lipgloss.Color("#06B6D4") // Cyan-500 (same as k10s)
	colorMuted   = lipgloss.Color("#6B7280") // Gray-500
	colorSuccess = lipgloss.Color("#10B981") // Emerald-500
	colorError   = lipgloss.Color("#EF4444") // Red-500
	colorWarning = lipgloss.Color("#F59E0B") // Amber-500
	colorRunning = lipgloss.Color("#3B82F6") // Blue-500
	colorNormal  = lipgloss.Color("#D1D5DB") // Gray-300

	StyleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(colorPrimary).
			Padding(0, 1)

	StyleSubHeader = lipgloss.NewStyle().
			Foreground(colorNormal)

	StyleSelected = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF"))

	StyleMuted = lipgloss.NewStyle().
			Foreground(colorMuted)

	StyleHelp = lipgloss.NewStyle().
			Foreground(colorMuted).
			Italic(true)

	// StyleHelpKey: key name in cyan bold (e.g. [enter])
	StyleHelpKey = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true)

	// StyleHelpDesc: description in muted gray
	StyleHelpDesc = lipgloss.NewStyle().
			Foreground(colorMuted)

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
