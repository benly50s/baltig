// internal/tui/keys.go
package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all key bindings for baltig.
type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Back   key.Binding
	New    key.Binding
	Add    key.Binding
	Delete key.Binding
	Open   key.Binding
	Follow key.Binding
	Help   key.Binding
	Quit   key.Binding
}

// Keys is the global key map used by all TUI screens.
var Keys = KeyMap{
	Up:     key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:   key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Select: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	Back:   key.NewBinding(key.WithKeys("esc", "b"), key.WithHelp("esc/b", "back")),
	New:    key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new pipeline")),
	Add:    key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "add repo")),
	Delete: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
	Open:   key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open in browser")),
	Follow: key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "follow log")),
	Help:   key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	Quit:   key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("q", "quit")),
}
