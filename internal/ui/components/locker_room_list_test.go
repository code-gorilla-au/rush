package components

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

func TestLockerRoomList(t *testing.T) {
	group := odize.NewGroup(t, nil)
	theme := styles.NewIceTheme()

	group.Test("NewLockerRoomList should have 3 items", func(t *testing.T) {
		l := NewLockerRoomList(theme)
		odize.AssertEqual(t, 3, len(l.Model.Items()))
		odize.AssertEqual(t, ItemPlayers, l.SelectedItem())
	})

	group.Test("Update should move cursor down", func(t *testing.T) {
		l := NewLockerRoomList(theme)
		odize.AssertEqual(t, 0, l.Model.Index())

		l.Update(tea.KeyPressMsg{Text: "down"})
		odize.AssertEqual(t, 1, l.Model.Index())
		odize.AssertEqual(t, ItemPlaybooks, l.SelectedItem())
	})

	group.Test("Update should move cursor up", func(t *testing.T) {
		l := NewLockerRoomList(theme)
		l.Model.Select(1)

		l.Update(tea.KeyPressMsg{Text: "up"})
		odize.AssertEqual(t, 0, l.Model.Index())
		odize.AssertEqual(t, ItemPlayers, l.SelectedItem())
	})

	group.Test("Update should not move cursor out of bounds", func(t *testing.T) {
		l := NewLockerRoomList(theme)

		l.Update(tea.KeyPressMsg{Text: "up"})
		odize.AssertEqual(t, 0, l.Model.Index())

		l.Model.Select(2)
		l.Update(tea.KeyPressMsg{Text: "down"})
		odize.AssertEqual(t, 2, l.Model.Index())
	})

	err := group.Run()
	odize.AssertNoError(t, err)
}
