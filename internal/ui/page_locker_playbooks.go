package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/playbooks"
)

type ModelLockerPlaybooks struct {
	width       int
	height      int
	theme       IceTheme
	globalState *GlobalState
	playbookSvc *playbooks.Service
}

func NewModelLockerPlaybooks(state *GlobalState, playbookSvc *playbooks.Service) *ModelLockerPlaybooks {
	return &ModelLockerPlaybooks{
		theme:       NewIceTheme(),
		globalState: state,
		playbookSvc: playbookSvc,
	}
}

func (m *ModelLockerPlaybooks) Init() tea.Cmd {
	return nil
}

func (m *ModelLockerPlaybooks) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *ModelLockerPlaybooks) View() tea.View {
	return tea.NewView("playbooks")
}
