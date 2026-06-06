// internal/tui/helpers.go
package tui

import (
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// renderHelp renders key-description pairs as "[key] desc   [key2] desc2"
// with keys in cyan bold and descriptions in muted gray.
// Usage: renderHelp("enter", "select", "esc", "back", "n", "new pipeline")
func renderHelp(pairs ...string) string {
	var parts []string
	for i := 0; i < len(pairs); i += 2 {
		k := StyleHelpKey.Render("[" + pairs[i] + "]")
		var desc string
		if i+1 < len(pairs) {
			desc = StyleHelpDesc.Render(pairs[i+1])
		}
		parts = append(parts, k+" "+desc)
	}
	return "  " + strings.Join(parts, "   ")
}

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
