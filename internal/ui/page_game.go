package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/ui/components"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

type PageGameModel struct {
	width       int
	height      int
	theme       styles.IceTheme
	globalState *GlobalState
	gameSvc     *games.Service
	gameID      int64
	game        *games.Game
	gameComp    components.Game
}

func NewModelGame(state *GlobalState, gameSvc *games.Service, theme styles.IceTheme) *PageGameModel {
	return &PageGameModel{
		theme:       theme,
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
		m.game = new(games.Game)
		*m.game = msg.Game
		teamA, teamB := getTeamNames(m.game.Name())
		m.gameComp = components.NewGame(m.game, teamA, teamB, nil)
		cmds = append(cmds, m.gameComp.Init())
	case components.MsgResolveRound:
		cmds = append(cmds, m.gameComp.Update(msg))
		cmds = append(cmds, func() tea.Msg {
			_, err := m.gameSvc.UpdateGame(m.globalState.Context(), *m.game)
			if err != nil {
				return MsgGameError{Err: err}
			}
			return nil
		})
	case components.MsgNextRound:
		if !m.game.IsGameComplete() {
			teamA, teamB := getTeamNames(m.game.Name())
			m.gameComp = components.NewGame(m.game, teamA, teamB, nil)
			cmds = append(cmds, m.gameComp.Init())
		} else {
			cmds = append(cmds, func() tea.Msg {
				_, err := m.gameSvc.CompleteGame(m.globalState.Context(), *m.game)
				if err != nil {
					return MsgGameError{Err: err}
				}
				return MsgSwitchPage{
					NewPage: PageGameComplete,
					GameID:  m.game.ID(),
				}
			})
		}
	default:
		cmds = append(cmds, m.gameComp.Update(msg))
	}

	return m, tea.Batch(cmds...)
}

func getTeamNames(gameName string) (string, string) {
	teamA, teamB := "Team A", "Team B"
	parts := strings.Split(gameName, " VS ")
	if len(parts) == 2 {
		teamA = parts[0]
		teamB = parts[1]
	}
	return teamA, teamB
}

func (m *PageGameModel) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true

	var mainContent string
	if m.game == nil {
		mainContent = "Loading Game..."
	} else {
		mainContent = m.gameComp.View(m.theme)
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
