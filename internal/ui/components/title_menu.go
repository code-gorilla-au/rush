package components

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type TitleItem int

const (
	ItemCreateCoach TitleItem = iota
	ItemLockerRoom
	ItemNewTournament
	ItemNewBattle
	ItemNewSettings
)

func (i TitleItem) String() string {
	switch i {
	case ItemCreateCoach:
		return "Create Coach"
	case ItemLockerRoom:
		return "Locker Room"
	case ItemNewTournament:
		return "New Tournament"
	case ItemNewBattle:
		return "New Battle"
	case ItemNewSettings:
		return "Settings"
	}
	return ""
}

type TitleMenu struct {
	cursor int
	items  []TitleItem
}

func NewTitleMenu(hasCoach bool) TitleMenu {
	var items []TitleItem
	if !hasCoach {
		items = []TitleItem{ItemCreateCoach}
	} else {
		items = []TitleItem{ItemLockerRoom, ItemNewTournament, ItemNewBattle, ItemNewSettings}
	}
	return TitleMenu{
		items: items,
	}
}

func (m *TitleMenu) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		}
	}
}

func (m *TitleMenu) View(itemStyle lipgloss.Style, selectedStyle lipgloss.Style) string {
	var s string
	for i, item := range m.items {
		if i == m.cursor {
			s += selectedStyle.Render("> " + item.String())
		} else {
			s += itemStyle.Render("  " + item.String())
		}
		if i < len(m.items)-1 {
			s += "\n"
		}
	}
	return s
}

func (m *TitleMenu) SelectedItem() TitleItem {
	if m.cursor < 0 || m.cursor >= len(m.items) {
		return -1
	}
	return m.items[m.cursor]
}

func (m *TitleMenu) SetHasCoach(hasCoach bool) {
	var items []TitleItem
	if !hasCoach {
		items = []TitleItem{ItemCreateCoach}
	} else {
		items = []TitleItem{ItemLockerRoom, ItemNewTournament, ItemNewBattle, ItemNewSettings}
	}
	m.items = items
	if m.cursor >= len(m.items) {
		m.cursor = 0
	}
}
