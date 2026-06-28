package components

import (
	"strings"
	"testing"

	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/games"
)

func TestRound(t *testing.T) {
	t.Parallel()
	group := odize.NewGroup(t, nil)

	round := games.Round{
		TeamA: games.TeamFormation{
			Lanes: [3][]int{
				{1, 2},
				{3},
				{4, 5, 6},
			},
		},
		TeamB: games.TeamFormation{
			Lanes: [3][]int{
				{7},
				{8, 9},
				{},
			},
		},
	}

	err := group.
		Test("NewRound should initialize with given round and names", func(t *testing.T) {
			rComp := NewRound(round, "Team A", "Team B")
			odize.AssertEqual(t, "Team A", rComp.teamAName)
			odize.AssertEqual(t, "Team B", rComp.teamBName)
		}).
		Test("View should render team names and players in side-by-side formation", func(t *testing.T) {
			rComp := NewRound(round, "Team A", "Team B")
			rendered := rComp.View()

			odize.AssertTrue(t, strings.Contains(rendered, "Team A"))
			odize.AssertTrue(t, strings.Contains(rendered, "Team B"))
			odize.AssertTrue(t, strings.Contains(rendered, "|"))
			odize.AssertTrue(t, strings.Contains(rendered, "Lane 1"))
			odize.AssertTrue(t, strings.Contains(rendered, "Lane 2"))
			odize.AssertTrue(t, strings.Contains(rendered, "Lane 3"))
			// Check for player icons (dots)
			odize.AssertTrue(t, strings.Contains(rendered, "●"))
		}).
		Run()

	odize.AssertNoError(t, err)
}
