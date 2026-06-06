// internal/tui/job_log.go
package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/benly/baltig/internal/config"
	"github.com/benly/baltig/internal/gitlab"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type logLoadedMsg struct{ content string }
type logErrMsg struct{ err error }
type logTickMsg struct{}

type JobLogModel struct {
	cfg       *config.Config
	client    *gitlab.Client
	projectID int64
	job       gitlab.Job
	viewport  viewport.Model
	content   string
	follow    bool
	loading   bool
	err       string
}

func NewJobLogModel(cfg *config.Config, client *gitlab.Client, projectID int64, job gitlab.Job) *JobLogModel {
	vp := viewport.New(0, 0)
	return &JobLogModel{
		cfg:       cfg,
		client:    client,
		projectID: projectID,
		job:       job,
		viewport:  vp,
		follow:    true,
		loading:   true,
	}
}

func (m *JobLogModel) fetchLog() tea.Cmd {
	client := m.client
	projectID := m.projectID
	jobID := m.job.ID
	return func() tea.Msg {
		content, err := client.GetJobLog(projectID, jobID)
		if err != nil {
			return logErrMsg{err}
		}
		return logLoadedMsg{content}
	}
}

func (m *JobLogModel) tick() tea.Cmd {
	if m.job.Status != "running" {
		return nil
	}
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return logTickMsg{}
	})
}

func (m *JobLogModel) Init() tea.Cmd {
	return tea.Batch(m.fetchLog(), m.tick())
}

func (m *JobLogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 4

	case logLoadedMsg:
		m.loading = false
		m.content = msg.content
		m.viewport.SetContent(m.content)
		if m.follow {
			_ = m.viewport.GotoBottom()
		}

	case logErrMsg:
		m.loading = false
		m.err = msg.err.Error()

	case logTickMsg:
		return m, tea.Batch(m.fetchLog(), m.tick())

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Back):
			return m, func() tea.Msg { return navigateBack{} }

		case key.Matches(msg, Keys.Follow):
			m.follow = !m.follow
			if m.follow {
				_ = m.viewport.GotoBottom()
			}

		case key.Matches(msg, Keys.Open):
			return m, openBrowser(m.job.WebURL)
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *JobLogModel) View() string {
	jobTitle := fmt.Sprintf("%s %s [%s]", StatusIcon(m.job.Status), m.job.Name, m.job.Stage)
	header := StyleHeader.Render(jobTitle)

	followIndicator := ""
	if m.follow {
		followIndicator = StyleSuccess.Render(" [follow]")
	}

	if m.loading {
		return header + "\n\n  " + StyleMuted.Render("로그 로딩 중...")
	}
	if m.err != "" {
		return header + "\n\n  " + StyleError.Render("오류: "+m.err)
	}

	lineCount := strings.Count(m.content, "\n")
	statusLine := StyleMuted.Render(fmt.Sprintf("  %d lines", lineCount)) + followIndicator
	help := renderHelp("f", "follow토글", "o", "브라우저", "esc", "뒤로")

	return header + "\n" + statusLine + "\n" + m.viewport.View() + "\n" + help
}
