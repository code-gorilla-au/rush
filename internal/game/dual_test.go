package game

import (
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestResolveLane(t *testing.T) {
	group := odize.NewGroup(t, nil)

	group.Test("Team A should win when Team B runs out of players", func(t *testing.T) {
		r := &Round{
			TeamA: Squad{},
			TeamB: Squad{},
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
			TeamA: Squad{},
			TeamB: Squad{},
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

	group.Test("Draw should not result in any player being eliminated", func(t *testing.T) {
		r := &Round{
			TeamA: Squad{},
			TeamB: Squad{},
		}
		lane := 0
		r.TeamA.LaneFill(lane, 1)
		r.TeamB.LaneFill(lane, 1)

		// Draw
		rolls := []int{3, 3, 6, 1} // aRoll=3, bRoll=3 -> nothing; then aRoll=6, bRoll=1 -> Team B pops
		idx := 0
		rollFn := func() int {
			val := rolls[idx]
			idx++
			return val
		}

		res := r.ResolveLane(lane, rollFn)

		odize.AssertTrue(t, res.TeamA)
		odize.AssertFalse(t, res.TeamB)
		odize.AssertEqual(t, 1, res.RemainingPlayers)
		odize.AssertEqual(t, 0, r.TeamB.LaneCount(lane))
		odize.AssertEqual(t, 1, r.TeamA.LaneCount(lane))
	})

	group.Test("Team A starts with 0 players should lose immediately", func(t *testing.T) {
		r := &Round{
			TeamA: Squad{},
			TeamB: Squad{},
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
			TeamA: Squad{},
			TeamB: Squad{},
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
