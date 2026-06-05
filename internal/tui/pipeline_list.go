package tui

import (
	"github.com/benly/baltig/internal/config"
	"github.com/benly/baltig/internal/gitlab"
	tea "github.com/charmbracelet/bubbletea"
)

type PipelineListModel struct {
	cfg     *config.Config
	client  *gitlab.Client
	project config.ProjectEntry
}

func NewPipelineListModel(cfg *config.Config, client *gitlab.Client, project config.ProjectEntry) *PipelineListModel {
	return &PipelineListModel{cfg: cfg, client: client, project: project}
}
func (m *PipelineListModel) Init() tea.Cmd                           { return nil }
func (m *PipelineListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m *PipelineListModel) View() string                            { return "pipeline list (stub)" }
