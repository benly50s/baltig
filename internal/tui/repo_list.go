package tui

import (
	"github.com/benly/baltig/internal/config"
	"github.com/benly/baltig/internal/gitlab"
	tea "github.com/charmbracelet/bubbletea"
)

type RepoListModel struct {
	cfg    *config.Config
	client *gitlab.Client
}

func NewRepoListModel(cfg *config.Config, client *gitlab.Client) *RepoListModel {
	return &RepoListModel{cfg: cfg, client: client}
}
func (m *RepoListModel) Init() tea.Cmd                           { return nil }
func (m *RepoListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m *RepoListModel) View() string                            { return "repo list (stub)" }
