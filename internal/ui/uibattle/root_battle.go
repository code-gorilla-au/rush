package uibattle

import (
	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"
)

type SubPageBattle int

const (
	SubPageBattleSelection SubPageBattle = iota
	SubPageBattleConfirm
)

type MsgSwitchBattlePage struct {
	NewPage          SubPageBattle
	SelectedPlaybook *playbooks.Playbook
	SelectedAITeam   *teams.AITeam
}

type BattleModel struct {
	currentPage            SubPageBattle
	subPageBattleSelection tea.Model
	subPageBattleConfirm   tea.Model
}

func NewBattleModel(state *uistate.GlobalState, teamsSvc *teams.Service, playbookSvc *playbooks.Service, gameSvc *games.Service, theme styles.IceTheme) *BattleModel {
	return &BattleModel{
		subPageBattleSelection: NewModelBattleSelection(state, teamsSvc, playbookSvc, gameSvc, theme),
		subPageBattleConfirm:   NewPageBattleConfirm(state, gameSvc, theme),
		currentPage:            SubPageBattleSelection,
	}
}

func (m *BattleModel) Init() tea.Cmd {
	switch m.currentPage {
	case SubPageBattleSelection:
		return m.subPageBattleSelection.Init()
	case SubPageBattleConfirm:
		return m.subPageBattleConfirm.Init()
	default:
		return nil
	}
}

func (m *BattleModel) SetData(playbook *playbooks.Playbook, aiTeam *teams.AITeam) {
	if p, ok := m.subPageBattleConfirm.(*PageBattleConfirmModel); ok {
		p.SetData(playbook, aiTeam)
	}
}

func (m *BattleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case uistate.MsgSwitchPage:
		if msg.NewPage == uistate.PageNewBattleSelection {
			m.currentPage = SubPageBattleSelection
			return m, m.Init()
		}
	case MsgSwitchBattlePage:
		m.currentPage = msg.NewPage
		if msg.SelectedPlaybook != nil && msg.SelectedAITeam != nil {
			m.SetData(msg.SelectedPlaybook, msg.SelectedAITeam)
		}
		return m, m.Init()
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		m.subPageBattleSelection, cmd = m.subPageBattleSelection.Update(msg)
		cmds = append(cmds, cmd)
		m.subPageBattleConfirm, cmd = m.subPageBattleConfirm.Update(msg)
		cmds = append(cmds, cmd)
	}

	var cmd tea.Cmd
	switch m.currentPage {
	case SubPageBattleSelection:
		m.subPageBattleSelection, cmd = m.subPageBattleSelection.Update(msg)
	case SubPageBattleConfirm:
		m.subPageBattleConfirm, cmd = m.subPageBattleConfirm.Update(msg)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *BattleModel) View() tea.View {
	switch m.currentPage {
	case SubPageBattleSelection:
		return m.subPageBattleSelection.View()
	case SubPageBattleConfirm:
		return m.subPageBattleConfirm.View()
	default:
		return tea.NewView("unknown battle page")
	}
}
