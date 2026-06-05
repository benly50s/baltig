package tui

import (
	"github.com/benly/baltig/internal/config"
	"github.com/benly/baltig/internal/gitlab"
	tea "github.com/charmbracelet/bubbletea"
)

type RepoSearchModel struct {
	cfg    *config.Config
	client *gitlab.Client
}

func NewRepoSearchModel(cfg *config.Config, client *gitlab.Client) *RepoSearchModel {
	return &RepoSearchModel{cfg: cfg, client: client}
}
func (m *RepoSearchModel) Init() tea.Cmd                           { return nil }
func (m *RepoSearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m *RepoSearchModel) View() string                            { return "repo search (stub)" }
