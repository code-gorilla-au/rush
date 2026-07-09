package games

import (
	"testing"

	"github.com/code-gorilla-au/odize"
)

func fillLaneWithIDs(t *testing.T, team *TeamFormation, lane int, playerIDs []int64) {
	remainder := team.LaneFill(lane, len(playerIDs), playerIDs)
	odize.AssertEqual(t, 0, len(remainder))
	odize.AssertEqual(t, len(playerIDs), len(team.Lanes[lane]))
	if len(playerIDs) > 0 {
		odize.AssertEqual(t, playerIDs, team.Lanes[lane])
	}
}

func TestFillSquad(t *testing.T) {
	group := odize.NewGroup(t, nil)

	group.Test("FillSquad should populate lanes with player IDs in order", func(t *testing.T) {
		r := NewRound()

		r.FillSquad(
			LanesConfig{
				TeamID:  10,
				Players: []int64{101, 102, 103, 104, 105, 106},
				Lane1:   2,
				Lane2:   1,
				Lane3:   3,
			},
			LanesConfig{
				TeamID:  20,
				Players: []int64{201, 202, 203, 204, 205, 206},
				Lane1:   1,
				Lane2:   3,
				Lane3:   2,
			},
		)

		odize.AssertEqual(t, int64(10), r.TeamA.TeamID)
		odize.AssertEqual(t, []int64{101, 102}, r.TeamA.Lanes[0])
		odize.AssertEqual(t, []int64{103}, r.TeamA.Lanes[1])
		odize.AssertEqual(t, []int64{104, 105, 106}, r.TeamA.Lanes[2])

		odize.AssertEqual(t, int64(20), r.TeamB.TeamID)
		odize.AssertEqual(t, []int64{201}, r.TeamB.Lanes[0])
		odize.AssertEqual(t, []int64{202, 203, 204}, r.TeamB.Lanes[1])
		odize.AssertEqual(t, []int64{205, 206}, r.TeamB.Lanes[2])
	})

	err := group.Run()
	odize.AssertNoError(t, err)
}

func TestResolveLane(t *testing.T) {
	group := odize.NewGroup(t, nil)

	group.Test("Team A should win when Team B runs out of players", func(t *testing.T) {
		r := &Round{
			TeamA: TeamFormation{},
			TeamB: TeamFormation{},
		}
		lane := 0
		fillLaneWithIDs(t, &r.TeamA, lane, []int64{11, 12})
		fillLaneWithIDs(t, &r.TeamB, lane, []int64{21})

		// Team A wins if aRoll > bRoll
		rolls := []int{1, 6} // bRoll=1, aRoll=6 -> Team B pops
		idx := 0
		rollFn := func() int {
			val := rolls[idx]
			idx++
			return val
		}

		res := r.ResolveLane(lane, &TestEngine{RollFn: rollFn})
		odize.AssertEqual(t, 2, res.RemainingPlayers)
		odize.AssertEqual(t, 0, r.TeamB.LaneCount(lane))
	})

	group.Test("Team B should win when Team A runs out of players", func(t *testing.T) {
		r := &Round{
			TeamA: TeamFormation{},
			TeamB: TeamFormation{},
		}
		lane := 0
		fillLaneWithIDs(t, &r.TeamA, lane, []int64{31})
		fillLaneWithIDs(t, &r.TeamB, lane, []int64{41, 42})

		// Team B wins if bRoll > aRoll
		rolls := []int{6, 1} // bRoll=6, aRoll=1 -> Team A pops
		idx := 0
		rollFn := func() int {
			val := rolls[idx]
			idx++
			return val
		}

		res := r.ResolveLane(lane, &TestEngine{RollFn: rollFn})
		odize.AssertEqual(t, 2, res.RemainingPlayers)
		odize.AssertEqual(t, 0, r.TeamA.LaneCount(lane))
	})

	group.Test("Team A starts with 0 players should lose immediately", func(t *testing.T) {
		r := &Round{
			TeamA: TeamFormation{},
			TeamB: TeamFormation{},
		}
		lane := 0
		fillLaneWithIDs(t, &r.TeamA, lane, []int64{})
		fillLaneWithIDs(t, &r.TeamB, lane, []int64{51, 52, 53})

		res := r.ResolveLane(lane, &TestEngine{RollFn: func() int { return 1 }})

		odize.AssertEqual(t, TeamB, res.Outcome)
		odize.AssertEqual(t, 3, res.RemainingPlayers)
	})

	group.Test("Team B starts with 0 players should lose immediately", func(t *testing.T) {
		r := &Round{
			TeamA: TeamFormation{},
			TeamB: TeamFormation{},
		}
		lane := 0
		fillLaneWithIDs(t, &r.TeamA, lane, []int64{61, 62, 63})
		fillLaneWithIDs(t, &r.TeamB, lane, []int64{})

		res := r.ResolveLane(lane, &TestEngine{RollFn: func() int { return 1 }})

		odize.AssertEqual(t, TeamA, res.Outcome)
		odize.AssertEqual(t, 3, res.RemainingPlayers)
	})

	err := group.Run()
	odize.AssertNoError(t, err)
}

