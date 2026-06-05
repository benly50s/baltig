package tui

import (
	"github.com/benly/baltig/internal/config"
	"github.com/benly/baltig/internal/gitlab"
	tea "github.com/charmbracelet/bubbletea"
)

type PipelineDetailModel struct {
	cfg       *config.Config
	client    *gitlab.Client
	projectID int64
	pipeline  gitlab.Pipeline
}

func NewPipelineDetailModel(cfg *config.Config, client *gitlab.Client, projectID int64, pipeline gitlab.Pipeline) *PipelineDetailModel {
	return &PipelineDetailModel{cfg: cfg, client: client, projectID: projectID, pipeline: pipeline}
}
func (m *PipelineDetailModel) Init() tea.Cmd                           { return nil }
func (m *PipelineDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m *PipelineDetailModel) View() string                            { return "pipeline detail (stub)" }
