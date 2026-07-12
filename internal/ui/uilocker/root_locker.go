package uilocker

import (
	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"
)

type SubPageLocker int

const (
	SubPageLockerRoom SubPageLocker = iota
	SubPageLockerPlayers
	SubPageLockerTeamStatistics
	SubPageLockerPlaybooksList
	SubPageLockerPlaybooksCreate
	SubPageLockerPlaybooksEdit
	SubPageLockerCoachEdit
)

type MsgSwitchLockerPage struct {
	NewPage  SubPageLocker
	Playbook *playbooks.Playbook
	GameID   int64
}

// LockerModel handles all locker room related pages.
type LockerModel struct {
	currentPage                  SubPageLocker
	subPageLockerRoom            tea.Model
	subPageLockerPlayers         tea.Model
	subPageLockerTeamStatistics  tea.Model
	subPageLockerPlaybooksList   tea.Model
	subPageLockerPlaybooksCreate tea.Model
	subPageLockerPlaybooksEdit   tea.Model
	subPageLockerCoachEdit       tea.Model
}

// NewLockerModel returns a new LockerModel.
func NewLockerModel(state *uistate.GlobalState, teamsSvc *teams.Service, playbookSvc *playbooks.Service, gameSvc *games.Service, theme styles.IceTheme) *LockerModel {
	return &LockerModel{
		subPageLockerRoom:            NewModelLockerRoom(state, theme),
		subPageLockerPlayers:         NewModelLockerPlayers(state, teamsSvc, theme),
		subPageLockerTeamStatistics:  NewModelLockerTeamStatistics(state, gameSvc, theme),
		subPageLockerPlaybooksList:   NewModelLockerPlaybooksList(state, playbookSvc, theme),
		subPageLockerPlaybooksCreate: NewModelLockerPlaybooksCreate(state, playbookSvc, theme),
		subPageLockerPlaybooksEdit:   NewModelLockerPlaybooksEdit(state, playbookSvc, theme),
		subPageLockerCoachEdit:       NewPageLockerCoachEdit(state, teamsSvc, theme),
	}
}

func (m *LockerModel) Init() tea.Cmd {
	return nil
}

func (m *LockerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case uistate.MsgSwitchPage:
		if msg.NewPage == uistate.PageLockerRoom {
			m.currentPage = SubPageLockerRoom
			return m, m.Init()
		}
	case MsgSwitchLockerPage:
		m.currentPage = msg.NewPage
	case tea.WindowSizeMsg:
		return m, m.handleWindowSize(msg)
	}

	return m.updateCurrentPage(msg)
}

func (m *LockerModel) handleWindowSize(msg tea.WindowSizeMsg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.subPageLockerRoom, cmd = m.subPageLockerRoom.Update(msg)
	cmds = append(cmds, cmd)
	m.subPageLockerPlayers, cmd = m.subPageLockerPlayers.Update(msg)
	cmds = append(cmds, cmd)
	m.subPageLockerTeamStatistics, cmd = m.subPageLockerTeamStatistics.Update(msg)
	cmds = append(cmds, cmd)
	m.subPageLockerPlaybooksList, cmd = m.subPageLockerPlaybooksList.Update(msg)
	cmds = append(cmds, cmd)
	m.subPageLockerPlaybooksCreate, cmd = m.subPageLockerPlaybooksCreate.Update(msg)
	cmds = append(cmds, cmd)
	m.subPageLockerPlaybooksEdit, cmd = m.subPageLockerPlaybooksEdit.Update(msg)
	cmds = append(cmds, cmd)
	m.subPageLockerCoachEdit, cmd = m.subPageLockerCoachEdit.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *LockerModel) updateCurrentPage(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.currentPage {
	case SubPageLockerRoom:
		m.subPageLockerRoom, cmd = m.subPageLockerRoom.Update(msg)
	case SubPageLockerPlayers:
		m.subPageLockerPlayers, cmd = m.subPageLockerPlayers.Update(msg)
	case SubPageLockerTeamStatistics:
		m.subPageLockerTeamStatistics, cmd = m.subPageLockerTeamStatistics.Update(msg)
	case SubPageLockerPlaybooksList:
		m.subPageLockerPlaybooksList, cmd = m.subPageLockerPlaybooksList.Update(msg)
	case SubPageLockerPlaybooksCreate:
		m.subPageLockerPlaybooksCreate, cmd = m.subPageLockerPlaybooksCreate.Update(msg)
	case SubPageLockerPlaybooksEdit:
		m.subPageLockerPlaybooksEdit, cmd = m.subPageLockerPlaybooksEdit.Update(msg)
	case SubPageLockerCoachEdit:
		m.subPageLockerCoachEdit, cmd = m.subPageLockerCoachEdit.Update(msg)
	}
	return m, cmd
}

func (m *LockerModel) View() tea.View {
	switch m.currentPage {
	case SubPageLockerRoom:
		return m.subPageLockerRoom.View()
	case SubPageLockerPlayers:
		return m.subPageLockerPlayers.View()
	case SubPageLockerTeamStatistics:
		return m.subPageLockerTeamStatistics.View()
	case SubPageLockerPlaybooksList:
		return m.subPageLockerPlaybooksList.View()
	case SubPageLockerPlaybooksCreate:
		return m.subPageLockerPlaybooksCreate.View()
	case SubPageLockerPlaybooksEdit:
		return m.subPageLockerPlaybooksEdit.View()
	case SubPageLockerCoachEdit:
		return m.subPageLockerCoachEdit.View()
	}

	return tea.NewView("unknown locker page")
}
