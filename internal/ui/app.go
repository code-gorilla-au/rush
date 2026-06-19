package ui

import (
	"context"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/teams"
)

type MsgStateLoaded struct {
	Coach *teams.Coach
	Team  *teams.Team
}

type MsgSwitchPage struct {
	NewPage Page
}

type Page int

const (
	PageTitle Page = iota + 1
	PageCreateCoach
	PageLockerRoom
)

type GlobalState struct {
	Coach *teams.Coach
	Team  *teams.Team
}

type RootModel struct {
	ctx             context.Context
	width           int
	height          int
	theme           IceTheme
	currentPage     Page
	pageTitle       tea.Model
	pageCreateCoach tea.Model
	pageLockerRoom  tea.Model
	globalState     *GlobalState
	teamsSvc        *teams.Service
}

// New returns a new UI model.
func New(teamsService *teams.Service) RootModel {
	state := &GlobalState{}

	return RootModel{
		ctx:             context.Background(),
		theme:           NewIceTheme(),
		currentPage:     PageTitle,
		pageTitle:       NewModelTitle(state),
		pageCreateCoach: NewModelCreateCoach(state),
		pageLockerRoom:  NewModelLockerRoom(state),
		globalState:     state,
		teamsSvc:        teamsService,
	}
}

func (m RootModel) Init() tea.Cmd {
	return func() tea.Msg {
		coach, err := m.teamsSvc.GetDefaultCoach(m.ctx)
		if err != nil {
			return MsgStateLoaded{Coach: nil}
		}

		team, err := m.teamsSvc.GetTeamByCoachID(m.ctx, coach.ID)
		if err != nil {
			return MsgStateLoaded{Coach: nil}
		}

		return MsgStateLoaded{Coach: &coach, Team: &team}
	}
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case MsgStateLoaded:
		m.globalState.Coach = msg.Coach
		m.globalState.Team = msg.Team
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case MsgSwitchPage:
		m.currentPage = msg.NewPage
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		var cmd tea.Cmd
		m.pageTitle, cmd = m.pageTitle.Update(msg)
		cmds = append(cmds, cmd)
		m.pageCreateCoach, cmd = m.pageCreateCoach.Update(msg)
		cmds = append(cmds, cmd)
		m.pageLockerRoom, cmd = m.pageLockerRoom.Update(msg)
		cmds = append(cmds, cmd)
	}

	var cmd tea.Cmd
	switch m.currentPage {
	case PageTitle:
		m.pageTitle, cmd = m.pageTitle.Update(msg)
	case PageCreateCoach:
		m.pageCreateCoach, cmd = m.pageCreateCoach.Update(msg)
	case PageLockerRoom:
		m.pageLockerRoom, cmd = m.pageLockerRoom.Update(msg)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m RootModel) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Initializing...")
	}

	switch m.currentPage {
	case PageTitle:
		return m.pageTitle.View()
	case PageCreateCoach:
		return m.pageCreateCoach.View()
	case PageLockerRoom:
		return m.pageLockerRoom.View()
	}

	return tea.NewView("unknown page")
}
