// internal/tui/pipeline_list.go
package tui

import (
	"fmt"
	"time"

	"github.com/benly/baltig/internal/config"
	"github.com/benly/baltig/internal/gitlab"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type pipelinesLoadedMsg struct{ pipelines []gitlab.Pipeline }
type pipelinesErrMsg struct{ err error }

// pipelineItem implements list.DefaultItem.
type pipelineItem struct{ p gitlab.Pipeline }

func (pi pipelineItem) Title() string {
	icon := StatusIcon(pi.p.Status)
	ago := formatAgo(pi.p.CreatedAt)
	return fmt.Sprintf("%s #%d  %s  %s  %s", icon, pi.p.ID, pi.p.Ref, pi.p.Status, ago)
}
func (pi pipelineItem) Description() string { return "" }
func (pi pipelineItem) FilterValue() string { return pi.p.Ref }

type PipelineListModel struct {
	cfg     *config.Config
	client  *gitlab.Client
	project config.ProjectEntry
	list    list.Model
	loading bool
	err     string
}

func NewPipelineListModel(cfg *config.Config, client *gitlab.Client, project config.ProjectEntry) *PipelineListModel {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = StyleSelected

	l := list.New(nil, delegate, 0, 0)
	l.Title = project.Namespace
	l.Styles.Title = StyleHeader
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	return &PipelineListModel{
		cfg:     cfg,
		client:  client,
		project: project,
		list:    l,
		loading: true,
	}
}

func (m *PipelineListModel) loadPipelines() tea.Cmd {
	client := m.client
	projectID := int64(m.project.ID)
	return func() tea.Msg {
		pipelines, err := client.ListPipelines(projectID)
		if err != nil {
			return pipelinesErrMsg{err}
		}
		return pipelinesLoadedMsg{pipelines}
	}
}

func (m *PipelineListModel) Init() tea.Cmd {
	return m.loadPipelines()
}

func (m *PipelineListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-4)

	case pipelinesLoadedMsg:
		m.loading = false
		items := make([]list.Item, len(msg.pipelines))
		for i, p := range msg.pipelines {
			items[i] = pipelineItem{p}
		}
		_ = m.list.SetItems(items)

	case pipelinesErrMsg:
		m.loading = false
		m.err = msg.err.Error()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Back):
			return m, func() tea.Msg { return navigateBack{} }

		case key.Matches(msg, Keys.New):
			project := m.project
			return m, func() tea.Msg { return navigateToPipelineRun{project: project} }

		case key.Matches(msg, Keys.Select):
			if i, ok := m.list.SelectedItem().(pipelineItem); ok {
				pID := int64(m.project.ID)
				p := i.p
				return m, func() tea.Msg {
					return navigateToPipelineDetail{projectID: pID, pipeline: p}
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *PipelineListModel) View() string {
	if m.loading {
		return StyleHeader.Render(m.project.Namespace) + "\n\n  " + StyleMuted.Render("로딩 중...")
	}
	if m.err != "" {
		return StyleHeader.Render(m.project.Namespace) + "\n\n  " + StyleError.Render("오류: "+m.err)
	}
	help := StyleHelp.Render("n 새 파이프라인  enter 상세  esc 뒤로")
	return m.list.View() + "\n" + help
}

func formatAgo(t *time.Time) string {
	if t == nil {
		return ""
	}
	d := time.Since(*t)
	switch {
	case d < time.Minute:
		return "방금 전"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}
