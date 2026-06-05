// internal/tui/pipeline_detail.go
package tui

import (
	"fmt"

	"github.com/benly/baltig/internal/config"
	"github.com/benly/baltig/internal/gitlab"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type jobsLoadedMsg struct{ jobs []gitlab.Job }
type jobsErrMsg struct{ err error }

// jobItem implements list.DefaultItem.
type jobItem struct{ j gitlab.Job }

func (ji jobItem) Title() string {
	icon := StatusIcon(ji.j.Status)
	return fmt.Sprintf("%s %s  [%s]", icon, ji.j.Name, ji.j.Stage)
}
func (ji jobItem) Description() string { return "" }
func (ji jobItem) FilterValue() string { return ji.j.Name }

type PipelineDetailModel struct {
	cfg       *config.Config
	client    *gitlab.Client
	projectID int64
	pipeline  gitlab.Pipeline
	list      list.Model
	jobs      []gitlab.Job
	loading   bool
	err       string
}

func NewPipelineDetailModel(cfg *config.Config, client *gitlab.Client, projectID int64, pipeline gitlab.Pipeline) *PipelineDetailModel {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = StyleSelected

	title := fmt.Sprintf("#%d  %s  %s", pipeline.ID, pipeline.Ref, pipeline.Status)
	l := list.New(nil, delegate, 0, 0)
	l.Title = title
	l.Styles.Title = StyleHeader
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	return &PipelineDetailModel{
		cfg:       cfg,
		client:    client,
		projectID: projectID,
		pipeline:  pipeline,
		list:      l,
		loading:   true,
	}
}

func (m *PipelineDetailModel) Init() tea.Cmd {
	client := m.client
	projectID := m.projectID
	pipelineID := m.pipeline.ID
	return func() tea.Msg {
		jobs, err := client.ListJobs(projectID, pipelineID)
		if err != nil {
			return jobsErrMsg{err}
		}
		return jobsLoadedMsg{jobs}
	}
}

func (m *PipelineDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-4)

	case jobsLoadedMsg:
		m.loading = false
		m.jobs = msg.jobs
		items := make([]list.Item, len(msg.jobs))
		for i, j := range msg.jobs {
			items[i] = jobItem{j}
		}
		_ = m.list.SetItems(items)

	case jobsErrMsg:
		m.loading = false
		m.err = msg.err.Error()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Back):
			return m, func() tea.Msg { return navigateBack{} }

		case key.Matches(msg, Keys.Open):
			return m, openBrowser(m.pipeline.WebURL)

		case key.Matches(msg, Keys.Select):
			if i, ok := m.list.SelectedItem().(jobItem); ok {
				pID := m.projectID
				j := i.j
				return m, func() tea.Msg {
					return navigateToJobLog{projectID: pID, job: j}
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *PipelineDetailModel) View() string {
	if m.loading {
		title := fmt.Sprintf("#%d  %s", m.pipeline.ID, m.pipeline.Ref)
		return StyleHeader.Render(title) + "\n\n  " + StyleMuted.Render("로딩 중...")
	}
	if m.err != "" {
		return StyleError.Render("오류: " + m.err)
	}
	help := StyleHelp.Render("enter 로그  o 브라우저  esc 뒤로")
	return m.list.View() + "\n" + help
}
