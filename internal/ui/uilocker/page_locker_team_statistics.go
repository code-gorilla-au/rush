package uilocker

import (
	"errors"
	"fmt"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/ui/components"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"
)

type lockerTeamStatisticsKeyMap struct {
	uistate.KeyMap
}

func (k lockerTeamStatisticsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Back, k.Quit}
}

func (k lockerTeamStatisticsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.Back, k.Quit},
	}
}

func newLockerTeamStatisticsKeyMap() lockerTeamStatisticsKeyMap {
	return lockerTeamStatisticsKeyMap{KeyMap: uistate.NewKeyMap()}
}

type MsgTeamStatisticsLoaded struct {
	Statistics games.TeamStatistics
	Err        error
}

type ModelLockerTeamStatistics struct {
	width       int
	height      int
	theme       styles.IceTheme
	globalState *uistate.GlobalState
	gameSvc     *games.Service
	keys        lockerTeamStatisticsKeyMap
	footer      components.Footer
	stats       games.TeamStatistics
	loadErr     error
}

func NewModelLockerTeamStatistics(state *uistate.GlobalState, gameSvc *games.Service, theme styles.IceTheme) *ModelLockerTeamStatistics {
	keys := newLockerTeamStatisticsKeyMap()

	return &ModelLockerTeamStatistics{
		theme:       theme,
		globalState: state,
		gameSvc:     gameSvc,
		keys:        keys,
		footer:      components.NewFooter(keys),
	}
}

func (m *ModelLockerTeamStatistics) Init() tea.Cmd {
	return nil
}

func (m *ModelLockerTeamStatistics) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case uistate.MsgStateUpdated:
		m.globalState.Coach = msg.Coach
		m.globalState.Team = msg.Team
	case MsgSwitchLockerPage:
		if msg.NewPage == SubPageLockerTeamStatistics {
			cmds = append(cmds, m.loadStatisticsCmd())
		}
	case MsgTeamStatisticsLoaded:
		m.stats = msg.Statistics
		m.loadErr = msg.Err
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Back):
			return m, func() tea.Msg {
				return MsgSwitchLockerPage{NewPage: SubPageLockerRoom}
			}
		case key.Matches(msg, m.keys.Enter):
			cmds = append(cmds, m.loadStatisticsCmd())
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.footer.Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m *ModelLockerTeamStatistics) loadStatisticsCmd() tea.Cmd {
	if m.globalState.Team == nil {
		return nil
	}
	if m.gameSvc == nil {
		return func() tea.Msg {
			return MsgTeamStatisticsLoaded{Err: errors.New("game service unavailable")}
		}
	}

	teamID := m.globalState.Team.ID
	return func() tea.Msg {
		stats, err := m.gameSvc.GetTeamStatistics(m.globalState.Context(), teamID)
		return MsgTeamStatisticsLoaded{Statistics: stats, Err: err}
	}
}

func (m *ModelLockerTeamStatistics) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true

	teamName := "Unknown Team"
	if m.globalState.Team != nil {
		teamName = m.globalState.Team.Name
	}

	recordTitle := m.theme.Header.Render("Record (W-D-L)")
	record := m.theme.Highlight.Render(fmt.Sprintf("%d - %d - %d", m.stats.Wins, m.stats.Draws, m.stats.Losses))

	statsBody := lipgloss.JoinVertical(
		lipgloss.Left,
		recordTitle,
		record,
		"",
		m.theme.Label.Render(fmt.Sprintf("Games played: %d", m.stats.GamesPlayed)),
		m.theme.Label.Render(fmt.Sprintf("Win rate: %.1f%%", m.stats.WinRate)),
		m.theme.Label.Render(fmt.Sprintf("Rounds won: %d", m.stats.RoundsWon)),
		m.theme.Label.Render(fmt.Sprintf("Rounds lost: %d", m.stats.RoundsLost)),
		m.theme.Label.Render(fmt.Sprintf("Round differential: %+d", m.stats.RoundDifferential)),
		m.theme.Label.Render(fmt.Sprintf("Avg rounds won/game: %.2f", m.stats.AverageRoundsWon)),
		m.theme.Label.Render(fmt.Sprintf("Avg rounds lost/game: %.2f", m.stats.AverageRoundsLost)),
	)

	if m.stats.GamesPlayed == 0 && m.loadErr == nil {
		statsBody = lipgloss.JoinVertical(
			lipgloss.Left,
			statsBody,
			"",
			m.theme.Muted.Render("No completed games yet."),
		)
	}

	if m.loadErr != nil {
		statsBody = lipgloss.JoinVertical(
			lipgloss.Left,
			statsBody,
			"",
			m.theme.Muted.Render(fmt.Sprintf("Could not load team statistics: %v", m.loadErr)),
		)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		m.theme.Logo.Render("TEAM STATISTICS: "+teamName),
		"",
		m.theme.ActiveBorder.Render(statsBody),
		"",
		m.footer.View(m.theme),
	)

	centeredContent := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)

	view.Content = m.theme.Base.
		Width(m.width).
		Height(m.height).
		Render(centeredContent)

	return view
}
