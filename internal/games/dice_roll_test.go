package games

import (
	"testing"

	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/augments"
)

func TestRuleTwistOfFate(t *testing.T) {
	group := odize.NewGroup(t, nil)

	group.Test("should not modify input struct (immutability)", func(t *testing.T) {
		input := DecisionInput{
			teamAAugment: TeamDecisionInput{
				augment: augments.Effect{Name: augments.TwistOfFate},
				roll:    1,
			},
			teamBAugment: TeamDecisionInput{
				augment: augments.Effect{Name: augments.TwistOfFate},
				roll:    1,
			},
		}

		inputCopy := input
		_ = RuleTwistOfFate(input)

		odize.AssertEqual(t, inputCopy, input)
	})

	group.Test("should increase roll for Team A when TwistOfFate provides a higher roll", func(t *testing.T) {
		input := DecisionInput{
			teamAAugment: TeamDecisionInput{
				augment: augments.Effect{Name: augments.TwistOfFate},
				roll:    2,
			},
			teamBAugment: TeamDecisionInput{
				roll: 3,
			},
		}

		mockRoll := func() int { return 5 }
		res := ruleTwistOfFate(input, mockRoll)

		odize.AssertEqual(t, 5, res.teamAAugment.roll)
		odize.AssertEqual(t, 3, res.teamBAugment.roll) // Team B unchanged
	})

	group.Test("should not change roll for Team A when TwistOfFate provides a lower roll", func(t *testing.T) {
		input := DecisionInput{
			teamAAugment: TeamDecisionInput{
				augment: augments.Effect{Name: augments.TwistOfFate},
				roll:    4,
			},
		}

		mockRoll := func() int { return 2 }
		res := ruleTwistOfFate(input, mockRoll)

		odize.AssertEqual(t, 4, res.teamAAugment.roll)
	})

	group.Test("should increase roll for Team B when TwistOfFate provides a higher roll", func(t *testing.T) {
		input := DecisionInput{
			teamBAugment: TeamDecisionInput{
				augment: augments.Effect{Name: augments.TwistOfFate},
				roll:    2,
			},
			teamAAugment: TeamDecisionInput{
				roll: 3,
			},
		}

		mockRoll := func() int { return 5 }
		res := ruleTwistOfFate(input, mockRoll)

		odize.AssertEqual(t, 5, res.teamBAugment.roll)
		odize.AssertEqual(t, 3, res.teamAAugment.roll) // Team A unchanged
	})

	group.Test("should not change rolls when TwistOfFate augment is not present", func(t *testing.T) {
		input := DecisionInput{
			teamAAugment: TeamDecisionInput{
				augment: augments.Effect{Name: "Some Other Effect"},
				roll:    3,
			},
			teamBAugment: TeamDecisionInput{
				roll: 4,
			},
		}

		mockRoll := func() int { return 6 }
		res := ruleTwistOfFate(input, mockRoll)

		odize.AssertEqual(t, 3, res.teamAAugment.roll)
		odize.AssertEqual(t, 4, res.teamBAugment.roll)
	})

	group.Test("should handle both teams having TwistOfFate", func(t *testing.T) {
		input := DecisionInput{
			teamAAugment: TeamDecisionInput{
				augment: augments.Effect{Name: augments.TwistOfFate},
				roll:    2,
			},
			teamBAugment: TeamDecisionInput{
				augment: augments.Effect{Name: augments.TwistOfFate},
				roll:    3,
			},
		}

		rolls := []int{5, 6} // Team B gets 5, Team A gets 6
		idx := 0
		mockRoll := func() int {
			val := rolls[idx]
			idx++
			return val
		}

		res := ruleTwistOfFate(input, mockRoll)

		odize.AssertEqual(t, 6, res.teamAAugment.roll)
		odize.AssertEqual(t, 5, res.teamBAugment.roll)
	})

	err := group.Run()
	odize.AssertNoError(t, err)
}
