package tui

import (
	"github.com/benly/baltig/internal/config"
	"github.com/benly/baltig/internal/gitlab"
	tea "github.com/charmbracelet/bubbletea"
)

type PipelineRunModel struct {
	cfg     *config.Config
	client  *gitlab.Client
	project config.ProjectEntry
}

func NewPipelineRunModel(cfg *config.Config, client *gitlab.Client, project config.ProjectEntry) *PipelineRunModel {
	return &PipelineRunModel{cfg: cfg, client: client, project: project}
}
func (m *PipelineRunModel) Init() tea.Cmd                           { return nil }
func (m *PipelineRunModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m *PipelineRunModel) View() string                            { return "pipeline run (stub)" }
