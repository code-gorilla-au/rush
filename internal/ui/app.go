package ui

import (
	"context"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

type MsgStateUpdated struct {
	Coach *teams.Coach
	Team  *teams.Team
}

type MsgSwitchPage struct {
	NewPage  Page
	Playbook *playbooks.Playbook
	GameID   int64
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
	PageGame
	PageGameComplete
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
	theme                     styles.IceTheme
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
	pageGame                  tea.Model
	pageGameComplete          tea.Model
	globalState               *GlobalState
	teamsSvc                  *teams.Service
	playbookSvc               *playbooks.Service
	gameSvc                   *games.Service
}

type Dependencies struct {
	TeamsSvc    *teams.Service
	PlaybookSvc *playbooks.Service
	GameSvc     *games.Service
}

// New returns a new UI model.
func New(deps Dependencies) *RootModel {
	state := &GlobalState{}
	theme := styles.NewIceTheme()

	return &RootModel{
		ctx:                       context.Background(),
		theme:                     theme,
		currentPage:               PageTitle,
		pageTitle:                 NewModelTitle(state, theme),
		pageCreateCoach:           NewModelCreateCoach(state, deps.TeamsSvc, theme),
		pageLockerRoom:            NewModelLockerRoom(state, theme),
		pageLockerPlayers:         NewModelLockerPlayers(state, deps.TeamsSvc, theme),
		pageLockerPlaybooksList:   NewModelLockerPlaybooksList(state, deps.PlaybookSvc, theme),
		pageLockerPlaybooksCreate: NewModelLockerPlaybooksCreate(state, deps.PlaybookSvc, theme),
		pageLockerPlaybooksEdit:   NewModelLockerPlaybooksEdit(state, deps.PlaybookSvc, theme),
		pageNewTournament:         NewModelNewTournament(state, theme),
		pageNewBattleSelection:    NewModelNewBattleSelection(state, deps.TeamsSvc, deps.PlaybookSvc, deps.GameSvc, theme),
		pageTitleSettings:         NewModelTitleSettings(state, theme),
		pageGame:                  NewModelGame(state, deps.GameSvc, theme),
		pageGameComplete:          NewPageGameComplete(state, deps.TeamsSvc, deps.GameSvc, theme),
		globalState:               state,
		teamsSvc:                  deps.TeamsSvc,
		playbookSvc:               deps.PlaybookSvc,
		gameSvc:                   deps.GameSvc,
	}
}

func (m *RootModel) Init() tea.Cmd {
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

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		var cmd tea.Cmd
		switch m.currentPage {
		case PageNewBattleSelection:
			cmd = m.pageNewBattleSelection.Init()
		case PageGame:
			if page, ok := m.pageGame.(*PageGameModel); ok {
				page.SetGameID(msg.GameID)
			}
			cmd = m.pageGame.Init()
		case PageGameComplete:
			if page, ok := m.pageGameComplete.(*PageGameCompleteModel); ok {
				page.SetGameID(msg.GameID)
			}
			cmd = m.pageGameComplete.Init()
		}
		cmds = append(cmds, cmd)
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
		m.pageGame, cmd = m.pageGame.Update(msg)
		cmds = append(cmds, cmd)
		m.pageGameComplete, cmd = m.pageGameComplete.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
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
	case PageGame:
		m.pageGame, cmd = m.pageGame.Update(msg)
	case PageGameComplete:
		m.pageGameComplete, cmd = m.pageGameComplete.Update(msg)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *RootModel) View() tea.View {
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
	case PageGame:
		return m.pageGame.View()
	case PageGameComplete:
		return m.pageGameComplete.View()
	}

	return tea.NewView("unknown page")
}
