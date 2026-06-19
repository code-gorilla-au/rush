package ui

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/teams"
)

type MsgCoachCreated struct {
	Coach teams.Coach
	Team  teams.Team
}

type ModelCreateCoach struct {
	width       int
	height      int
	theme       IceTheme
	globalState *GlobalState
	teamsSvc    *teams.Service

	coachInput textinput.Model
	teamInput  textinput.Model
	focusIndex int
	err        error
}

func NewModelCreateCoach(state *GlobalState, teamsSvc *teams.Service) *ModelCreateCoach {
	c := textinput.New()
	c.Placeholder = "Coach Name"
	c.Focus()
	c.CharLimit = 156
	c.SetWidth(20)

	t := textinput.New()
	t.Placeholder = "Team Name"
	t.CharLimit = 156
	c.SetWidth(20)

	return &ModelCreateCoach{
		globalState: state,
		teamsSvc:    teamsSvc,
		coachInput:  c,
		teamInput:   t,
		theme:       NewIceTheme(),
	}
}

func (m *ModelCreateCoach) Init() tea.Cmd {
	return func() tea.Msg { return textinput.Blink() }
}

func (m *ModelCreateCoach) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "up", "down":
			s := msg.String()

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > 1 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = 1
			}

			m.updateFocus()

		case "enter":
			if m.focusIndex == 1 {
				return m, m.submit()
			}
			m.focusIndex++
			m.updateFocus()

		case "esc":
			return m, func() tea.Msg {
				return MsgSwitchPage{NewPage: PageTitle}
			}
		}

	case MsgCoachCreated:
		m.globalState.Coach = &msg.Coach
		m.globalState.Team = &msg.Team
		return m, func() tea.Msg {
			return MsgSwitchPage{NewPage: PageLockerRoom}
		}
	}

	var cmd tea.Cmd
	m.coachInput, cmd = m.coachInput.Update(msg)
	cmds = append(cmds, cmd)

	m.teamInput, cmd = m.teamInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *ModelCreateCoach) updateFocus() {
	if m.focusIndex == 0 {
		m.coachInput.Focus()
		m.teamInput.Blur()
	} else {
		m.coachInput.Blur()
		m.teamInput.Focus()
	}
}

func (m *ModelCreateCoach) submit() tea.Cmd {
	return func() tea.Msg {
		ctx := m.globalState.Context()
		coach, err := m.teamsSvc.CreateCoach(ctx, m.coachInput.Value(), true)
		if err != nil {
			return err
		}

		team, err := m.teamsSvc.CreateTeam(ctx, m.teamInput.Value(), coach.ID, true)
		if err != nil {
			return err
		}

		return MsgCoachCreated{
			Coach: coach,
			Team:  team,
		}
	}
}

func (m *ModelCreateCoach) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true

	s := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center)

	form := lipgloss.JoinVertical(
		lipgloss.Left,
		m.theme.Logo.Render("RUSH - NEW CAREER"),
		"",
		"Coach Details",
		m.coachInput.View(),
		"",
		"Team Details",
		m.teamInput.View(),
		"",
		m.theme.Hotkey.Render("enter")+" continue • "+m.theme.Hotkey.Render("esc")+" back",
	)

	view.Content = s.Render(form)
	return view
}
