package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

type ModelTitleSettings struct {
	width       int
	height      int
	theme       styles.IceTheme
	globalState *GlobalState
}

func NewModelTitleSettings(globalState *GlobalState, theme styles.IceTheme) *ModelTitleSettings {
	return &ModelTitleSettings{
		globalState: globalState,
		theme:       theme,
	}
}

func (m *ModelTitleSettings) Init() tea.Cmd {
	return nil
}

func (m *ModelTitleSettings) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch vMsg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = vMsg.Width
		m.height = vMsg.Height
	case MsgStateUpdated:
		m.globalState.Coach = vMsg.Coach
		m.globalState.Team = vMsg.Team
	}

	return m, nil
}

func (m *ModelTitleSettings) View() tea.View {
	return tea.NewView("Settings")
}
