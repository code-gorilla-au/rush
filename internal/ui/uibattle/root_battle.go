package uibattle

import (
	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"
)

type BattleModel struct {
	subPageBattleSelection tea.Model
}

func NewBattleModel(state *uistate.GlobalState, teamsSvc *teams.Service, playbookSvc *playbooks.Service, gameSvc *games.Service, theme styles.IceTheme) *BattleModel {
	return &BattleModel{
		subPageBattleSelection: NewModelBattleSelection(state, teamsSvc, playbookSvc, gameSvc, theme),
	}
}

func (m *BattleModel) Init() tea.Cmd {
	return m.subPageBattleSelection.Init()
}

func (m *BattleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case uistate.MsgSwitchPage:
		if msg.NewPage == uistate.PageNewBattleSelection {
			return m, m.Init()
		}
	}

	var cmd tea.Cmd
	m.subPageBattleSelection, cmd = m.subPageBattleSelection.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *BattleModel) View() tea.View {
	return m.subPageBattleSelection.View()
}
