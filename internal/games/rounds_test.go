package games

import (
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestResolveLane(t *testing.T) {
	group := odize.NewGroup(t, nil)

	group.Test("Team A should win when Team B runs out of players", func(t *testing.T) {
		r := &Round{
			TeamA: TeamFormation{},
			TeamB: TeamFormation{},
		}
		lane := 0
		r.TeamA.LaneFill(lane, 2)
		r.TeamB.LaneFill(lane, 1)

		// Team A wins if aRoll > bRoll
		rolls := []int{6, 1} // aRoll=6, bRoll=1 -> Team B pops
		idx := 0
		rollFn := func() int {
			val := rolls[idx]
			idx++
			return val
		}

		res := r.ResolveLane(lane, rollFn)

		odize.AssertTrue(t, res.TeamA)
		odize.AssertFalse(t, res.TeamB)
		odize.AssertEqual(t, 2, res.RemainingPlayers)
		odize.AssertEqual(t, 0, r.TeamB.LaneCount(lane))
	})

	group.Test("Team B should win when Team A runs out of players", func(t *testing.T) {
		r := &Round{
			TeamA: TeamFormation{},
			TeamB: TeamFormation{},
		}
		lane := 0
		r.TeamA.LaneFill(lane, 1)
		r.TeamB.LaneFill(lane, 2)

		// Team B wins if bRoll > aRoll
		rolls := []int{1, 6} // aRoll=1, bRoll=6 -> Team A pops
		idx := 0
		rollFn := func() int {
			val := rolls[idx]
			idx++
			return val
		}

		res := r.ResolveLane(lane, rollFn)

		odize.AssertFalse(t, res.TeamA)
		odize.AssertTrue(t, res.TeamB)
		odize.AssertEqual(t, 2, res.RemainingPlayers)
		odize.AssertEqual(t, 0, r.TeamA.LaneCount(lane))
	})

	group.Test("Team A starts with 0 players should lose immediately", func(t *testing.T) {
		r := &Round{
			TeamA: TeamFormation{},
			TeamB: TeamFormation{},
		}
		lane := 0
		r.TeamA.LaneFill(lane, 0)
		r.TeamB.LaneFill(lane, 3)

		res := r.ResolveLane(lane, func() int { return 1 })

		odize.AssertFalse(t, res.TeamA)
		odize.AssertTrue(t, res.TeamB)
		odize.AssertEqual(t, 3, res.RemainingPlayers)
	})

	group.Test("Team B starts with 0 players should lose immediately", func(t *testing.T) {
		r := &Round{
			TeamA: TeamFormation{},
			TeamB: TeamFormation{},
		}
		lane := 0
		r.TeamA.LaneFill(lane, 3)
		r.TeamB.LaneFill(lane, 0)

		res := r.ResolveLane(lane, func() int { return 1 })

		odize.AssertTrue(t, res.TeamA)
		odize.AssertFalse(t, res.TeamB)
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
		r.TeamA.LaneFill(0, 2)
		r.TeamB.LaneFill(0, 1)

		// Lane 1: Team B wins (1 remaining)
		r.TeamA.LaneFill(1, 1)
		r.TeamB.LaneFill(1, 2)

		// Lane 2: Team A wins (3 remaining)
		r.TeamA.LaneFill(2, 3)
		r.TeamB.LaneFill(2, 0)

		// Rolls for Lane 0: A(6), B(1) -> B loses 1
		// Rolls for Lane 1: A(1), B(6) -> A loses 1
		// Lane 2: No rolls needed as B has 0 players
		rolls := []int{6, 1, 1, 6}
		idx := 0
		rollFn := func() int {
			val := rolls[idx]
			idx++
			return val
		}

		res := r.ResolveLanes(rollFn)

		// Team A players: Lane 0 (2), Lane 1 (0), Lane 2 (3) = 5
		// Team B players: Lane 0 (0), Lane 1 (2), Lane 2 (0) = 2
		// Total A (5) > Total B (2) -> Team A wins
		odize.AssertTrue(t, res.TeamA)
		odize.AssertFalse(t, res.TeamB)
		odize.AssertEqual(t, 5, res.RemainingPlayers)
	})

	group.Test("Team B should win the round if they have more remaining players across all lanes", func(t *testing.T) {
		r := &Round{
			TeamA: TeamFormation{},
			TeamB: TeamFormation{},
		}
		// Lane 0: Team B wins (3 remaining)
		r.TeamA.LaneFill(0, 0)
		r.TeamB.LaneFill(0, 3)

		// Lane 1: Team B wins (2 remaining)
		r.TeamA.LaneFill(1, 1)
		r.TeamB.LaneFill(1, 2)

		// Lane 2: Team A wins (1 remaining)
		r.TeamA.LaneFill(2, 1)
		r.TeamB.LaneFill(2, 0)

		// Rolls for Lane 1: A(1), B(6) -> A loses 1
		rolls := []int{1, 6}
		idx := 0
		rollFn := func() int {
			val := rolls[idx]
			idx++
			return val
		}

		res := r.ResolveLanes(rollFn)

		// Team A players: Lane 0 (0), Lane 1 (0), Lane 2 (1) = 1
		// Team B players: Lane 0 (3), Lane 1 (2), Lane 2 (0) = 5
		// Total B (5) > Total A (1) -> Team B wins
		odize.AssertFalse(t, res.TeamA)
		odize.AssertTrue(t, res.TeamB)
		odize.AssertEqual(t, 5, res.RemainingPlayers)
	})

	group.Test("A tie in total remaining players should default to Team B win (or current logic)", func(t *testing.T) {
		r := &Round{
			TeamA: TeamFormation{},
			TeamB: TeamFormation{},
		}
		// Lane 0: Team A wins (1 remaining)
		r.TeamA.LaneFill(0, 1)
		r.TeamB.LaneFill(0, 0)

		// Lane 1: Team B wins (1 remaining)
		r.TeamA.LaneFill(1, 0)
		r.TeamB.LaneFill(1, 1)

		// Lane 2: Both empty
		r.TeamA.LaneFill(2, 0)
		r.TeamB.LaneFill(2, 0)

		res := r.ResolveLanes(func() int { return 1 })

		// Total A (1) == Total B (1)
		// Current logic: if teamAPlayers > teamBPlayers { A wins } else { B wins }
		// So it should be Team B win.
		odize.AssertFalse(t, res.TeamA)
		odize.AssertTrue(t, res.TeamB)
		odize.AssertEqual(t, 1, res.RemainingPlayers)
	})

	err := group.Run()
	odize.AssertNoError(t, err)
}
