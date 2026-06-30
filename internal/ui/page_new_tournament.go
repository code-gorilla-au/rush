package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"
)

type ModelNewTournament struct {
	width       int
	height      int
	theme       styles.IceTheme
	globalState *uistate.GlobalState
}

func NewModelNewTournament(globalState *uistate.GlobalState, theme styles.IceTheme) *ModelNewTournament {
	return &ModelNewTournament{
		globalState: globalState,
		theme:       theme,
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
	case uistate.MsgStateUpdated:
		m.globalState.Coach = vMsg.Coach
		m.globalState.Team = vMsg.Team
	}
	return m, nil
}

func (m *ModelNewTournament) View() tea.View {
	return tea.NewView("new tournament")
}
