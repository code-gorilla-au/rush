package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/teams"
)

type ModelLockerPlayers struct {
	width       int
	height      int
	theme       IceTheme
	globalState *GlobalState
	teamsSvc    *teams.Service
}

func NewModelLockerPlayers(state *GlobalState, teamsSvc *teams.Service) *ModelLockerPlayers {
	return &ModelLockerPlayers{
		theme:       NewIceTheme(),
		globalState: state,
		teamsSvc:    teamsSvc,
	}
}

func (m *ModelLockerPlayers) Init() tea.Cmd {
	return nil
}

func (m *ModelLockerPlayers) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m *ModelLockerPlayers) View() tea.View {
	return tea.NewView("locker players")
}
