package ui

import tea "charm.land/bubbletea/v2"

type ModelLockerRoom struct {
	width       int
	height      int
	theme       IceTheme
	globalState GlobalState
}

func NewModelLockerRoom(globalState GlobalState) *ModelLockerRoom {
	return &ModelLockerRoom{
		globalState: globalState,
	}
}

func (m ModelLockerRoom) Init() tea.Cmd {
	return nil
}

func (m ModelLockerRoom) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ModelLockerRoom) View() tea.View {
	return tea.NewView("locker room")
}
