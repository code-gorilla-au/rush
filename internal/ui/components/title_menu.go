package components

import (
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

type TitleItem int

const (
	TitleItemCreateCoach TitleItem = iota
	TitleItemLockerRoom
	TitleItemNewTournament
	TitleItemNewBattleSelection
)

func (i TitleItem) String() string {
	switch i {
	case TitleItemCreateCoach:
		return "Create Coach"
	case TitleItemLockerRoom:
		return "Locker Room"
	case TitleItemNewTournament:
		return "New Tournament"
	case TitleItemNewBattleSelection:
		return "New Battle"
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
		items = []TitleItem{TitleItemCreateCoach}
	} else {
		items = []TitleItem{TitleItemLockerRoom, TitleItemNewTournament, TitleItemNewBattleSelection}
	}
	return TitleMenu{
		items: items,
	}
}

func (m *TitleMenu) MoveUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *TitleMenu) MoveDown() {
	if m.cursor < len(m.items)-1 {
		m.cursor++
	}
}

func (m *TitleMenu) View(theme styles.IceTheme) string {
	var s string
	for i, item := range m.items {
		if i == m.cursor {
			s += theme.ListSelected.Render(">  " + item.String())
		} else {
			s += theme.Base.Render("   " + item.String())
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
		items = []TitleItem{TitleItemCreateCoach}
	} else {
		items = []TitleItem{TitleItemLockerRoom, TitleItemNewTournament, TitleItemNewBattleSelection}
	}
	m.items = items
	if m.cursor >= len(m.items) {
		m.cursor = 0
	}
}
