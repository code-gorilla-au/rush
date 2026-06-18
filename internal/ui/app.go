package ui

import (
	tea "charm.land/bubbletea/v2"
)

type MsgSwitchPage struct {
	NewPage Page
}

type Page int

const (
	PageTitle Page = iota + 1
	PageCreateCoach
	PageLockerRoom
)

type RootModel struct {
	width           int
	height          int
	theme           IceTheme
	currentPage     Page
	pageTitle       tea.Model
	pageCreateCoach tea.Model
	pageLockerRoom  tea.Model
}

// New returns a new UI model.
func New() RootModel {
	return RootModel{
		theme:           NewIceTheme(),
		currentPage:     PageTitle,
		pageTitle:       NewModelTitle(),
		pageCreateCoach: NewModelCreateCoach(),
		pageLockerRoom:  NewModelLockerRoom(),
	}
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case MsgSwitchPage:
		m.currentPage = msg.NewPage
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.pageTitle, _ = m.pageTitle.Update(msg)
		m.pageCreateCoach, _ = m.pageCreateCoach.Update(msg)
		m.pageLockerRoom, _ = m.pageLockerRoom.Update(msg)
	}

	switch m.currentPage {
	case PageTitle:
		m.pageTitle, _ = m.pageTitle.Update(msg)
	case PageCreateCoach:
		m.pageCreateCoach, _ = m.pageCreateCoach.Update(msg)
	case PageLockerRoom:
		m.pageLockerRoom, _ = m.pageLockerRoom.Update(msg)
	}

	return m, nil
}

func (m RootModel) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Initializing...")
	}

	switch m.currentPage {
	case PageTitle:
		return m.pageTitle.View()
	case PageCreateCoach:
		return m.pageCreateCoach.View()
	case PageLockerRoom:
		return m.pageLockerRoom.View()
	}

	return tea.NewView("unknown page")
}
