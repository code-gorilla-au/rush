package ui

import (
	"fmt"

	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/components"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type selectionState int

const (
	stateSelectingCoach selectionState = iota
)

type msgAICoachesLoaded struct {
	aiTeams []AITeamItem
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

type AITeamItem struct {
	coach teams.Coach
	team  teams.Team
}

type ModelNewBattleSelection struct {
	width            int
	height           int
	theme            IceTheme
	globalState      *GlobalState
	teamsSvc         *teams.Service
	state            selectionState
	aiCoaches        []AITeamItem
	selectedCoachIdx int
	keys             battleSelectionKeyMap
	footer           components.Footer
	err              error
}

func NewModelNewBattleSelection(globalState *GlobalState, teamsSvc *teams.Service) *ModelNewBattleSelection {
	keys := newBattleSelectionKeyMap()
	return &ModelNewBattleSelection{
		globalState: globalState,
		teamsSvc:    teamsSvc,
		theme:       NewIceTheme(),
		keys:        keys,
		footer:      components.NewFooter(keys),
	}
}

func (m *ModelNewBattleSelection) Init() tea.Cmd {
	return m.loadAICoaches
}

func (m *ModelNewBattleSelection) loadAICoaches() tea.Msg {
	aiTeams, err := m.teamsSvc.ListAITeams(m.globalState.Context())
	if err != nil {
		return err
	}

	items := make([]AITeamItem, 0, len(aiTeams))
	for _, aiTeam := range aiTeams {
		items = append(items, AITeamItem{
			coach: aiTeam.Coach,
			team:  aiTeam.Team,
		})
	}

	return msgAICoachesLoaded{aiTeams: items}
}

func (m *ModelNewBattleSelection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	m.footer.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case msgAICoachesLoaded:
		m.aiCoaches = msg.aiTeams
		if len(m.aiCoaches) > 0 {
			m.selectedCoachIdx = 0
		}
	case error:
		m.err = msg
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Back):
			return m, func() tea.Msg {
				return MsgSwitchPage{NewPage: PageTitle}
			}
		case key.Matches(msg, m.keys.Select):
			if len(m.aiCoaches) > 0 {
				// TODO: Start battle with selected coach
				return m, func() tea.Msg {
					return MsgSwitchPage{NewPage: PageTitle}
				}
			}
		case msg.String() == "up", msg.String() == "k":
			if m.selectedCoachIdx > 0 {
				m.selectedCoachIdx--
			}
		case msg.String() == "down", msg.String() == "j":
			if m.selectedCoachIdx < len(m.aiCoaches)-1 {
				m.selectedCoachIdx++
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *ModelNewBattleSelection) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Initializing...")
	}

	view := tea.NewView("")
	view.AltScreen = true

	var content string
	if m.err != nil {
		content = m.theme.Logo.Render(fmt.Sprintf("Error: %v", m.err))
	} else {
		content = m.viewCoaches()
	}

	centeredContent := lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center,
			m.theme.Logo.Render("NEW BATTLE"),
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

func (m *ModelNewBattleSelection) viewCoaches() string {
	if len(m.aiCoaches) == 0 {
		return "No AI coaches available"
	}

	var s string
	s += m.theme.Logo.Render("Select your opponent:") + "\n\n"

	for i, item := range m.aiCoaches {
		avatar := components.NewCoachAvatar(&item.coach, &item.team)
		avatarView := avatar.View(m.theme.CoachTeam, m.theme.CoachName)
		if i == m.selectedCoachIdx {
			s += m.theme.ListSelected.Render("> " + avatarView)
		} else {
			s += "  " + avatarView
		}
		s += "\n\n"
	}

	return s
}
