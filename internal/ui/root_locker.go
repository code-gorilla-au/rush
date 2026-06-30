package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

type SubPageLocker int

const (
	SubPageLockerRoom SubPageLocker = iota
	SubPageLockerPlayers
	SubPageLockerPlaybooksList
	SubPageLockerPlaybooksCreate
	SubPageLockerPlaybooksEdit
)

type MsgSwitchLockerPage struct {
	NewPage  SubPageLocker
	Playbook *playbooks.Playbook
	GameID   int64
}

// LockerModel handles all locker room related pages.
type LockerModel struct {
	currentPage                  SubPageLocker
	subPageLocker                tea.Model
	subPageLockerRoom            tea.Model
	subPageLockerPlayers         tea.Model
	subPageLockerPlaybooksList   tea.Model
	subPageLockerPlaybooksCreate tea.Model
	subPageLockerPlaybooksEdit   tea.Model
	globalState                  *GlobalState
	theme                        styles.IceTheme
}

// NewLockerModel returns a new LockerModel.
func NewLockerModel(state *GlobalState, teamsSvc *teams.Service, playbookSvc *playbooks.Service, theme styles.IceTheme) *LockerModel {
	return &LockerModel{
		globalState:                  state,
		theme:                        theme,
		subPageLockerRoom:            NewModelLockerRoom(state, theme),
		subPageLockerPlayers:         NewModelLockerPlayers(state, teamsSvc, theme),
		subPageLockerPlaybooksList:   NewModelLockerPlaybooksList(state, playbookSvc, theme),
		subPageLockerPlaybooksCreate: NewModelLockerPlaybooksCreate(state, playbookSvc, theme),
		subPageLockerPlaybooksEdit:   NewModelLockerPlaybooksEdit(state, playbookSvc, theme),
	}
}

func (m *LockerModel) Init() tea.Cmd {
	return nil
}

func (m *LockerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case MsgSwitchLockerPage:
		switch msg.NewPage {
		case SubPageLockerRoom, SubPageLockerPlayers, SubPageLockerPlaybooksList, SubPageLockerPlaybooksCreate, SubPageLockerPlaybooksEdit:
			m.currentPage = msg.NewPage
		}

	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		m.subPageLockerRoom, cmd = m.subPageLockerRoom.Update(msg)
		cmds = append(cmds, cmd)
		m.subPageLockerPlayers, cmd = m.subPageLockerPlayers.Update(msg)
		cmds = append(cmds, cmd)
		m.subPageLockerPlaybooksList, cmd = m.subPageLockerPlaybooksList.Update(msg)
		cmds = append(cmds, cmd)
		m.subPageLockerPlaybooksCreate, cmd = m.subPageLockerPlaybooksCreate.Update(msg)
		cmds = append(cmds, cmd)
		m.subPageLockerPlaybooksEdit, cmd = m.subPageLockerPlaybooksEdit.Update(msg)
		cmds = append(cmds, cmd)
	}

	var cmd tea.Cmd
	switch m.currentPage {
	case SubPageLockerRoom:
		m.subPageLockerRoom, cmd = m.subPageLockerRoom.Update(msg)
	case SubPageLockerPlayers:
		m.subPageLockerPlayers, cmd = m.subPageLockerPlayers.Update(msg)
	case SubPageLockerPlaybooksList:
		m.subPageLockerPlaybooksList, cmd = m.subPageLockerPlaybooksList.Update(msg)
	case SubPageLockerPlaybooksCreate:
		m.subPageLockerPlaybooksCreate, cmd = m.subPageLockerPlaybooksCreate.Update(msg)
	case SubPageLockerPlaybooksEdit:
		m.subPageLockerPlaybooksEdit, cmd = m.subPageLockerPlaybooksEdit.Update(msg)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *LockerModel) View() tea.View {
	switch m.currentPage {
	case SubPageLockerRoom:
		return m.subPageLockerRoom.View()
	case SubPageLockerPlayers:
		return m.subPageLockerPlayers.View()
	case SubPageLockerPlaybooksList:
		return m.subPageLockerPlaybooksList.View()
	case SubPageLockerPlaybooksCreate:
		return m.subPageLockerPlaybooksCreate.View()
	case SubPageLockerPlaybooksEdit:
		return m.subPageLockerPlaybooksEdit.View()
	}

	return tea.NewView("unknown locker page")
}
