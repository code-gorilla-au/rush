package ui

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/components"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"
)

type createCoachKeyMap struct {
	uistate.KeyMap
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
		KeyMap: uistate.NewKeyMap(),
	}
}

type ModelCreateCoach struct {
	width       int
	height      int
	theme       styles.IceTheme
	globalState *uistate.GlobalState
	teamsSvc    *teams.Service

	coachInput textinput.Model
	teamInput  textinput.Model
	focusIndex int
	err        error
	keys       createCoachKeyMap
	footer     components.Footer
}

func NewModelCreateCoach(state *uistate.GlobalState, teamsSvc *teams.Service, theme styles.IceTheme) *ModelCreateCoach {
	c := textinput.New()
	c.Placeholder = "Coach Name"
	c.Focus()
	c.CharLimit = 156
	c.SetWidth(20)

	t := textinput.New()
	t.Placeholder = "Team Name"
	t.CharLimit = 156
	t.SetWidth(20)

	keys := newCreateCoachKeyMap()

	return &ModelCreateCoach{
		globalState: state,
		teamsSvc:    teamsSvc,
		coachInput:  c,
		teamInput:   t,
		theme:       theme,
		keys:        keys,
		footer:      components.NewFooter(keys),
	}
}

func (m *ModelCreateCoach) Init() tea.Cmd {
	return func() tea.Msg { return textinput.Blink() }
}

func (m *ModelCreateCoach) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case uistate.MsgSwitchPage:
		if msg.NewPage == uistate.PageCreateCoach {
			return m, m.Init()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.footer.Update(msg)

	case tea.KeyMsg:
		if cmd, done := m.handleKeyMsg(msg); done {
			return m, cmd
		}

	case uistate.MsgStateUpdated:
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
			return uistate.MsgSwitchPage{NewPage: uistate.PageTitle}
		}, true
	}

	return nil, false
}

func (m *ModelCreateCoach) handleStateUpdated(msg uistate.MsgStateUpdated) (tea.Model, tea.Cmd) {
	m.globalState.Coach = msg.Coach
	m.globalState.Team = msg.Team

	return m, func() tea.Msg {
		return uistate.MsgSwitchPage{NewPage: uistate.PageLockerRoom}
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
		coach, err := m.teamsSvc.CreateCoach(ctx, teams.CreateCoachParams{
			Name:      m.coachInput.Value(),
			IsHuman:   true,
			IsDefault: true,
		})
		if err != nil {
			return err
		}

		team, err := m.teamsSvc.CreateTeam(ctx, m.teamInput.Value(), coach.ID, true)
		if err != nil {
			return err
		}

		return uistate.MsgStateUpdated{
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
		m.theme.SecondaryHeader.Render("Coach Details"),
		m.coachInput.View(),
		"",
		m.theme.SecondaryHeader.Render("Team Details"),
		m.teamInput.View(),
		"",
		m.footer.View(m.theme),
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
