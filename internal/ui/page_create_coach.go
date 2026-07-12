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

	coachInput   textinput.Model
	teamInput    textinput.Model
	personaIndex int
	personas     []teams.CoachPersona
	focusIndex   int
	err          error
	keys         createCoachKeyMap
	footer       components.Footer
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
		personas: []teams.CoachPersona{
			teams.CoachPersonaVanguard,
			teams.CoachPersonaBastion,
			teams.CoachPersonaTrickster,
			teams.CoachPersonaWildcard,
		},
		theme:  theme,
		keys:   keys,
		footer: components.NewFooter(keys),
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

	case error:
		m.err = msg

	case uistate.MsgStateUpdated:
		return m.handleStateUpdated(msg)
	}

	return m.updateInputs(msg)
}

func (m *ModelCreateCoach) handleKeyMsg(msg tea.KeyMsg) (tea.Cmd, bool) {
	if key.Matches(msg, m.keys.Quit) {
		return tea.Quit, true
	}

	if m.handleFocusNavigation(msg) {
		return nil, false
	}

	if cmd, handled := m.handleEnter(msg); handled {
		return cmd, true
	}

	if key.Matches(msg, m.keys.Back) {
		return func() tea.Msg {
			return uistate.MsgSwitchPage{NewPage: uistate.PageTitle}
		}, true
	}

	if m.handlePersonaNavigation(msg) {
		return nil, false
	}

	return nil, false
}

func (m *ModelCreateCoach) handleFocusNavigation(msg tea.KeyMsg) bool {
	if key.Matches(msg, m.keys.Up) || key.Matches(msg, m.keys.ShiftTab) {
		m.focusIndex--
		if m.focusIndex < 0 {
			m.focusIndex = 2
		}
		m.updateFocus()
		return true
	}

	if key.Matches(msg, m.keys.Down) || key.Matches(msg, m.keys.Tab) {
		m.focusIndex++
		if m.focusIndex > 2 {
			m.focusIndex = 0
		}
		m.updateFocus()
		return true
	}

	return false
}

func (m *ModelCreateCoach) handleEnter(msg tea.KeyMsg) (tea.Cmd, bool) {
	if !key.Matches(msg, m.keys.Enter) {
		return nil, false
	}

	if m.focusIndex == 2 {
		return m.submit(), true
	}

	m.focusIndex++
	m.updateFocus()

	return nil, true
}

func (m *ModelCreateCoach) handlePersonaNavigation(msg tea.KeyMsg) bool {
	if m.focusIndex != 1 {
		return false
	}

	if key.Matches(msg, m.keys.Left) {
		m.personaIndex--
		if m.personaIndex < 0 {
			m.personaIndex = len(m.personas) - 1
		}
		return true
	}

	if key.Matches(msg, m.keys.Right) {
		m.personaIndex++
		if m.personaIndex >= len(m.personas) {
			m.personaIndex = 0
		}
		return true
	}

	return false
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
	switch m.focusIndex {
	case 0:
		m.coachInput.Focus()
		m.teamInput.Blur()
	case 1:
		m.coachInput.Blur()
		m.teamInput.Blur()
	case 2:
		m.coachInput.Blur()
		m.teamInput.Focus()
	}
}

func (m *ModelCreateCoach) submit() tea.Cmd {
	return func() tea.Msg {
		ctx := m.globalState.Context()
		coach, err := m.teamsSvc.CreateCoach(ctx, teams.CreateCoachParams{
			Name:      m.coachInput.Value(),
			Persona:   m.personas[m.personaIndex],
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

func (m *ModelCreateCoach) formatPersona() string {
	persona := m.personas[m.personaIndex]

	style := m.theme.Muted
	if m.focusIndex == 1 {
		style = m.theme.Highlight
	}

	return style.Render("< " + persona.Name() + " >")
}

func (m *ModelCreateCoach) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true

	errorView := ""
	if m.err != nil {
		errorView = m.theme.Muted.Render(m.err.Error())
	}

	form := lipgloss.JoinVertical(
		lipgloss.Left,
		m.theme.Logo.Render("RUSH - NEW CAREER"),
		"",
		m.theme.SecondaryHeader.Render("Coach Details"),
		m.coachInput.View(),
		"",
		m.theme.SecondaryHeader.Render("Persona"),
		m.formatPersona(),
		"",
		m.theme.SecondaryHeader.Render("Team Details"),
		m.teamInput.View(),
		"",
		errorView,
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
