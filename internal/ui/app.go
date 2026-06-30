package ui

import (
	"context"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uilocker"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"
)

type RootModel struct {
	ctx                    context.Context
	width                  int
	height                 int
	theme                  styles.IceTheme
	currentPage            uistate.Page
	pageTitle              tea.Model
	pageCreateCoach        tea.Model
	pageLocker             tea.Model
	pageNewTournament      tea.Model
	pageNewBattleSelection tea.Model
	pageTitleSettings      tea.Model
	pageGame               tea.Model
	pageGameComplete       tea.Model
	globalState            *uistate.GlobalState
	teamsSvc               *teams.Service
	playbookSvc            *playbooks.Service
	gameSvc                *games.Service
}

type Dependencies struct {
	TeamsSvc    *teams.Service
	PlaybookSvc *playbooks.Service
	GameSvc     *games.Service
}

// New returns a new UI model.
func New(deps Dependencies) *RootModel {
	state := &uistate.GlobalState{}
	theme := styles.NewIceTheme()

	return &RootModel{
		ctx:                    context.Background(),
		theme:                  theme,
		currentPage:            uistate.PageTitle,
		pageTitle:              NewModelTitle(state, theme),
		pageCreateCoach:        NewModelCreateCoach(state, deps.TeamsSvc, theme),
		pageLocker:             uilocker.NewLockerModel(state, deps.TeamsSvc, deps.PlaybookSvc, theme),
		pageNewTournament:      NewModelNewTournament(state, theme),
		pageNewBattleSelection: NewModelNewBattleSelection(state, deps.TeamsSvc, deps.PlaybookSvc, deps.GameSvc, theme),
		pageTitleSettings:      NewModelTitleSettings(state, theme),
		pageGame:               NewModelGame(state, deps.GameSvc, theme),
		pageGameComplete:       NewPageGameComplete(state, deps.TeamsSvc, deps.GameSvc, theme),
		globalState:            state,
		teamsSvc:               deps.TeamsSvc,
		playbookSvc:            deps.PlaybookSvc,
		gameSvc:                deps.GameSvc,
	}
}

func (m *RootModel) Init() tea.Cmd {
	return func() tea.Msg {
		coach, err := m.teamsSvc.GetDefaultCoach(m.ctx)
		if err != nil {
			return uistate.MsgStateUpdated{Coach: nil}
		}

		team, err := m.teamsSvc.GetTeamByCoachID(m.ctx, coach.ID)
		if err != nil {
			return uistate.MsgStateUpdated{Coach: &coach, Team: nil}
		}

		return uistate.MsgStateUpdated{Coach: &coach, Team: &team}
	}
}

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case uistate.MsgStateUpdated:
		m.globalState.Coach = msg.Coach
		m.globalState.Team = msg.Team
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case uistate.MsgSwitchPage:
		m.currentPage = msg.NewPage
		var cmd tea.Cmd
		switch m.currentPage {
		case uistate.PageNewBattleSelection:
			cmd = m.pageNewBattleSelection.Init()
		case uistate.PageGame:
			if page, ok := m.pageGame.(*PageGameModel); ok {
				page.SetGameID(msg.GameID)
			}
			cmd = m.pageGame.Init()
		case uistate.PageGameComplete:
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
		m.pageLocker, cmd = m.pageLocker.Update(msg)
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
	case uistate.PageTitle:
		m.pageTitle, cmd = m.pageTitle.Update(msg)
	case uistate.PageCreateCoach:
		m.pageCreateCoach, cmd = m.pageCreateCoach.Update(msg)
	case uistate.PageLockerRoom:
		m.pageLocker, cmd = m.pageLocker.Update(msg)
	case uistate.PageNewTournament:
		m.pageNewTournament, cmd = m.pageNewTournament.Update(msg)
	case uistate.PageNewBattleSelection:
		m.pageNewBattleSelection, cmd = m.pageNewBattleSelection.Update(msg)
	case uistate.PageTitleSettings:
		m.pageTitleSettings, cmd = m.pageTitleSettings.Update(msg)
	case uistate.PageGame:
		m.pageGame, cmd = m.pageGame.Update(msg)
	case uistate.PageGameComplete:
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
	case uistate.PageTitle:
		return m.pageTitle.View()
	case uistate.PageCreateCoach:
		return m.pageCreateCoach.View()
	case uistate.PageLockerRoom:
		return m.pageLocker.View()
	case uistate.PageNewTournament:
		return m.pageNewTournament.View()
	case uistate.PageNewBattleSelection:
		return m.pageNewBattleSelection.View()
	case uistate.PageTitleSettings:
		return m.pageTitleSettings.View()
	case uistate.PageGame:
		return m.pageGame.View()
	case uistate.PageGameComplete:
		return m.pageGameComplete.View()
	}

	return tea.NewView("unknown page")
}
