package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/games"
)

type PageGameModel struct {
	width       int
	height      int
	theme       IceTheme
	globalState *GlobalState
	gameSvc     *games.Service
	gameID      int64
}

func NewModelGame(state *GlobalState, gameSvc *games.Service) *PageGameModel {
	return &PageGameModel{
		theme:       NewIceTheme(),
		globalState: state,
		gameSvc:     gameSvc,
	}
}

func (m *PageGameModel) SetGameID(id int64) {
	m.gameID = id
}

func (m *PageGameModel) Init() tea.Cmd {
	return nil
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
	}

	return m, tea.Batch(cmds...)
}

func (m *PageGameModel) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true

	mainContent := "Game"

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
