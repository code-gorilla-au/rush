package ui

import tea "charm.land/bubbletea/v2"

type ModelNewBattle struct {
	width       int
	height      int
	theme       IceTheme
	globalState *GlobalState
}

func NewModelNewBattle(globalState *GlobalState) *ModelNewBattle {
	return &ModelNewBattle{
		globalState: globalState,
	}
}

func (m *ModelNewBattle) Init() tea.Cmd {
	return nil
}

func (m *ModelNewBattle) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch vMsg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = vMsg.Width
		m.height = vMsg.Height
	case MsgStateUpdated:
		m.globalState.Coach = vMsg.Coach
		m.globalState.Team = vMsg.Team
	}

	return m, nil
}

func (m *ModelNewBattle) View() tea.View {
	return tea.NewView("new battle")
}
