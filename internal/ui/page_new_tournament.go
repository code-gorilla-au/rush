package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

type ModelNewTournament struct {
	width       int
	height      int
	theme       styles.IceTheme
	globalState *GlobalState
}

func NewModelNewTournament(globalState *GlobalState) *ModelNewTournament {
	return &ModelNewTournament{
		globalState: globalState,
	}
}

func (m *ModelNewTournament) Init() tea.Cmd {
	return nil
}

func (m *ModelNewTournament) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *ModelNewTournament) View() tea.View {
	return tea.NewView("new tournament")
}