func TestResolveLanes(t *testing.T) {
	group := odize.NewGroup(t, nil)

	group.Test("Team A should win the round if they have more remaining players across all lanes", func(t *testing.T) {
		r := &Round{
			TeamA: TeamFormation{},
			TeamB: TeamFormation{},
		}
		// Lane 0: Team A wins (2 remaining)
		fillLaneWithIDs(t, &r.TeamA, 0, []int64{101, 102})
		fillLaneWithIDs(t, &r.TeamB, 0, []int64{201})

		// Lane 1: Team B wins (1 remaining)
		fillLaneWithIDs(t, &r.TeamA, 1, []int64{103})
		fillLaneWithIDs(t, &r.TeamB, 1, []int64{202, 203})

		// Lane 2: Team A wins (3 remaining)
		fillLaneWithIDs(t, &r.TeamA, 2, []int64{104, 105, 106})
		fillLaneWithIDs(t, &r.TeamB, 2, []int64{})

		// Rolls for Lane 0: A(6), B(1) -> B loses 1
		// Rolls for Lane 1: A(1), B(6) -> A loses 1
		// Lane 2: No rolls needed as B has 0 players
		rolls := []int{1, 6, 6, 1}
		idx := 0
		rollFn := func() int {
			val := rolls[idx]
			idx++
			return val
		}

		res := r.ResolveLanes(&TestEngine{RollFn: rollFn})

		// Team A players: Lane 0 (2), Lane 1 (0), Lane 2 (3) = 5
		// Team B players: Lane 0 (0), Lane 1 (2), Lane 2 (0) = 2
		// Total A (5) > Total B (2) -> Team A wins
		odize.AssertEqual(t, TeamA, res.Outcome)
		odize.AssertEqual(t, 5, res.RemainingPlayers)
	})

	group.Test("Team B should win the round if they have more remaining players across all lanes", func(t *testing.T) {
		r := &Round{
			TeamA: TeamFormation{},
			TeamB: TeamFormation{},
		}
		// Lane 0: Team B wins (3 remaining)
		fillLaneWithIDs(t, &r.TeamA, 0, []int64{})
		fillLaneWithIDs(t, &r.TeamB, 0, []int64{301, 302, 303})

		// Lane 1: Team B wins (2 remaining)
		fillLaneWithIDs(t, &r.TeamA, 1, []int64{304})
		fillLaneWithIDs(t, &r.TeamB, 1, []int64{401, 402})

		// Lane 2: Team A wins (1 remaining)
		fillLaneWithIDs(t, &r.TeamA, 2, []int64{305})
		fillLaneWithIDs(t, &r.TeamB, 2, []int64{})

		// Rolls for Lane 1: A(1), B(6) -> A loses 1
		rolls := []int{6, 1}
		idx := 0
		rollFn := func() int {
			val := rolls[idx]
			idx++
			return val
		}

		res := r.ResolveLanes(&TestEngine{RollFn: rollFn})

		// Team A players: Lane 0 (0), Lane 1 (0), Lane 2 (1) = 1
		// Team B players: Lane 0 (3), Lane 1 (2), Lane 2 (0) = 5
		// Total B (5) > Total A (1) -> Team B wins
		odize.AssertEqual(t, TeamB, res.Outcome)
		odize.AssertEqual(t, 5, res.RemainingPlayers)
	})

	group.Test("A tie in total remaining players should default to Team B win (or current logic)", func(t *testing.T) {
		r := &Round{
			TeamA: TeamFormation{},
			TeamB: TeamFormation{},
		}
		// Lane 0: Team A wins (1 remaining)
		fillLaneWithIDs(t, &r.TeamA, 0, []int64{501})
		fillLaneWithIDs(t, &r.TeamB, 0, []int64{})

		// Lane 1: Team B wins (1 remaining)
		fillLaneWithIDs(t, &r.TeamA, 1, []int64{})
		fillLaneWithIDs(t, &r.TeamB, 1, []int64{601})

		// Lane 2: Both empty
		fillLaneWithIDs(t, &r.TeamA, 2, []int64{})
		fillLaneWithIDs(t, &r.TeamB, 2, []int64{})

		res := r.ResolveLanes(&TestEngine{RollFn: func() int { return 1 }})

		// Total A (1) == Total B (1)
		// Current logic: if teamAPlayers > teamBPlayers { A wins } else if teamBPlayers > teamAPlayers { B wins } else { Draw }
		// So it should be Draw.
		odize.AssertEqual(t, Draw, res.Outcome)
		odize.AssertEqual(t, 0, res.RemainingPlayers)
	})

	err := group.Run()
	odize.AssertNoError(t, err)
}
