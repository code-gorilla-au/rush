package ui

import (
	"fmt"

	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/components"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type msgDataLoaded struct {
	playbooks []playbooks.Playbook
	aiTeams   []teams.AITeam
}

type battleSelectionKeyMap struct {
	components.CommonKeys
	Back   key.Binding
	Select key.Binding
}

func (k battleSelectionKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Select, k.Back, k.Quit}
}

func (k battleSelectionKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Select, k.Back, k.Quit},
	}
}

func newBattleSelectionKeyMap() battleSelectionKeyMap {
	return battleSelectionKeyMap{
		CommonKeys: components.NewCommonKeys(),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}

type selectionState int

const (
	stateSelectingPlaybook selectionState = iota
	stateSelectingOpponent
	stateConfirming
)

type ModelNewBattleSelection struct {
	width            int
	height           int
	theme            IceTheme
	globalState      *GlobalState
	teamsSvc         *teams.Service
	playbookSvc      *playbooks.Service
	gameSvc          *games.Service
	state            selectionState
	playbookList     components.PlaybookList
	aiTeamList       components.AITeamList
	selectedPlaybook *playbooks.Playbook
	selectedAITeam   *teams.AITeam
	keys             battleSelectionKeyMap
	footer           components.Footer
	err              error
}

func NewModelNewBattleSelection(globalState *GlobalState, teamsSvc *teams.Service, playbookSvc *playbooks.Service, gameSvc *games.Service) *ModelNewBattleSelection {
	keys := newBattleSelectionKeyMap()
	return &ModelNewBattleSelection{
		globalState:  globalState,
		teamsSvc:     teamsSvc,
		playbookSvc:  playbookSvc,
		gameSvc:      gameSvc,
		theme:        NewIceTheme(),
		keys:         keys,
		footer:       components.NewFooter(keys),
		playbookList: components.NewPlaybookList(nil),
		aiTeamList:   components.NewAITeamList(nil),
	}
}

func (m *ModelNewBattleSelection) Init() tea.Cmd {
	m.reset()
	return m.loadData
}

func (m *ModelNewBattleSelection) loadData() tea.Msg {
	if m.globalState.Team == nil {
		return fmt.Errorf("no team loaded")
	}

	playbooks, err := m.playbookSvc.GetTeamPlaybooks(m.globalState.Context(), m.globalState.Team.ID)
	if err != nil {
		return err
	}

	aiTeams, err := m.teamsSvc.ListAITeams(m.globalState.Context())
	if err != nil {
		return err
	}

	return msgDataLoaded{
		playbooks: playbooks,
		aiTeams:   aiTeams,
	}
}

func (m *ModelNewBattleSelection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.footer.Update(msg)
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.playbookList.SetSize(m.width/2-4, m.height-20)
		m.aiTeamList.SetSize(m.width/2-4, m.height-20)
	case msgDataLoaded:
		m.playbookList.SetItems(msg.playbooks)
		m.aiTeamList.SetItems(msg.aiTeams)
	case error:
		m.err = msg
	case tea.KeyMsg:
		model, cmd := m.handleKey(msg)
		if cmd != nil {
			return model, cmd
		}
	}

	m.playbookList.SetActive(m.state == stateSelectingPlaybook)
	m.aiTeamList.SetActive(m.state == stateSelectingOpponent)

	var listCmd tea.Cmd
	if m.state == stateSelectingPlaybook {
		m.playbookList, listCmd = m.playbookList.Update(msg)
	} else if m.state == stateSelectingOpponent {
		m.aiTeamList, listCmd = m.aiTeamList.Update(msg)
	}
	cmds = append(cmds, listCmd)

	return m, tea.Batch(cmds...)
}

func (m *ModelNewBattleSelection) handleKey(msg tea.KeyMsg) (*ModelNewBattleSelection, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.Back):
		if m.state == stateSelectingOpponent {
			m.state = stateSelectingPlaybook
			return m, nil
		}
		if m.state == stateConfirming {
			m.state = stateSelectingOpponent
			return m, nil
		}
		return m, func() tea.Msg {
			return MsgSwitchPage{NewPage: PageTitle}
		}
	case key.Matches(msg, m.keys.Select):
		if m.state == stateSelectingPlaybook {
			m.selectedPlaybook = m.playbookList.SelectedItem()
			if m.selectedPlaybook != nil {
				m.state = stateSelectingOpponent
			}
			return m, nil
		}
		if m.state == stateSelectingOpponent {
			m.selectedAITeam = m.aiTeamList.SelectedAITeam()
			if m.selectedAITeam != nil {
				m.state = stateConfirming
			}
			return m, nil
		}
		if m.state == stateConfirming {
			return m, m.createGame
		}
	}
	return m, nil
}

func (m *ModelNewBattleSelection) createGame() tea.Msg {
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

	m.reset()

	return MsgSwitchPage{
		NewPage: PageGame,
		GameID:  game.ID(),
	}
}

func (m *ModelNewBattleSelection) reset() {
	m.state = stateSelectingPlaybook
	m.selectedPlaybook = nil
	m.selectedAITeam = nil
	m.err = nil
	m.playbookList.Reset()
	m.aiTeamList.Reset()
}

func (m *ModelNewBattleSelection) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Initializing...")
	}

	view := tea.NewView("")
	view.AltScreen = true

	var content string
	var header string
	if m.err != nil {
		content = m.theme.Logo.Render(fmt.Sprintf("Error: %v", m.err))
		header = "ERROR"
	} else if m.state == stateConfirming {
		content = m.viewConfirmation()
		header = "STEP 3: CONFIRMATION"
	} else {
		content = m.viewSelection()
		if m.state == stateSelectingPlaybook {
			header = "STEP 1: SELECT YOUR PLAYBOOK"
		} else {
			header = "STEP 2: SELECT YOUR OPPONENT"
		}
	}

	centeredContent := lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center,
			m.theme.Logo.Render("NEW BATTLE"),
			m.theme.Footer.Render(header),
			"",
			content,
			"",
			m.footer.View(m.theme.Footer),
		),
	)

	view.Content = m.theme.Base.
		Width(m.width).
		Height(m.height).
		Render(centeredContent)

	return view
}

func (m *ModelNewBattleSelection) viewSelection() string {
	return lipgloss.JoinHorizontal(lipgloss.Top,
		m.playbookList.View(),
		lipgloss.NewStyle().Width(2).Render(""),
		m.aiTeamList.View(),
	)
}

func (m *ModelNewBattleSelection) viewConfirmation() string {
	return lipgloss.JoinVertical(lipgloss.Center,
		m.theme.Logo.Render("CONFIRM BATTLE"),
		"",
		fmt.Sprintf("Your Playbook: %s", m.selectedPlaybook.Name),
		fmt.Sprintf("Opponent: %s", m.selectedAITeam.Team.Name),
		"",
		m.theme.ListSelected.Render("Press ENTER to start the game"),
		"Press ESC to go back",
	)
}
