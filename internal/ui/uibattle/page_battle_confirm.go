package uibattle

import (
	"fmt"

	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/components"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type PageBattleConfirmModel struct {
	width            int
	height           int
	theme            styles.IceTheme
	globalState      *uistate.GlobalState
	gameSvc          *games.Service
	selectedPlaybook *playbooks.Playbook
	selectedAITeam   *teams.AITeam
	keys             battleSelectionKeyMap
	footer           components.Footer
	err              error
}

func NewPageBattleConfirm(globalState *uistate.GlobalState, gameSvc *games.Service, theme styles.IceTheme) *PageBattleConfirmModel {
	keys := newBattleSelectionKeyMap()
	return &PageBattleConfirmModel{
		globalState: globalState,
		gameSvc:     gameSvc,
		theme:       theme,
		keys:        keys,
		footer:      components.NewFooter(keys),
	}
}

func (m *PageBattleConfirmModel) SetData(playbook *playbooks.Playbook, aiTeam *teams.AITeam) {
	m.selectedPlaybook = playbook
	m.selectedAITeam = aiTeam
}

func (m *PageBattleConfirmModel) Init() tea.Cmd {
	return nil
}

func (m *PageBattleConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.footer.Update(msg)
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case error:
		m.err = msg
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Back):
			return m, func() tea.Msg {
				return MsgSwitchBattlePage{NewPage: SubPageBattleSelection}
			}
		case key.Matches(msg, m.keys.Select):
			return m, m.createGame
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *PageBattleConfirmModel) createGame() tea.Msg {
	if m.selectedPlaybook == nil || m.selectedAITeam == nil {
		return fmt.Errorf("missing playbook or opponent")
	}

	params := games.NewGameParams{
		TeamA: games.TeamConfig{
			TeamID:     m.globalState.Team.ID,
			TeamName:   m.globalState.Team.Name,
			Formations: m.selectedPlaybook.Formations,
		},
		TeamB: games.TeamConfig{
			TeamID:     m.selectedAITeam.Team.ID,
			TeamName:   m.selectedAITeam.Team.Name,
			Formations: m.selectedAITeam.Playbook.Formations,
		},
	}

	game, err := m.gameSvc.NewGame(m.globalState.Context(), params)
	if err != nil {
		return err
	}

	return uistate.MsgSwitchPage{
		NewPage: uistate.PageGame,
		GameID:  game.ID(),
	}
}

func (m *PageBattleConfirmModel) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Initializing...")
	}

	view := tea.NewView("")
	view.AltScreen = true

	var content string
	header := "STEP 3: CONFIRMATION"
	if m.err != nil {
		content = m.theme.Logo.Render(fmt.Sprintf("Error: %v", m.err))
		header = "ERROR"
	} else {
		content = lipgloss.JoinVertical(lipgloss.Center,
			m.theme.Logo.Render("CONFIRM BATTLE"),
			"",
			fmt.Sprintf("Your Playbook: %s", m.selectedPlaybook.Name),
			fmt.Sprintf("Opponent: %s", m.selectedAITeam.Team.Name),
			"",
			m.theme.ListSelected.Render("Press ENTER to start the game"),
			"Press ESC to go back",
		)
	}

	centeredContent := lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center,
			m.theme.Logo.Render("NEW BATTLE"),
			m.theme.Muted.Render(header),
			"",
			content,
			"",
			m.footer.View(m.theme),
		),
	)

	view.Content = m.theme.Base.
		Width(m.width).
		Height(m.height).
		Render(centeredContent)

	return view
}
