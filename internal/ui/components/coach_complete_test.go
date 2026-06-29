package components

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/teams"
)

func TestCoachWinnerComponents(t *testing.T) {
	group := odize.NewGroup(t, nil)

	team := &teams.Team{
		Name: "Test Team",
		Players: []teams.Player{
			{Name: "Player 1"},
			{Name: "Player 2"},
		},
	}

	coachHuman := &teams.Coach{
		Name:    "Human Coach",
		IsHuman: true,
	}

	coachAI := &teams.Coach{
		Name:    "AI Coach",
		IsHuman: false,
	}

	coachStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))

	err := group.
		Test("CoachWinnerHuman renders correctly", func(t *testing.T) {
			c := NewCoachWinnerHuman(team, coachHuman)
			view := c.View(coachStyle)

			odize.AssertTrue(t, contains(view, "Winner: Test Team"))
			odize.AssertTrue(t, contains(view, "Human Coach (Human Coach)"))
			odize.AssertTrue(t, contains(view, "Winning Roster:"))
			odize.AssertTrue(t, contains(view, "• Player 1"))
			odize.AssertTrue(t, contains(view, "• Player 2"))
		}).
		Test("CoachWinnerAI renders correctly", func(t *testing.T) {
			c := NewCoachWinnerAI(team, coachAI)
			view := c.View(coachStyle)

			odize.AssertTrue(t, contains(view, "Winner: Test Team"))
			odize.AssertTrue(t, contains(view, "AI Coach (AI Coach)"))
			odize.AssertTrue(t, contains(view, "Winning Roster:"))
			odize.AssertTrue(t, contains(view, "• Player 1"))
			odize.AssertTrue(t, contains(view, "• Player 2"))
		}).
		Run()

	odize.AssertNoError(t, err)
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
