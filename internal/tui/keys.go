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
	Run    key.Binding
	Delete key.Binding
	Open   key.Binding
	Follow key.Binding
	Help   key.Binding
	Quit   key.Binding
}

// Keys is the global key map used by all TUI screens.
var Keys = KeyMap{
	Up:     key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "이동")),
	Down:   key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "이동")),
	Select: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "선택")),
	Back:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "뒤로")),
	New:    key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "새 파이프라인")),
	Run:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "실행")),
	Delete: key.NewBinding(key.WithKeys("ctrl+d"), key.WithHelp("ctrl+d", "삭제")),
	Open:   key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "브라우저")),
	Follow: key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "follow")),
	Help:   key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "도움말")),
	Quit:   key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("q", "종료")),
}
