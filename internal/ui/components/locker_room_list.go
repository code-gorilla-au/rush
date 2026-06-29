package components

import (
	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

type LockerRoomItem int

const (
	ItemPlayers LockerRoomItem = iota
	ItemPlaybooks
	ItemSettings
)

func (i LockerRoomItem) String() string {
	switch i {
	case ItemPlayers:
		return "Players"
	case ItemPlaybooks:
		return "Playbooks"
	case ItemSettings:
		return "Settings"
	}
	return ""
}

type LockerRoomList struct {
	cursor int
	items  []LockerRoomItem
}

func NewLockerRoomList() LockerRoomList {
	return LockerRoomList{
		items: []LockerRoomItem{ItemPlayers, ItemPlaybooks, ItemSettings},
	}
}

func (l *LockerRoomList) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if l.cursor > 0 {
				l.cursor--
			}
		case "down", "j":
			if l.cursor < len(l.items)-1 {
				l.cursor++
			}
		}
	}
}

func (l *LockerRoomList) View(theme styles.IceTheme) string {
	var s string
	for i, item := range l.items {
		if i == l.cursor {
			s += theme.ListSelected.Render(">  " + item.String())
		} else {
			s += theme.Base.Render("   " + item.String())
		}
		if i < len(l.items)-1 {
			s += "\n"
		}
	}
	return s
}

func (l *LockerRoomList) SelectedItem() LockerRoomItem {
	return l.items[l.cursor]
}
