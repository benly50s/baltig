// internal/tui/repo_search.go
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

type searchResultMsg struct{ projects []gitlab.Project }
type searchErrMsg struct{ err error }

// searchResultItem implements list.DefaultItem.
type searchResultItem struct{ p gitlab.Project }

func (s searchResultItem) Title() string       { return s.p.NameWithNamespace }
func (s searchResultItem) Description() string { return s.p.WebURL }
func (s searchResultItem) FilterValue() string { return s.p.NameWithNamespace }

type RepoSearchModel struct {
	cfg     *config.Config
	client  *gitlab.Client
	input   textinput.Model
	list    list.Model
	status  string
	loading bool
}

func NewRepoSearchModel(cfg *config.Config, client *gitlab.Client) *RepoSearchModel {
	ti := textinput.New()
	ti.Placeholder = "저장소 이름 검색..."
	ti.Width = 40

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = StyleSelected

	l := list.New(nil, delegate, 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	return &RepoSearchModel{
		cfg:    cfg,
		client: client,
		input:  ti,
		list:   l,
	}
}

func (m *RepoSearchModel) doSearch() tea.Cmd {
	query := m.input.Value()
	if query == "" {
		return nil
	}
	client := m.client
	return func() tea.Msg {
		projects, err := client.SearchProjects(query)
		if err != nil {
			return searchErrMsg{err}
		}
		return searchResultMsg{projects}
	}
}

func (m *RepoSearchModel) Init() tea.Cmd {
	return m.input.Focus()
}

func (m *RepoSearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-8)

	case searchResultMsg:
		m.loading = false
		items := make([]list.Item, len(msg.projects))
		for i, p := range msg.projects {
			items[i] = searchResultItem{p}
		}
		_ = m.list.SetItems(items)
		if len(msg.projects) == 0 {
			m.status = "결과 없음"
		} else {
			m.status = fmt.Sprintf("%d개 발견", len(msg.projects))
		}

	case searchErrMsg:
		m.loading = false
		m.status = StyleError.Render("오류: " + msg.err.Error())

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Back):
			return m, func() tea.Msg { return navigateBack{} }

		case msg.String() == "enter" && !m.input.Focused():
			// result list: Enter → add repo
			if i, ok := m.list.SelectedItem().(searchResultItem); ok {
				entry := config.ProjectEntry{
					ID:        i.p.ID,
					Namespace: i.p.NameWithNamespace,
					Name:      i.p.Name,
				}
				m.cfg.AddProject(entry)
				_ = config.Save(m.cfg)
				return m, func() tea.Msg { return repoAdded{} }
			}

		case msg.String() == "enter" && m.input.Focused():
			// search input: Enter → run search
			m.input.Blur()
			m.loading = true
			m.status = "검색 중..."
			return m, m.doSearch()

		case msg.String() == "tab":
			if m.input.Focused() {
				m.input.Blur()
			} else {
				cmd := m.input.Focus()
				return m, cmd
			}
		}
	}

	var cmds []tea.Cmd
	if m.input.Focused() {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *RepoSearchModel) View() string {
	header := StyleHeader.Render("저장소 추가")
	inputLine := "  " + m.input.View()
	statusLine := StyleMuted.Render("  " + m.status)
	help := StyleHelp.Render("enter 검색  tab 목록이동  enter 추가  esc 뒤로")

	return header + "\n\n" + inputLine + "\n" + statusLine + "\n" + m.list.View() + "\n" + help
}
