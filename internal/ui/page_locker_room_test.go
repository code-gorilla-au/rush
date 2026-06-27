package ui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/ui/components"
)

func TestModelLockerRoom_Selection(t *testing.T) {
	group := odize.NewGroup(t, nil)

	group.Test("should route to locker players when players item is selected", func(t *testing.T) {
		state := &GlobalState{}
		m := NewModelLockerRoom(state)

		// Ensure ItemPlayers is selected (it is by default)
		odize.AssertEqual(t, components.ItemPlayers, m.list.SelectedItem())

		// Simulate Enter key press
		_, cmd := m.Update(tea.KeyPressMsg{Text: "enter"})

		odize.AssertTrue(t, cmd != nil)
		msg := cmd()
		switch v := msg.(type) {
		case MsgSwitchPage:
			odize.AssertEqual(t, PageLockerPlayers, v.NewPage)
		default:
			t.Fatalf("expected MsgSwitchPage, got %T", msg)
		}
	})

	group.Test("should route to locker playbooks when playbooks item is selected", func(t *testing.T) {
		state := &GlobalState{}
		m := NewModelLockerRoom(state)

		// Select Playbooks (it's the second item)
		m.Update(tea.KeyPressMsg{Text: "down"})
		odize.AssertEqual(t, components.ItemPlaybooks, m.list.SelectedItem())

		// Simulate Enter key press
		_, cmd := m.Update(tea.KeyPressMsg{Text: "enter"})

		odize.AssertTrue(t, cmd != nil)
		msg := cmd()
		switch v := msg.(type) {
		case MsgSwitchPage:
			odize.AssertEqual(t, PageLockerPlaybooksList, v.NewPage)
		default:
			t.Fatalf("expected MsgSwitchPage, got %T", msg)
		}
	})

	err := group.Run()
	odize.AssertNoError(t, err)
}
