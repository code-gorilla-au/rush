package games

import (
	"encoding/json"
	"testing"

	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/playbooks"
)

func TestGame_ToGameModel(t *testing.T) {
	group := odize.NewGroup(t, nil)

	group.Test("should marshal rounds correctly", func(t *testing.T) {
		teamA := TeamConfig{TeamID: 1, TeamName: "A", Formations: make([]playbooks.Formation, 10)}
		teamB := TeamConfig{TeamID: 2, TeamName: "B", Formations: make([]playbooks.Formation, 10)}
		game := Game{
			rounds: generateRounds(teamA, teamB),
		}

		model, err := toGameModel(game)
		odize.AssertNoError(t, err)

		var rounds [10]Round
		err = json.Unmarshal(model.Rounds, &rounds)
		odize.AssertNoError(t, err)

		odize.AssertEqual(t, teamA.TeamID, rounds[0].TeamA.TeamID)
	})

	err := group.Run()
	odize.AssertNoError(t, err)
}

func TestGame_ResolveRound(t *testing.T) {
	group := odize.NewGroup(t, nil)

	group.Test("should resolve the first round (index 0)", func(t *testing.T) {
		teamA := TeamConfig{TeamID: 1, TeamName: "A", Formations: make([]playbooks.Formation, 10)}
		teamB := TeamConfig{TeamID: 2, TeamName: "B", Formations: make([]playbooks.Formation, 10)}

		// Fill some players in round 0
		teamA.Formations[0] = playbooks.Formation{Lane1: 1, Lane2: 1, Lane3: 1}
		teamB.Formations[0] = playbooks.Formation{Lane1: 1, Lane2: 1, Lane3: 1}

		game := Game{
			currentRound: 0,
			rounds:       generateRounds(teamA, teamB),
		}

		rolls := []int{6, 1}
		idx := 0
		rollFn := func() int {
			val := rolls[idx%2]
			idx++
			return val
		}
		res, err := game.ResolveRound(rollFn)

		odize.AssertNoError(t, err)
		odize.AssertTrue(t, res.Outcome == ResultTeamA || res.Outcome == ResultTeamB)
		odize.AssertEqual(t, int64(1), game.currentRound)

		// Check if players were removed from round 0
		round0 := game.Rounds()[0]
		// Since it's 1v1 in each lane and Team A rolled 6 vs Team B rolled 1,
		// Team B should have 0 players in Lane 1, 2, 3.
		for i := 0; i < 3; i++ {
			odize.AssertEqual(t, 0, len(round0.TeamB.Lanes[i]))
		}
	})

	err := group.Run()
	odize.AssertNoError(t, err)
}
