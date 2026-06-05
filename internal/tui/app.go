// internal/tui/app.go
package tui

import (
	"github.com/benly/baltig/internal/config"
	"github.com/benly/baltig/internal/gitlab"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	stateRepoList state = iota
	stateRepoSearch
	statePipelineList
	statePipelineDetail
	statePipelineRun
	stateJobLog
)

// AppModel is the root BubbleTea model. It owns the state machine
// and delegates Init/Update/View to the active child model.
type AppModel struct {
	state  state
	cfg    *config.Config
	client *gitlab.Client

	repoList       *RepoListModel
	repoSearch     *RepoSearchModel
	pipelineList   *PipelineListModel
	pipelineDetail *PipelineDetailModel
	pipelineRun    *PipelineRunModel
	jobLog         *JobLogModel

	width  int
	height int
}

// NewApp creates the root AppModel. cfg and client must be non-nil.
func NewApp(cfg *config.Config, client *gitlab.Client) *AppModel {
	m := &AppModel{
		cfg:    cfg,
		client: client,
		state:  stateRepoList,
	}
	m.repoList = NewRepoListModel(cfg, client)
	return m
}

func (m *AppModel) Init() tea.Cmd {
	return m.repoList.Init()
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	// Navigation messages from child models
	case navigateToRepoSearch:
		m.state = stateRepoSearch
		m.repoSearch = NewRepoSearchModel(m.cfg, m.client)
		return m, m.repoSearch.Init()

	case navigateToPipelineList:
		m.state = statePipelineList
		m.pipelineList = NewPipelineListModel(m.cfg, m.client, msg.project)
		return m, m.pipelineList.Init()

	case navigateToPipelineDetail:
		m.state = statePipelineDetail
		m.pipelineDetail = NewPipelineDetailModel(m.cfg, m.client, msg.projectID, msg.pipeline)
		return m, m.pipelineDetail.Init()

	case navigateToPipelineRun:
		m.state = statePipelineRun
		m.pipelineRun = NewPipelineRunModel(m.cfg, m.client, msg.project)
		return m, m.pipelineRun.Init()

	case navigateToJobLog:
		m.state = stateJobLog
		m.jobLog = NewJobLogModel(m.cfg, m.client, msg.projectID, msg.job)
		return m, m.jobLog.Init()

	case navigateBack:
		return m.handleBack()

	case repoAdded:
		m.state = stateRepoList
		m.repoList = NewRepoListModel(m.cfg, m.client)
		return m, m.repoList.Init()
	}

	return m, m.updateChild(msg)
}

func (m *AppModel) updateChild(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.state {
	case stateRepoList:
		_, cmd = m.repoList.Update(msg)
	case stateRepoSearch:
		_, cmd = m.repoSearch.Update(msg)
	case statePipelineList:
		_, cmd = m.pipelineList.Update(msg)
	case statePipelineDetail:
		_, cmd = m.pipelineDetail.Update(msg)
	case statePipelineRun:
		_, cmd = m.pipelineRun.Update(msg)
	case stateJobLog:
		_, cmd = m.jobLog.Update(msg)
	}
	return cmd
}

func (m *AppModel) handleBack() (tea.Model, tea.Cmd) {
	switch m.state {
	case stateRepoSearch:
		m.state = stateRepoList
		return m, nil
	case statePipelineList:
		m.state = stateRepoList
		return m, nil
	case statePipelineDetail, statePipelineRun:
		m.state = statePipelineList
		return m, m.pipelineList.Init()
	case stateJobLog:
		m.state = statePipelineDetail
		return m, nil
	}
	return m, nil
}

func (m *AppModel) View() string {
	switch m.state {
	case stateRepoList:
		return m.repoList.View()
	case stateRepoSearch:
		return m.repoSearch.View()
	case statePipelineList:
		return m.pipelineList.View()
	case statePipelineDetail:
		return m.pipelineDetail.View()
	case statePipelineRun:
		return m.pipelineRun.View()
	case stateJobLog:
		return m.jobLog.View()
	}
	return ""
}

// Navigation messages — child models emit these to trigger state transitions.
type navigateToRepoSearch struct{}
type navigateToPipelineList struct{ project config.ProjectEntry }
type navigateToPipelineDetail struct {
	projectID int64
	pipeline  gitlab.Pipeline
}
type navigateToPipelineRun struct{ project config.ProjectEntry }
type navigateToJobLog struct {
	projectID int64
	job       gitlab.Job
}
type navigateBack struct{}
type repoAdded struct{}
