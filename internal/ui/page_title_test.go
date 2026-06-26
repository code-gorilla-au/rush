package ui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/teams"
)

func TestModelTitle(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should route to create coach when coach is nil and enter is pressed", func(t *testing.T) {
			m := NewModelTitle(&GlobalState{Coach: nil})
			m.width = 100
			m.height = 50

			_, cmd := m.Update(tea.KeyPressMsg{Text: "enter"})
			odize.AssertTrue(t, cmd != nil)

			msg := cmd()
			switch switchMsg := msg.(type) {
			case MsgSwitchPage:
				odize.AssertEqual(t, PageCreateCoach, switchMsg.NewPage)
			default:
				t.Fatalf("expected MsgSwitchPage, got %T", msg)
			}
		}).
		Test("should route to locker room when coach is not nil and enter is pressed", func(t *testing.T) {
			m := NewModelTitle(&GlobalState{Coach: &teams.Coach{Name: "Coach Carter"}})
			m.width = 100
			m.height = 50

			_, cmd := m.Update(tea.KeyPressMsg{Text: "enter"})
			odize.AssertTrue(t, cmd != nil)

			msg := cmd()
			switch switchMsg := msg.(type) {
			case MsgSwitchPage:
				odize.AssertEqual(t, PageLockerRoom, switchMsg.NewPage)
			default:
				t.Fatalf("expected MsgSwitchPage, got %T", msg)
			}
		}).
		Run()

	odize.AssertNoError(t, err)
}
