package ui

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/components"
)

type createCoachKeyMap struct {
	components.CommonKeys
	Back     key.Binding
	Enter    key.Binding
	Up       key.Binding
	Down     key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
}

func (k createCoachKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Back, k.Quit}
}

func (k createCoachKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.Back},
		{k.Up, k.Down, k.Tab, k.ShiftTab},
		{k.Quit},
	}
}

func newCreateCoachKeyMap() createCoachKeyMap {
	return createCoachKeyMap{
		CommonKeys: components.NewCommonKeys(),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "continue"),
		),
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next"),
		),
		ShiftTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev"),
		),
	}
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
	keys       createCoachKeyMap
	footer     components.Footer
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

	keys := newCreateCoachKeyMap()

	return &ModelCreateCoach{
		globalState: state,
		teamsSvc:    teamsSvc,
		coachInput:  c,
		teamInput:   t,
		theme:       NewIceTheme(),
		keys:        keys,
		footer:      components.NewFooter(keys),
	}
}

func (m *ModelCreateCoach) Init() tea.Cmd {
	return func() tea.Msg { return textinput.Blink() }
}

func (m *ModelCreateCoach) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.footer.Update(msg)

	case tea.KeyMsg:
		if cmd, done := m.handleKeyMsg(msg); done {
			return m, cmd
		}

	case MsgStateUpdated:
		return m.handleStateUpdated(msg)
	}

	return m.updateInputs(msg)
}

func (m *ModelCreateCoach) handleKeyMsg(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		return tea.Quit, true

	case key.Matches(msg, m.keys.Up), key.Matches(msg, m.keys.ShiftTab):
		m.focusIndex--
		if m.focusIndex < 0 {
			m.focusIndex = 1
		}
		m.updateFocus()
		return nil, false

	case key.Matches(msg, m.keys.Down), key.Matches(msg, m.keys.Tab):
		m.focusIndex++
		if m.focusIndex > 1 {
			m.focusIndex = 0
		}
		m.updateFocus()
		return nil, false

	case key.Matches(msg, m.keys.Enter):
		if m.focusIndex == 1 {
			return m.submit(), true
		}
		m.focusIndex++
		m.updateFocus()
		return nil, false

	case key.Matches(msg, m.keys.Back):
		return func() tea.Msg {
			return MsgSwitchPage{NewPage: PageTitle}
		}, true
	}

	return nil, false
}

func (m *ModelCreateCoach) handleStateUpdated(msg MsgStateUpdated) (tea.Model, tea.Cmd) {
	m.globalState.Coach = msg.Coach
	m.globalState.Team = msg.Team

	return m, func() tea.Msg {
		return MsgSwitchPage{NewPage: PageLockerRoom}
	}
}

func (m *ModelCreateCoach) updateInputs(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
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

		return MsgStateUpdated{
			Coach: &coach,
			Team:  &team,
		}
	}
}

func (m *ModelCreateCoach) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true

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
		m.footer.View(m.theme.Footer),
	)

	centeredContent := lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		form,
	)

	view.Content = m.theme.Base.
		Width(m.width).
		Height(m.height).
		Render(centeredContent)

	return view
}
