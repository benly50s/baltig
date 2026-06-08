// internal/tui/pipeline_detail.go
package tui

import (
	"fmt"
	"strings"

	"github.com/benly/baltig/internal/config"
	"github.com/benly/baltig/internal/gitlab"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type jobsLoadedMsg struct{ jobs []gitlab.Job }
type jobsErrMsg struct{ err error }

type PipelineDetailModel struct {
	cfg        *config.Config
	client     *gitlab.Client
	projectID  int64
	pipeline   gitlab.Pipeline
	vp         viewport.Model
	jobs       []gitlab.Job // flat, for cursor indexing
	stageOrder []string     // ordered unique stages
	cursorIdx  int          // index into jobs slice
	loading    bool
	err        string
	width      int
	height     int
}

func NewPipelineDetailModel(cfg *config.Config, client *gitlab.Client, projectID int64, pipeline gitlab.Pipeline) *PipelineDetailModel {
	vp := viewport.New(0, 0)
	return &PipelineDetailModel{
		cfg:       cfg,
		client:    client,
		projectID: projectID,
		pipeline:  pipeline,
		vp:        vp,
		loading:   true,
	}
}

func (m *PipelineDetailModel) Init() tea.Cmd {
	client := m.client
	projectID := m.projectID
	pipelineID := m.pipeline.ID
	return func() tea.Msg {
		jobs, err := client.ListJobs(projectID, pipelineID)
		if err != nil {
			return jobsErrMsg{err}
		}
		return jobsLoadedMsg{jobs}
	}
}

// buildStageOrder extracts unique stages preserving insertion order.
func buildStageOrder(jobs []gitlab.Job) []string {
	seen := make(map[string]bool)
	var order []string
	for _, j := range jobs {
		if !seen[j.Stage] {
			seen[j.Stage] = true
			order = append(order, j.Stage)
		}
	}
	return order
}

// renderContent builds the viewport content and returns (content, jobLineOffsets).
// jobLineOffsets[i] = the line number (0-based) in the content where jobs[i] is rendered.
func (m *PipelineDetailModel) renderContent() (string, []int) {
	var sb strings.Builder
	lineOffsets := make([]int, len(m.jobs))
	lineNum := 0

	// Map stage → jobs in that stage
	stageJobs := make(map[string][]int) // stage → indices into m.jobs
	for i, j := range m.jobs {
		stageJobs[j.Stage] = append(stageJobs[j.Stage], i)
	}

	isLast := func(stage string) bool {
		return len(m.stageOrder) > 0 && m.stageOrder[len(m.stageOrder)-1] == stage
	}

	for _, stage := range m.stageOrder {
		// Compute stage aggregate status
		stageStatus := aggregateStatus(m.jobs, stageJobs[stage])

		stageHeader := fmt.Sprintf("  %s %s",
			StatusIcon(stageStatus),
			StyleSubHeader.Bold(true).Render(stage),
		)
		sb.WriteString(stageHeader + "\n")
		lineNum++

		for _, jobIdx := range stageJobs[stage] {
			j := m.jobs[jobIdx]
			lineOffsets[jobIdx] = lineNum

			cursor := "  "
			nameStyle := StyleMuted
			if jobIdx == m.cursorIdx {
				cursor = StylePrimary.Render("▶")
				nameStyle = StyleSelected.Bold(true)
			}

			row := fmt.Sprintf("  %s  %s  %s",
				cursor,
				StatusIcon(j.Status),
				nameStyle.Render(j.Name),
			)
			sb.WriteString(row + "\n")
			lineNum++
		}

		if !isLast(stage) {
			sb.WriteString("\n")
			lineNum++
		}
	}

	return sb.String(), lineOffsets
}

// aggregateStatus returns a representative status for a stage given its job indices.
func aggregateStatus(jobs []gitlab.Job, indices []int) string {
	has := func(s string) bool {
		for _, i := range indices {
			if jobs[i].Status == s {
				return true
			}
		}
		return false
	}
	switch {
	case has("failed"):
		return "failed"
	case has("running"):
		return "running"
	case has("pending"):
		return "pending"
	case has("canceled"):
		return "canceled"
	default:
		return "success"
	}
}

// syncViewport updates viewport content and scrolls to keep cursor visible.
func (m *PipelineDetailModel) syncViewport() {
	content, offsets := m.renderContent()
	m.vp.SetContent(content)

	if len(offsets) == 0 || m.cursorIdx >= len(offsets) {
		return
	}
	cursorLine := offsets[m.cursorIdx]
	top := m.vp.YOffset
	bottom := top + m.vp.Height - 1
	if cursorLine < top {
		m.vp.SetYOffset(cursorLine)
	} else if cursorLine > bottom {
		m.vp.SetYOffset(cursorLine - m.vp.Height + 1)
	}
}

func (m *PipelineDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	const headerLines = 3 // header + status bar + blank

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.vp.Width = msg.Width
		m.vp.Height = msg.Height - headerLines - 2 // reserve for help bar
		if !m.loading {
			m.syncViewport()
		}

	case jobsLoadedMsg:
		m.loading = false
		m.jobs = msg.jobs
		m.stageOrder = buildStageOrder(msg.jobs)
		m.cursorIdx = 0
		m.vp.Width = m.width
		m.vp.Height = m.height - headerLines - 2
		m.syncViewport()

	case jobsErrMsg:
		m.loading = false
		m.err = msg.err.Error()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Back):
			return m, func() tea.Msg { return navigateBack{} }

		case key.Matches(msg, Keys.Open):
			return m, openBrowser(m.pipeline.WebURL)

		case key.Matches(msg, Keys.Select):
			if m.cursorIdx < len(m.jobs) {
				j := m.jobs[m.cursorIdx]
				pID := m.projectID
				return m, func() tea.Msg {
					return navigateToJobLog{projectID: pID, job: j}
				}
			}

		case key.Matches(msg, Keys.Up):
			if m.cursorIdx > 0 {
				m.cursorIdx--
				m.syncViewport()
			}
			return m, nil

		case key.Matches(msg, Keys.Down):
			if m.cursorIdx < len(m.jobs)-1 {
				m.cursorIdx++
				m.syncViewport()
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

func (m *PipelineDetailModel) View() string {
	titleStr := fmt.Sprintf("#%d  %s", m.pipeline.ID, m.pipeline.Ref)
	statusStr := StatusIcon(m.pipeline.Status) + " " + m.pipeline.Status
	header := StyleHeader.Render(titleStr) + "  " + StyleMuted.Render(statusStr)

	if m.loading {
		return header + "\n\n  " + StyleMuted.Render("로딩 중...")
	}
	if m.err != "" {
		return header + "\n\n  " + StyleError.Render("오류: "+m.err)
	}
	if len(m.jobs) == 0 {
		return header + "\n\n  " + StyleMuted.Render("job이 없습니다")
	}

	help := renderHelp("↑↓", "이동", "enter", "로그", "o", "브라우저", "esc", "뒤로")
	return header + "\n\n" + m.vp.View() + "\n" + help
}
