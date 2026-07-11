package components

import (
	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

type LockerRoomItem int

const (
	ItemPlayers LockerRoomItem = iota
	ItemTeamStatistics
	ItemPlaybooks
	ItemCoachEdit
)

func (i LockerRoomItem) String() string {
	switch i {
	case ItemPlayers:
		return "Players"
	case ItemTeamStatistics:
		return "Team Statistics"
	case ItemPlaybooks:
		return "Playbooks"
	case ItemCoachEdit:
		return "Coach Edit"
	}
	return ""
}

type LockerRoomList struct {
	List[LockerRoomItem]
}

func NewLockerRoomList(theme styles.IceTheme) LockerRoomList {
	return LockerRoomList{
		List: NewList(ListConfig[LockerRoomItem]{
			Items: []LockerRoomItem{ItemPlayers, ItemTeamStatistics, ItemPlaybooks},
			ItemMapper: func(i LockerRoomItem) ListItem[LockerRoomItem] {
				return ListItem[LockerRoomItem]{
					Data:     i,
					TitleVal: i.String(),
				}
			},
			EnableFiltering:   false,
			DisableAutoResize: true,
		}, theme),
	}
}

func (l *LockerRoomList) Update(msg tea.Msg) (LockerRoomList, tea.Cmd) {
	var cmd tea.Cmd
	l.List, cmd = l.List.Update(msg)
	return *l, cmd
}

func (l *LockerRoomList) View(theme styles.IceTheme) string {
	return l.List.View()
}

func (l *LockerRoomList) SelectedItem() LockerRoomItem {
	if item, ok := l.List.SelectedItem(); ok {
		return item
	}
	return -1
}
