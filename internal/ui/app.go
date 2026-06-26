package ui

import (
	"context"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/tournament"
)

type MsgStateUpdated struct {
	Coach *teams.Coach
	Team  *teams.Team
}

type MsgSwitchPage struct {
	NewPage  Page
	Playbook *playbooks.Playbook
}

type Page int

const (
	PageTitle Page = iota + 1
	PageCreateCoach
	PageLockerRoom
	PageLockerPlayers
	PageLockerPlaybooksList
	PageLockerPlaybooksCreate
	PageLockerPlaybooksEdit
	PageNewTournament
	PageNewBattleSelection
	PageTitleSettings
)

type GlobalState struct {
	Coach *teams.Coach
	Team  *teams.Team
}

func (m *GlobalState) Context() context.Context {
	return context.Background()
}

type RootModel struct {
	ctx                       context.Context
	width                     int
	height                    int
	theme                     IceTheme
	currentPage               Page
	pageTitle                 tea.Model
	pageCreateCoach           tea.Model
	pageLockerRoom            tea.Model
	pageLockerPlayers         tea.Model
	pageLockerPlaybooksList   tea.Model
	pageLockerPlaybooksCreate tea.Model
	pageLockerPlaybooksEdit   tea.Model
	pageNewTournament         tea.Model
	pageNewBattleSelection    tea.Model
	pageTitleSettings         tea.Model
	globalState               *GlobalState
	teamsSvc                  *teams.Service
	playbookSvc               *playbooks.Service
	gameSvc                   *games.Service
	aiTeamsSvc                *tournament.AITeamService
}

type Dependencies struct {
	teamsSvc    *teams.Service
	playbookSvc *playbooks.Service
	gameSvc     *games.Service
	aiTeamsSvc  *tournament.AITeamService
}

// New returns a new UI model.
func New(deps Dependencies) RootModel {
	state := &GlobalState{}

	return RootModel{
		ctx:                       context.Background(),
		theme:                     NewIceTheme(),
		currentPage:               PageTitle,
		pageTitle:                 NewModelTitle(state),
		pageCreateCoach:           NewModelCreateCoach(state, deps.teamsSvc),
		pageLockerRoom:            NewModelLockerRoom(state),
		pageLockerPlayers:         NewModelLockerPlayers(state, deps.teamsSvc),
		pageLockerPlaybooksList:   NewModelLockerPlaybooksList(state, deps.playbookSvc),
		pageLockerPlaybooksCreate: NewModelLockerPlaybooksCreate(state, deps.playbookSvc),
		pageLockerPlaybooksEdit:   NewModelLockerPlaybooksEdit(state, deps.playbookSvc),
		pageNewTournament:         NewModelNewTournament(state),
		pageNewBattleSelection:    NewModelNewBattleSelection(state, deps.teamsSvc, deps.playbookSvc),
		pageTitleSettings:         NewModelTitleSettings(state),
		globalState:               state,
		teamsSvc:                  deps.teamsSvc,
		playbookSvc:               deps.playbookSvc,
		gameSvc:                   deps.gameSvc,
		aiTeamsSvc:                deps.aiTeamsSvc,
	}
}

func (m RootModel) Init() tea.Cmd {
	return func() tea.Msg {
		coach, err := m.teamsSvc.GetDefaultCoach(m.ctx)
		if err != nil {
			return MsgStateUpdated{Coach: nil}
		}

		team, err := m.teamsSvc.GetTeamByCoachID(m.ctx, coach.ID)
		if err != nil {
			return MsgStateUpdated{Coach: nil}
		}

		return MsgStateUpdated{Coach: &coach, Team: &team}
	}
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case MsgStateUpdated:
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
		m.pageLockerPlayers, cmd = m.pageLockerPlayers.Update(msg)
		cmds = append(cmds, cmd)
		m.pageLockerPlaybooksList, cmd = m.pageLockerPlaybooksList.Update(msg)
		cmds = append(cmds, cmd)
		m.pageLockerPlaybooksCreate, cmd = m.pageLockerPlaybooksCreate.Update(msg)
		cmds = append(cmds, cmd)
		m.pageLockerPlaybooksEdit, cmd = m.pageLockerPlaybooksEdit.Update(msg)
		cmds = append(cmds, cmd)
		m.pageNewTournament, cmd = m.pageNewTournament.Update(msg)
		cmds = append(cmds, cmd)
		m.pageNewBattleSelection, cmd = m.pageNewBattleSelection.Update(msg)
		cmds = append(cmds, cmd)
		m.pageTitleSettings, cmd = m.pageTitleSettings.Update(msg)
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
	case PageLockerPlayers:
		m.pageLockerPlayers, cmd = m.pageLockerPlayers.Update(msg)
	case PageLockerPlaybooksList:
		m.pageLockerPlaybooksList, cmd = m.pageLockerPlaybooksList.Update(msg)
	case PageLockerPlaybooksCreate:
		m.pageLockerPlaybooksCreate, cmd = m.pageLockerPlaybooksCreate.Update(msg)
	case PageLockerPlaybooksEdit:
		m.pageLockerPlaybooksEdit, cmd = m.pageLockerPlaybooksEdit.Update(msg)
	case PageNewTournament:
		m.pageNewTournament, cmd = m.pageNewTournament.Update(msg)
	case PageNewBattleSelection:
		m.pageNewBattleSelection, cmd = m.pageNewBattleSelection.Update(msg)
	case PageTitleSettings:
		m.pageTitleSettings, cmd = m.pageTitleSettings.Update(msg)
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
	case PageLockerPlayers:
		return m.pageLockerPlayers.View()
	case PageLockerPlaybooksList:
		return m.pageLockerPlaybooksList.View()
	case PageLockerPlaybooksCreate:
		return m.pageLockerPlaybooksCreate.View()
	case PageLockerPlaybooksEdit:
		return m.pageLockerPlaybooksEdit.View()
	case PageNewTournament:
		return m.pageNewTournament.View()
	case PageNewBattleSelection:
		return m.pageNewBattleSelection.View()
	case PageTitleSettings:
		return m.pageTitleSettings.View()
	}

	return tea.NewView("unknown page")
}
