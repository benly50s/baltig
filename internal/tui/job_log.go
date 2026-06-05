package tui

import (
	"github.com/benly/baltig/internal/config"
	"github.com/benly/baltig/internal/gitlab"
	tea "github.com/charmbracelet/bubbletea"
)

type JobLogModel struct {
	cfg       *config.Config
	client    *gitlab.Client
	projectID int64
	job       gitlab.Job
}

func NewJobLogModel(cfg *config.Config, client *gitlab.Client, projectID int64, job gitlab.Job) *JobLogModel {
	return &JobLogModel{cfg: cfg, client: client, projectID: projectID, job: job}
}
func (m *JobLogModel) Init() tea.Cmd                           { return nil }
func (m *JobLogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m *JobLogModel) View() string                            { return "job log (stub)" }
