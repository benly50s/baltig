// internal/tui/pipeline_run.go
package tui

import (
	"fmt"

	"github.com/benly/baltig/internal/config"
	"github.com/benly/baltig/internal/gitlab"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type pipelineCreatedMsg struct{ p gitlab.Pipeline }
type pipelineCreateErrMsg struct{ err error }
type branchesLoadedMsg struct{ branches []string }
type ciVarsLoadedMsg struct{ vars []gitlab.CIVariable }

type runPhase int

const (
	phaseBranchSelect runPhase = iota
	phaseVarForm
)

// branchItem implements list.DefaultItem.
type branchItem struct{ name string }

func (b branchItem) Title() string       { return b.name }
func (b branchItem) Description() string { return "" }
func (b branchItem) FilterValue() string { return b.name }

// varField holds a key-value pair for a pipeline variable input row.
type varField struct {
	keyInput   textinput.Model
	valueInput textinput.Model
	desc       string // from .gitlab-ci.yml description
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

func newVarFieldFromCI(cv gitlab.CIVariable) varField {
	k := textinput.New()
	k.Width = 15
	k.SetValue(cv.Key)

	v := textinput.New()
	v.Width = 25
	v.SetValue(cv.Value)

	return varField{keyInput: k, valueInput: v, desc: cv.Description}
}

type PipelineRunModel struct {
	cfg     *config.Config
	client  *gitlab.Client
	project config.ProjectEntry
	phase   runPhase
	// phase 1: branch selection
	branchList      list.Model
	loadingBranches bool
	selectedBranch  string
	// phase 2: variable form
	vars           []varField
	focusIdx       int
	loadingCIVars  bool
	submitting     bool
	status         string
}

func NewPipelineRunModel(cfg *config.Config, client *gitlab.Client, project config.ProjectEntry) *PipelineRunModel {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = StyleSelected

	l := list.New(nil, delegate, 0, 0)
	l.Title = "브랜치 선택"
	l.Styles.Title = StyleHeader
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)

	return &PipelineRunModel{
		cfg:             cfg,
		client:          client,
		project:         project,
		phase:           phaseBranchSelect,
		branchList:      l,
		loadingBranches: true,
	}
}

func (m *PipelineRunModel) loadBranches() tea.Cmd {
	client := m.client
	projectID := m.project.ID
	return func() tea.Msg {
		branches, err := client.ListBranches(projectID)
		if err != nil {
			// fallback: empty list, user can still filter/type
			return branchesLoadedMsg{branches: []string{m.cfg.Global.DefaultRef}}
		}
		return branchesLoadedMsg{branches: branches}
	}
}

func (m *PipelineRunModel) loadCIVars(ref string) tea.Cmd {
	client := m.client
	projectID := m.project.ID
	return func() tea.Msg {
		vars, _ := client.GetCIVariables(projectID, ref)
		return ciVarsLoadedMsg{vars: vars}
	}
}

func (m *PipelineRunModel) Init() tea.Cmd {
	return m.loadBranches()
}

func (m *PipelineRunModel) totalVarFields() int {
	return len(m.vars) * 2
}

func (m *PipelineRunModel) updateVarFocus() tea.Cmd {
	for i := range m.vars {
		m.vars[i].keyInput.Blur()
		m.vars[i].valueInput.Blur()
	}
	if len(m.vars) == 0 {
		return nil
	}
	varIdx := m.focusIdx / 2
	isKey := m.focusIdx%2 == 0
	if varIdx < len(m.vars) {
		if isKey {
			return m.vars[varIdx].keyInput.Focus()
		}
		return m.vars[varIdx].valueInput.Focus()
	}
	return nil
}

func (m *PipelineRunModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.branchList.SetSize(msg.Width, msg.Height-4)

	case branchesLoadedMsg:
		m.loadingBranches = false
		items := make([]list.Item, len(msg.branches))
		for i, b := range msg.branches {
			items[i] = branchItem{b}
		}
		_ = m.branchList.SetItems(items)
		// Pre-select default branch
		for i, b := range msg.branches {
			if b == m.cfg.Global.DefaultRef {
				m.branchList.Select(i)
				break
			}
		}

	case ciVarsLoadedMsg:
		m.loadingCIVars = false
		if len(msg.vars) > 0 {
			m.vars = make([]varField, len(msg.vars))
			for i, cv := range msg.vars {
				m.vars[i] = newVarFieldFromCI(cv)
			}
		} else {
			// No CI vars: start with one empty field
			m.vars = []varField{newVarField()}
		}
		m.focusIdx = 0
		return m, m.updateVarFocus()

	case pipelineCreatedMsg:
		m.submitting = false
		return m, func() tea.Msg { return navigateBack{} }

	case pipelineCreateErrMsg:
		m.submitting = false
		m.status = StyleError.Render("오류: " + msg.err.Error())

	case tea.KeyMsg:
		switch m.phase {
		case phaseBranchSelect:
			switch {
			case key.Matches(msg, Keys.Back):
				return m, func() tea.Msg { return navigateBack{} }

			case key.Matches(msg, Keys.Select):
				if i, ok := m.branchList.SelectedItem().(branchItem); ok {
					m.selectedBranch = i.name
					m.phase = phaseVarForm
					m.loadingCIVars = true
					return m, m.loadCIVars(m.selectedBranch)
				}
			}

		case phaseVarForm:
			switch {
			case key.Matches(msg, Keys.Back):
				// Go back to branch selection
				m.phase = phaseBranchSelect
				m.status = ""
				return m, nil

			case msg.String() == "ctrl+r" && !m.submitting:
				return m, m.submit()

			case msg.String() == "tab":
				if m.totalVarFields() > 0 {
					m.focusIdx = (m.focusIdx + 1) % m.totalVarFields()
					return m, m.updateVarFocus()
				}

			case msg.String() == "shift+tab":
				if m.totalVarFields() > 0 {
					m.focusIdx = (m.focusIdx - 1 + m.totalVarFields()) % m.totalVarFields()
					return m, m.updateVarFocus()
				}

			case msg.String() == "ctrl+a":
				m.vars = append(m.vars, newVarField())
				return m, nil
			}
		}
	}

	// Delegate to active sub-model
	var cmd tea.Cmd
	switch m.phase {
	case phaseBranchSelect:
		m.branchList, cmd = m.branchList.Update(msg)
	case phaseVarForm:
		if len(m.vars) > 0 {
			varIdx := m.focusIdx / 2
			isKey := m.focusIdx%2 == 0
			if varIdx < len(m.vars) {
				if isKey {
					m.vars[varIdx].keyInput, cmd = m.vars[varIdx].keyInput.Update(msg)
				} else {
					m.vars[varIdx].valueInput, cmd = m.vars[varIdx].valueInput.Update(msg)
				}
			}
		}
	}
	return m, cmd
}

