package ui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/teams"
)

func TestModelNewBattleSelection_Rendering(t *testing.T) {
	group := odize.NewGroup(t, nil)

	group.Test("should load AI coaches and render them", func(t *testing.T) {
		state := &GlobalState{}

		m := NewModelNewBattleSelection(state, nil)
		m.width = 100
		m.height = 40

		aiTeams := []AITeamItem{
			{
				coach: teams.Coach{Name: "Coach A"},
				team:  teams.Team{Name: "Team A"},
			},
		}

		msg := msgAICoachesLoaded{aiTeams: aiTeams}
		m.Update(msg)

		odize.AssertEqual(t, 1, len(m.aiCoaches))

		view := m.View()
		content := view.Content

		odize.AssertTrue(t, strings.Contains(content, "Select your opponent"))
		odize.AssertTrue(t, !strings.Contains(content, "No AI coaches available"))
	})

	group.Test("should handle select coach and return to title", func(t *testing.T) {
		state := &GlobalState{}
		m := NewModelNewBattleSelection(state, nil)
		m.width = 100
		m.height = 40
		m.aiCoaches = []AITeamItem{
			{
				coach: teams.Coach{ID: 1, Name: "Coach A"},
				team:  teams.Team{ID: 1, Name: "Team A"},
			},
		}
		m.selectedCoachIdx = 0

		// Simulate Enter key
		_, cmd := m.Update(tea.KeyPressMsg{Text: "enter"})

		odize.AssertTrue(t, cmd != nil)
		msg := cmd()
		switch v := msg.(type) {
		case MsgSwitchPage:
			odize.AssertEqual(t, PageTitle, v.NewPage)
		default:
			t.Fatalf("expected MsgSwitchPage, got %T", msg)
		}
	})

	group.Test("should handle back navigation from coach selection", func(t *testing.T) {
		state := &GlobalState{}
		m := NewModelNewBattleSelection(state, nil)

		_, cmd := m.Update(tea.KeyPressMsg{Text: "esc"})

		odize.AssertTrue(t, cmd != nil)
		msg := cmd()
		switch v := msg.(type) {
		case MsgSwitchPage:
			odize.AssertEqual(t, PageTitle, v.NewPage)
		default:
			t.Fatalf("expected MsgSwitchPage, got %T", msg)
		}
	})

	group.Test("should handle navigation between coaches", func(t *testing.T) {
		state := &GlobalState{}
		m := NewModelNewBattleSelection(state, nil)
		m.aiCoaches = []AITeamItem{
			{coach: teams.Coach{Name: "Coach A"}},
			{coach: teams.Coach{Name: "Coach B"}},
			{coach: teams.Coach{Name: "Coach C"}},
		}
		m.selectedCoachIdx = 0

		// Move down
		m.Update(tea.KeyPressMsg{Text: "down"})
		odize.AssertEqual(t, 1, m.selectedCoachIdx)

		// Move down again
		m.Update(tea.KeyPressMsg{Text: "j"})
		odize.AssertEqual(t, 2, m.selectedCoachIdx)

		// Move down at the end (should stay)
		m.Update(tea.KeyPressMsg{Text: "down"})
		odize.AssertEqual(t, 2, m.selectedCoachIdx)

		// Move up
		m.Update(tea.KeyPressMsg{Text: "up"})
		odize.AssertEqual(t, 1, m.selectedCoachIdx)

		// Move up again
		m.Update(tea.KeyPressMsg{Text: "k"})
		odize.AssertEqual(t, 0, m.selectedCoachIdx)

		// Move up at the beginning (should stay)
		m.Update(tea.KeyPressMsg{Text: "up"})
		odize.AssertEqual(t, 0, m.selectedCoachIdx)
	})

	err := group.Run()
	odize.AssertNoError(t, err)
}
