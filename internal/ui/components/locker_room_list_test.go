package components

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/odize"
)

func TestLockerRoomList(t *testing.T) {
	group := odize.NewGroup(t, nil)

	group.Test("NewLockerRoomList should have 3 items", func(t *testing.T) {
		l := NewLockerRoomList()
		odize.AssertEqual(t, 3, len(l.items))
		odize.AssertEqual(t, ItemPlayers, l.items[0])
		odize.AssertEqual(t, ItemPlaybooks, l.items[1])
		odize.AssertEqual(t, TitleItemSettings, l.items[2])
	})

	group.Test("Update should move cursor down", func(t *testing.T) {
		l := NewLockerRoomList()
		odize.AssertEqual(t, 0, l.cursor)

		l.Update(tea.KeyPressMsg{Text: "down"})
		odize.AssertEqual(t, 1, l.cursor)
		odize.AssertEqual(t, ItemPlaybooks, l.SelectedItem())
	})

	group.Test("Update should move cursor up", func(t *testing.T) {
		l := NewLockerRoomList()
		l.cursor = 1

		l.Update(tea.KeyPressMsg{Text: "up"})
		odize.AssertEqual(t, 0, l.cursor)
		odize.AssertEqual(t, ItemPlayers, l.SelectedItem())
	})

	group.Test("Update should not move cursor out of bounds", func(t *testing.T) {
		l := NewLockerRoomList()

		l.Update(tea.KeyPressMsg{Text: "up"})
		odize.AssertEqual(t, 0, l.cursor)

		l.cursor = 2
		l.Update(tea.KeyPressMsg{Text: "down"})
		odize.AssertEqual(t, 2, l.cursor)
	})

	err := group.Run()
	odize.AssertNoError(t, err)
}