func (m *PipelineRunModel) submit() tea.Cmd {
	m.submitting = true
	m.status = StyleMuted.Render("파이프라인 생성 중...")

	vars := make([]gitlab.PipelineVariable, 0, len(m.vars))
	for _, vf := range m.vars {
		k := vf.keyInput.Value()
		v := vf.valueInput.Value()
		if k != "" {
			vars = append(vars, gitlab.PipelineVariable{Key: k, Value: v})
		}
	}

	client := m.client
	projectID := m.project.ID
	ref := m.selectedBranch

	return func() tea.Msg {
		p, err := client.CreatePipeline(projectID, ref, vars)
		if err != nil {
			return pipelineCreateErrMsg{err}
		}
		return pipelineCreatedMsg{*p}
	}
}

func (m *PipelineRunModel) View() string {
	header := StyleHeader.Render("새 파이프라인") + "  " + StyleMuted.Render(m.project.Namespace)

	switch m.phase {
	case phaseBranchSelect:
		if m.loadingBranches {
			return header + "\n\n  " + StyleMuted.Render("브랜치 로딩 중...")
		}
		help := renderHelp("enter", "선택", "esc", "취소")
		return header + "\n" + m.branchList.View() + "\n" + help

	case phaseVarForm:
		if m.loadingCIVars {
			return header + "\n\n  " + StyleMuted.Render("CI 변수 로딩 중...")
		}

		branchLine := "  Branch   " + StyleSuccess.Render(m.selectedBranch)
		varSection := "\n  Variables"
		if len(m.vars) == 0 {
			varSection += "\n  " + StyleMuted.Render("(없음)")
		}
		for _, vf := range m.vars {
			descStr := ""
			if vf.desc != "" {
				descStr = "  " + StyleMuted.Render("# "+vf.desc)
			}
			varSection += fmt.Sprintf("\n  %s  =  %s%s", vf.keyInput.View(), vf.valueInput.View(), descStr)
		}

		status := ""
		if m.status != "" {
			status = "\n  " + m.status
		}
		help := renderHelp("tab", "이동", "ctrl+a", "변수추가", "ctrl+r", "실행", "esc", "뒤로")
		return header + "\n\n" + branchLine + "\n" + varSection + "\n" + status + "\n" + help
	}
	return header
}
