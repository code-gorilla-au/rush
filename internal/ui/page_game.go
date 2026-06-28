package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/ui/components"
)

type PageGameModel struct {
	width       int
	height      int
	theme       IceTheme
	globalState *GlobalState
	gameSvc     *games.Service
	gameID      int64
	game        *games.Game
	roundComp   components.Round
}

func NewModelGame(state *GlobalState, gameSvc *games.Service) *PageGameModel {
	return &PageGameModel{
		theme:       NewIceTheme(),
		globalState: state,
		gameSvc:     gameSvc,
	}
}

type MsgGameLoaded struct {
	Game games.Game
}

type MsgGameError struct {
	Err error
}

func (m *PageGameModel) SetGameID(id int64) {
	m.gameID = id
}

func (m *PageGameModel) Init() tea.Cmd {
	return func() tea.Msg {
		game, err := m.gameSvc.GetGame(m.globalState.Context(), m.gameID)
		if err != nil {
			return MsgGameError{Err: err}
		}
		return MsgGameLoaded{Game: game}
	}
}

func (m *PageGameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case MsgStateUpdated:
		m.globalState.Coach = msg.Coach
		m.globalState.Team = msg.Team
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case MsgGameLoaded:
		m.game = &msg.Game
		teamA, teamB := "Team A", "Team B"
		parts := strings.Split(m.game.Name(), " VS ")
		if len(parts) == 2 {
			teamA = parts[0]
			teamB = parts[1]
		}
		m.roundComp = components.NewRound(m.game.Rounds()[m.game.CurrentRound()], teamA, teamB)
	}

	return m, tea.Batch(cmds...)
}

func (m *PageGameModel) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true

	var mainContent string
	if m.game == nil {
		mainContent = "Loading Game..."
	} else {
		mainContent = m.roundComp.View()
	}

	centeredContent := lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		mainContent,
	)

	view.Content = m.theme.Base.
		Width(m.width).
		Height(m.height).
		Render(centeredContent)

	return view
}
