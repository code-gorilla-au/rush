package ui

import tea "charm.land/bubbletea/v2"

type PageGameCompleteModel struct {
}

func NewPageGameComplete() *PageGameCompleteModel {
	return &PageGameCompleteModel{}
}

func (m *PageGameCompleteModel) Init() tea.Cmd {
	return nil
}

func (m *PageGameCompleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *PageGameCompleteModel) View() tea.View {
	return tea.NewView("Game Complete")
}
