package ui

import (
	tea "charm.land/bubbletea/v2"
)

type ModelCreateCoach struct {
	width       int
	height      int
	theme       IceTheme
	globalState *GlobalState
}

func NewModelCreateCoach(state *GlobalState) *ModelCreateCoach {
	return &ModelCreateCoach{
		globalState: state,
	}
}

func (m ModelCreateCoach) Init() tea.Cmd {
	return nil
}

func (m ModelCreateCoach) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return m, nil
}

func (m ModelCreateCoach) View() tea.View {

	view := tea.NewView("create coach")
	view.AltScreen = true
	return view
}
