package ui

import (
	tea "charm.land/bubbletea/v2"
)

type ModelCreateCoach struct {
	width  int
	height int
	theme  IceTheme
}

func NewModelCreateCoach() *ModelCreateCoach {
	return &ModelCreateCoach{}
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
