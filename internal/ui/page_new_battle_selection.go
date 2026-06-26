package ui

import (
	"fmt"

	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/components"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type selectionState int

const (
	stateSelectingCoach selectionState = iota
	stateSelectingPlaybook
)

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

type AICoachItem struct {
	coach teams.Coach
	team  teams.Team
}

type ModelNewBattleSelection struct {
	width            int
	height           int
	theme            IceTheme
	globalState      *GlobalState
	teamsSvc         *teams.Service
	playbookSvc      *playbooks.Service
	state            selectionState
	aiCoaches        []AICoachItem
	selectedCoachIdx int
	playbooks        []playbooks.Playbook
	playbookList     components.PlaybookList
	keys             battleSelectionKeyMap
	footer           components.Footer
	err              error
}

func NewModelNewBattleSelection(globalState *GlobalState, teamsSvc *teams.Service, playbookSvc *playbooks.Service) *ModelNewBattleSelection {
	keys := newBattleSelectionKeyMap()
	return &ModelNewBattleSelection{
		globalState:  globalState,
		teamsSvc:     teamsSvc,
		playbookSvc:  playbookSvc,
		theme:        NewIceTheme(),
		keys:         keys,
		footer:       components.NewFooter(keys),
		playbookList: components.NewPlaybookList(nil),
	}
}

func (m *ModelNewBattleSelection) Init() tea.Cmd {
	return m.loadAICoaches
}

func (m *ModelNewBattleSelection) loadAICoaches() tea.Msg {
	coaches, err := m.teamsSvc.ListAICoaches(m.globalState.Context())
	if err != nil {
		return err
	}

	items := make([]AICoachItem, 0, len(coaches))
	for _, coach := range coaches {
		team, err := m.teamsSvc.GetTeamByCoachID(m.globalState.Context(), coach.ID)
		if err != nil {
			continue
		}
		items = append(items, AICoachItem{coach: coach, team: team})
	}
	return msgAICoachesLoaded{coaches: items}
}

type msgAICoachesLoaded struct {
	coaches []AICoachItem
}

type msgPlaybooksLoaded struct {
	playbooks []playbooks.Playbook
}

func (m *ModelNewBattleSelection) loadPlaybooks(teamID int64) tea.Cmd {
	return func() tea.Msg {
		pbks, err := m.playbookSvc.GetTeamPlaybooks(m.globalState.Context(), teamID)
		if err != nil {
			return err
		}
		return msgPlaybooksLoaded{playbooks: pbks}
	}
}

func (m *ModelNewBattleSelection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.footer.Update(msg)
		m.playbookList.SetSize(m.width, m.height-10)
	case msgAICoachesLoaded:
		m.aiCoaches = msg.coaches
	case msgPlaybooksLoaded:
		m.playbooks = msg.playbooks
		cmds = append(cmds, m.playbookList.SetItems(msg.playbooks))
	case error:
		m.err = msg
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Back):
			if m.state == stateSelectingPlaybook {
				m.state = stateSelectingCoach
				return m, nil
			}
			return m, func() tea.Msg {
				return MsgSwitchPage{NewPage: PageTitle}
			}
		case key.Matches(msg, m.keys.Select):
			if m.state == stateSelectingCoach && len(m.aiCoaches) > 0 {
				m.state = stateSelectingPlaybook
				cmds = append(cmds, m.loadPlaybooks(m.aiCoaches[m.selectedCoachIdx].team.ID))
			} else if m.state == stateSelectingPlaybook {
				selectedPlaybook := m.playbookList.SelectedItem()
				if selectedPlaybook != nil {
					// TODO: Start battle with selected coach and playbook
					return m, func() tea.Msg {
						return MsgSwitchPage{NewPage: PageTitle}
					}
				}
			}
		case msg.String() == "up", msg.String() == "k":
			if m.state == stateSelectingCoach && m.selectedCoachIdx > 0 {
				m.selectedCoachIdx--
			}
		case msg.String() == "down", msg.String() == "j":
			if m.state == stateSelectingCoach && m.selectedCoachIdx < len(m.aiCoaches)-1 {
				m.selectedCoachIdx++
			}
		}
	}

	if m.state == stateSelectingPlaybook {
		var cmd tea.Cmd
		m.playbookList, cmd = m.playbookList.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *ModelNewBattleSelection) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true

	var content string
	if m.err != nil {
		content = m.theme.Logo.Render(fmt.Sprintf("Error: %v", m.err))
	} else if m.state == stateSelectingCoach {
		content = m.viewCoaches()
	} else {
		content = m.viewPlaybooks()
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

func (m *ModelNewBattleSelection) viewPlaybooks() string {
	if len(m.playbooks) == 0 {
		return "This coach has no playbooks"
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		m.theme.Logo.Render(fmt.Sprintf("Select %s's playbook:", m.aiCoaches[m.selectedCoachIdx].coach.Name)),
		"",
		m.playbookList.View(),
	)
}
