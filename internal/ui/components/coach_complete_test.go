package components

import (
	"strings"
	"testing"

	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
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

	theme := styles.NewIceTheme()

	err := group.
		Test("CoachWinnerHuman renders correctly", func(t *testing.T) {
			c := NewCoachWinnerHuman(team, coachHuman)
			view := c.View(theme)

			odize.AssertTrue(t, contains(view, "Winner: Test Team"))
			odize.AssertTrue(t, contains(view, "Human Coach (Human Coach)"))
			odize.AssertTrue(t, contains(view, "Winning Roster:"))
			odize.AssertTrue(t, contains(view, "• Player 1"))
			odize.AssertTrue(t, contains(view, "• Player 2"))
		}).
		Test("CoachWinnerAI renders correctly", func(t *testing.T) {
			c := NewCoachWinnerAI(team, coachAI)
			view := c.View(theme)

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
