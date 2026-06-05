// internal/tui/pipeline_run.go
package tui

import (
	"fmt"

	"github.com/benly/baltig/internal/config"
	"github.com/benly/baltig/internal/gitlab"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type pipelineCreatedMsg struct{ p gitlab.Pipeline }
type pipelineCreateErrMsg struct{ err error }

// varField holds a key-value pair for a pipeline variable input row.
type varField struct {
	keyInput   textinput.Model
	valueInput textinput.Model
}

func newVarField() varField {
	k := textinput.New()
	k.Placeholder = "KEY"
	k.Width = 15

	v := textinput.New()
	v.Placeholder = "value"
	v.Width = 25

	return varField{keyInput: k, valueInput: v}
}

type PipelineRunModel struct {
	cfg        *config.Config
	client     *gitlab.Client
	project    config.ProjectEntry
	refInput   textinput.Model
	vars       []varField
	focusIdx   int // 0=ref, 1=var[0].key, 2=var[0].value, 3=var[1].key, ...
	submitting bool
	status     string
}

func NewPipelineRunModel(cfg *config.Config, client *gitlab.Client, project config.ProjectEntry) *PipelineRunModel {
	ref := textinput.New()
	ref.Placeholder = "브랜치 또는 태그"
	ref.SetValue(cfg.Global.DefaultRef)
	ref.Width = 30

	return &PipelineRunModel{
		cfg:      cfg,
		client:   client,
		project:  project,
		refInput: ref,
		vars:     []varField{newVarField()},
		focusIdx: 0,
	}
}

func (m *PipelineRunModel) totalFields() int {
	return 1 + len(m.vars)*2
}

func (m *PipelineRunModel) updateFocus() tea.Cmd {
	m.refInput.Blur()
	for i := range m.vars {
		m.vars[i].keyInput.Blur()
		m.vars[i].valueInput.Blur()
	}

	switch {
	case m.focusIdx == 0:
		return m.refInput.Focus()
	default:
		varIdx := (m.focusIdx - 1) / 2
		isKey := (m.focusIdx-1)%2 == 0
		if varIdx < len(m.vars) {
			if isKey {
				return m.vars[varIdx].keyInput.Focus()
			}
			return m.vars[varIdx].valueInput.Focus()
		}
	}
	return nil
}

func (m *PipelineRunModel) Init() tea.Cmd {
	return m.refInput.Focus() // focus ref input on start
}

func (m *PipelineRunModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case pipelineCreatedMsg:
		m.submitting = false
		return m, func() tea.Msg { return navigateBack{} }

	case pipelineCreateErrMsg:
		m.submitting = false
		m.status = StyleError.Render("오류: " + msg.err.Error())

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Back):
			return m, func() tea.Msg { return navigateBack{} }

		case msg.String() == "ctrl+r" && !m.submitting:
			return m, m.submit()

		case msg.String() == "tab":
			m.focusIdx = (m.focusIdx + 1) % m.totalFields()
			cmd := m.updateFocus()
			return m, cmd

		case msg.String() == "shift+tab":
			m.focusIdx = (m.focusIdx - 1 + m.totalFields()) % m.totalFields()
			cmd := m.updateFocus()
			return m, cmd

		case msg.String() == "ctrl+a":
			m.vars = append(m.vars, newVarField())
			return m, nil
		}
	}

	// Update the focused input
	var cmd tea.Cmd
	switch {
	case m.focusIdx == 0:
		m.refInput, cmd = m.refInput.Update(msg)
	default:
		varIdx := (m.focusIdx - 1) / 2
		isKey := (m.focusIdx-1)%2 == 0
		if varIdx < len(m.vars) {
			if isKey {
				m.vars[varIdx].keyInput, cmd = m.vars[varIdx].keyInput.Update(msg)
			} else {
				m.vars[varIdx].valueInput, cmd = m.vars[varIdx].valueInput.Update(msg)
			}
		}
	}
	return m, cmd
}

func (m *PipelineRunModel) submit() tea.Cmd {
	m.submitting = true
	m.status = StyleMuted.Render("파이프라인 생성 중...")

	ref := m.refInput.Value()
	vars := make([]gitlab.PipelineVariable, 0, len(m.vars))
	for _, vf := range m.vars {
		k := vf.keyInput.Value()
		v := vf.valueInput.Value()
		if k != "" {
			vars = append(vars, gitlab.PipelineVariable{Key: k, Value: v})
		}
	}

	client := m.client
	projectID := int64(m.project.ID)

	return func() tea.Msg {
		p, err := client.CreatePipeline(projectID, ref, vars)
		if err != nil {
			return pipelineCreateErrMsg{err}
		}
		return pipelineCreatedMsg{*p}
	}
}

func (m *PipelineRunModel) View() string {
	header := StyleHeader.Render("새 파이프라인 실행") + "  " + StyleMuted.Render(m.project.Namespace)
	refLine := fmt.Sprintf("  Branch   %s", m.refInput.View())

	varSection := "\n  Variables"
	for _, vf := range m.vars {
		varSection += fmt.Sprintf("\n  %s  =  %s", vf.keyInput.View(), vf.valueInput.View())
	}

	help := StyleHelp.Render("tab 이동  ctrl+a 변수추가  ctrl+r 실행  esc 취소")
	status := ""
	if m.status != "" {
		status = "\n  " + m.status
	}

	return header + "\n\n" + refLine + "\n" + varSection + "\n" + status + "\n" + help
}
