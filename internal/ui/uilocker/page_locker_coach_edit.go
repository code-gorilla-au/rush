package uilocker

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

type coachEditKeyMap struct {
	uistate.KeyMap
}

func (k coachEditKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Back, k.Quit}
}

func (k coachEditKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.Back},
		{k.Up, k.Down, k.Tab, k.ShiftTab},
		{k.Quit},
	}
}

func newCoachEditKeyMap() coachEditKeyMap {
	return coachEditKeyMap{
		KeyMap: uistate.NewKeyMap(),
	}
}

type PageLockerCoachEdit struct {
	width       int
	height      int
	theme       styles.IceTheme
	globalState *uistate.GlobalState
	teamsSvc    *teams.Service

	coachInput   textinput.Model
	personaIndex int
	personas     []teams.CoachPersona
	focusIndex   int
	err          error
	keys         coachEditKeyMap
	footer       components.Footer
}

func NewPageLockerCoachEdit(state *uistate.GlobalState, teamsSvc *teams.Service, theme styles.IceTheme) *PageLockerCoachEdit {
	c := textinput.New()
	c.Placeholder = "Coach Name"
	c.Focus()
	c.CharLimit = 156
	c.SetWidth(20)
	if state.Coach != nil {
		c.SetValue(state.Coach.Name)
	}

	keys := newCoachEditKeyMap()

	personaIndex := 0
	personas := []teams.CoachPersona{
		teams.CoachPersonaVanguard,
		teams.CoachPersonaBastion,
		teams.CoachPersonaTrickster,
		teams.CoachPersonaWildcard,
	}
	if state.Coach != nil {
		for i, p := range personas {
			if p == state.Coach.Persona {
				personaIndex = i
				break
			}
		}
	}

	return &PageLockerCoachEdit{
		globalState:  state,
		teamsSvc:     teamsSvc,
		coachInput:   c,
		personas:     personas,
		personaIndex: personaIndex,
		theme:        theme,
		keys:         keys,
		footer:       components.NewFooter(keys),
	}
}

func (m *PageLockerCoachEdit) Init() tea.Cmd {
	return textinput.Blink
}

func (m *PageLockerCoachEdit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.footer.Update(msg)

	case tea.KeyMsg:
		if cmd, done := m.handleKeyMsg(msg); done {
			return m, cmd
		}

	case MsgSwitchLockerPage:
		if msg.NewPage == SubPageLockerCoachEdit {
			m.refresh()
		}

	case error:
		m.err = msg
	}

	return m.updateInputs(msg)
}

func (m *PageLockerCoachEdit) handleKeyMsg(msg tea.KeyMsg) (tea.Cmd, bool) {
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
			return MsgSwitchLockerPage{NewPage: SubPageLockerRoom}
		}, true

	case key.Matches(msg, m.keys.Left):
		if m.focusIndex == 1 {
			m.personaIndex--
			if m.personaIndex < 0 {
				m.personaIndex = len(m.personas) - 1
			}
			return nil, false
		}

	case key.Matches(msg, m.keys.Right):
		if m.focusIndex == 1 {
			m.personaIndex++
			if m.personaIndex >= len(m.personas) {
				m.personaIndex = 0
			}
			return nil, false
		}
	}

	return nil, false
}

func (m *PageLockerCoachEdit) updateInputs(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.coachInput, cmd = m.coachInput.Update(msg)
	return m, cmd
}

func (m *PageLockerCoachEdit) updateFocus() {
	switch m.focusIndex {
	case 0:
		m.coachInput.Focus()
	case 1:
		m.coachInput.Blur()
	}
}

func (m *PageLockerCoachEdit) submit() tea.Cmd {
	return func() tea.Msg {
		if m.globalState.Coach == nil {
			return nil
		}

		ctx := m.globalState.Context()
		err := m.teamsSvc.UpdateCoach(ctx, m.globalState.Coach.ID, m.coachInput.Value(), m.personas[m.personaIndex])
		if err != nil {
			return err
		}

		m.globalState.Coach.Name = m.coachInput.Value()
		m.globalState.Coach.Persona = m.personas[m.personaIndex]

		return MsgSwitchLockerPage{NewPage: SubPageLockerRoom}
	}
}

func (m *PageLockerCoachEdit) formatPersona() string {
	persona := m.personas[m.personaIndex]

	style := m.theme.Muted
	if m.focusIndex == 1 {
		style = m.theme.Highlight
	}

	return style.Render("< " + persona.Name() + " >")
}

func (m *PageLockerCoachEdit) refresh() {
	if m.globalState.Coach == nil {
		return
	}

	if m.coachInput.Value() != m.globalState.Coach.Name {
		m.coachInput.SetValue(m.globalState.Coach.Name)
	}

	for i, p := range m.personas {
		if p == m.globalState.Coach.Persona {
			if m.personaIndex != i {
				m.personaIndex = i
			}
			break
		}
	}
}

func (m *PageLockerCoachEdit) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true

	errorView := ""
	if m.err != nil {
		errorView = m.theme.Muted.Render(m.err.Error())
	}

	form := lipgloss.JoinVertical(
		lipgloss.Left,
		m.theme.Logo.Render("RUSH - EDIT COACH"),
		"",
		m.theme.SecondaryHeader.Render("Coach Details"),
		m.coachInput.View(),
		"",
		m.theme.SecondaryHeader.Render("Persona"),
		m.formatPersona(),
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
