package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/components"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"
)

type PageGameCompleteModel struct {
	width       int
	height      int
	theme       styles.IceTheme
	globalState *uistate.GlobalState
	teamsSvc    *teams.Service
	gameSvc     *games.Service
	gameID      int64
	winnerTeam  *teams.Team
	winnerCoach *teams.Coach
	isDraw      bool
	err         error
}

func NewPageGameComplete(state *uistate.GlobalState, teamsSvc *teams.Service, gameSvc *games.Service, theme styles.IceTheme) *PageGameCompleteModel {
	return &PageGameCompleteModel{
		theme:       theme,
		globalState: state,
		teamsSvc:    teamsSvc,
		gameSvc:     gameSvc,
	}
}

func (m *PageGameCompleteModel) SetGameID(id int64) {
	m.gameID = id
}

type MsgWinnerLoaded struct {
	Team   *teams.Team
	Coach  *teams.Coach
	IsDraw bool
}

func (m *PageGameCompleteModel) Init() tea.Cmd {
	return func() tea.Msg {
		ctx := m.globalState.Context()
		game, err := m.gameSvc.GetGame(ctx, m.gameID)
		if err != nil {
			return MsgGameError{Err: err}
		}

		winnerID, err := game.CalculateWinner()
		if err != nil {
			return MsgGameError{Err: err}
		}

		if winnerID == 0 {
			return MsgWinnerLoaded{IsDraw: true}
		}

		team, err := m.teamsSvc.GetTeamByID(ctx, winnerID)
		if err != nil {
			return MsgGameError{Err: err}
		}

		coach, err := m.teamsSvc.GetCoachByID(ctx, int64(team.CoachID))
		if err != nil {
			return MsgGameError{Err: err}
		}

		return MsgWinnerLoaded{Team: &team, Coach: &coach}
	}
}

func (m *PageGameCompleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case MsgWinnerLoaded:
		m.winnerTeam = msg.Team
		m.winnerCoach = msg.Coach
		m.isDraw = msg.IsDraw
	case MsgGameError:
		m.err = msg.Err
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, func() tea.Msg {
				return uistate.MsgSwitchPage{NewPage: uistate.PageTitle}
			}
		}
	}
	return m, nil
}

func (m *PageGameCompleteModel) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true

	var content string
	if m.err != nil {
		content = fmt.Sprintf("Error: %v", m.err)
		view.Content = m.theme.Base.
			Width(m.width).
			Height(m.height).
			Render(content)
		return view
	}

	if m.winnerTeam == nil && !m.isDraw {
		content = "Loading Results..."
	} else {
		content = m.renderMainContent(content)
	}

	centered := lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)

	view.Content = m.theme.Base.
		Width(m.width).
		Height(m.height).
		Render(centered)

	return view
}

func (m *PageGameCompleteModel) renderMainContent(content string) string {
	winMsg := m.theme.Logo.
		Padding(1, 0).
		Render("🏆 GAME COMPLETE 🏆")

	var mainContent string
	if m.isDraw {
		mainContent = m.theme.SecondaryHeader.
			Render("It's a DRAW!")
	} else {
		if m.winnerCoach.IsHuman {
			mainContent = components.NewCoachWinnerHuman(m.winnerTeam, m.winnerCoach).View(m.theme)
		} else {
			mainContent = components.NewCoachWinnerAI(m.winnerTeam, m.winnerCoach).View(m.theme)
		}
	}

	footer := m.theme.Footer.Render("\nPress Enter to continue")

	content = lipgloss.JoinVertical(
		lipgloss.Center,
		winMsg,
		mainContent,
		footer,
	)
	return content
}
