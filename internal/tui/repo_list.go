// internal/tui/repo_list.go
package tui

import (
	"fmt"
	"strings"

	"github.com/benly/baltig/internal/config"
	"github.com/benly/baltig/internal/gitlab"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// repoItem implements list.DefaultItem for a ProjectEntry.
type repoItem struct {
	project  config.ProjectEntry
	isRecent bool
}

func (r repoItem) Title() string {
	prefix := "  "
	if r.project.Starred {
		prefix = "★ "
	}
	return prefix + r.project.Namespace
}

func (r repoItem) Description() string {
	if r.isRecent {
		return StyleMuted.Render("최근 사용")
	}
	return ""
}

func (r repoItem) FilterValue() string { return r.project.Namespace }

type RepoListModel struct {
	cfg    *config.Config
	client *gitlab.Client
	list   list.Model
}

func NewRepoListModel(cfg *config.Config, client *gitlab.Client) *RepoListModel {
	items := buildRepoItems(cfg)
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = StyleSelected
	delegate.Styles.SelectedDesc = StyleSelected.Faint(true)

	l := list.New(items, delegate, 0, 0)
	l.Title = "baltig"
	l.Styles.Title = StyleHeader
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)

	return &RepoListModel{cfg: cfg, client: client, list: l}
}

func buildRepoItems(cfg *config.Config) []list.Item {
	recentSet := make(map[string]bool)
	for _, r := range cfg.Global.Recents {
		recentSet[r.Namespace] = true
	}

	var starred, rest []list.Item
	for _, p := range cfg.Projects {
		item := repoItem{project: p, isRecent: recentSet[p.Namespace]}
		if p.Starred {
			starred = append(starred, item)
		} else {
			rest = append(rest, item)
		}
	}
	return append(starred, rest...)
}

func (m *RepoListModel) Init() tea.Cmd { return nil }

func (m *RepoListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-4)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Add):
			return m, func() tea.Msg { return navigateToRepoSearch{} }

		case key.Matches(msg, Keys.Delete):
			if i, ok := m.list.SelectedItem().(repoItem); ok {
				m.cfg.RemoveProject(i.project.ID)
				_ = config.Save(m.cfg)
				_ = m.list.SetItems(buildRepoItems(m.cfg))
			}
			return m, nil

		case key.Matches(msg, Keys.Select):
			if i, ok := m.list.SelectedItem().(repoItem); ok {
				m.cfg.AddRecent(i.project.Namespace)
				_ = config.Save(m.cfg)
				project := i.project
				return m, func() tea.Msg {
					return navigateToPipelineList{project: project}
				}
			}

		case key.Matches(msg, Keys.Quit):
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *RepoListModel) View() string {
	rightAlign := strings.Repeat(" ", max(0, 40-len("baltig")))
	header := lipgloss.JoinHorizontal(
		lipgloss.Top,
		StyleHeader.Render("baltig"),
		rightAlign,
		StyleMuted.Render(m.cfg.Global.GitLabURL),
	)
	help := renderHelp("r", "추가", "d", "삭제", "enter", "선택", "q", "종료")

	if len(m.cfg.Projects) == 0 {
		empty := fmt.Sprintf("\n  %s\n  %s",
			StyleMuted.Render("등록된 저장소가 없습니다."),
			StyleMuted.Render("r 키를 눌러 저장소를 추가하세요."),
		)
		return header + "\n" + empty + "\n\n" + help
	}

	return header + "\n" + m.list.View() + "\n" + help
}
