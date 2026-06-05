// internal/tui/helpers.go
package tui

import (
	"os/exec"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
)

// openBrowser opens a URL in the system default browser.
func openBrowser(url string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("open", url)
		case "linux":
			cmd = exec.Command("xdg-open", url)
		default:
			cmd = exec.Command("cmd", "/c", "start", url)
		}
		_ = cmd.Start()
		return nil
	}
}
