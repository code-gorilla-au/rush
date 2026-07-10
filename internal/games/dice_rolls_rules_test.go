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
			lastRound: &DuelResult{Outcome: Draw},
			teamA: TeamDecisionInput{
				triggeredAugment: augments.TwistOfFate,
				roll:             1,
			},
			teamB: TeamDecisionInput{
				triggeredAugment: augments.TwistOfFate,
				roll:             1,
			},
		}

		inputCopy := input
		_ = RuleTwistOfFate(input)

		odize.AssertEqual(t, inputCopy, input)
	})

	group.Test("should increase roll for Team A when TwistOfFate provides a higher roll", func(t *testing.T) {
		input := DecisionInput{
			lastRound: &DuelResult{Outcome: TeamB},
			teamA: TeamDecisionInput{
				triggeredAugment: augments.NoAugment,
				passivesAugments: []augments.Effect{{Name: augments.TwistOfFate}},
				roll:             2,
			},
			teamB: TeamDecisionInput{
				roll: 3,
			},
		}

		mockRoll := func() int { return 5 }
		res := ruleTwistOfFate(input, mockRoll)

		odize.AssertEqual(t, 5, res.teamA.roll) // max(2, 5) + 1
		odize.AssertEqual(t, 3, res.teamB.roll) // Team B unchanged
	})

	group.Test("should not change roll for Team A when TwistOfFate provides a lower roll", func(t *testing.T) {
		input := DecisionInput{
			lastRound: &DuelResult{Outcome: TeamB},
			teamA: TeamDecisionInput{
				passivesAugments: []augments.Effect{{Name: augments.TwistOfFate}},
				roll:             4,
			},
		}

		mockRoll := func() int { return 2 }
		res := ruleTwistOfFate(input, mockRoll)

		odize.AssertEqual(t, 4, res.teamA.roll) // max(4, 2) + 1
	})

	group.Test("should increase roll for Team B when TwistOfFate provides a higher roll", func(t *testing.T) {
		input := DecisionInput{
			lastRound: &DuelResult{Outcome: TeamA},
			teamB: TeamDecisionInput{
				passivesAugments: []augments.Effect{{Name: augments.TwistOfFate}},
				roll:             2,
			},
			teamA: TeamDecisionInput{
				roll: 3,
			},
		}

		mockRoll := func() int { return 5 }
		res := ruleTwistOfFate(input, mockRoll)

		odize.AssertEqual(t, 5, res.teamB.roll) // max(2, 5) + 1
		odize.AssertEqual(t, 3, res.teamA.roll) // Team A unchanged
	})

	group.Test("should not change rolls when TwistOfFate augment is not present", func(t *testing.T) {
		input := DecisionInput{
			lastRound: &DuelResult{Outcome: TeamB},
			teamA: TeamDecisionInput{
				triggeredAugment: "Some Other Effect",
				roll:             3,
			},
			teamB: TeamDecisionInput{
				roll: 4,
			},
		}

		mockRoll := func() int { return 6 }
		res := ruleTwistOfFate(input, mockRoll)

		odize.AssertEqual(t, 3, res.teamA.roll)
		odize.AssertEqual(t, 4, res.teamB.roll)
	})

	group.Test("should handle both teams having TwistOfFate", func(t *testing.T) {
		input := DecisionInput{
			lastRound: &DuelResult{Outcome: TeamB},
			teamA: TeamDecisionInput{
				passivesAugments: []augments.Effect{{Name: augments.TwistOfFate}},
				roll:             2,
			},
			teamB: TeamDecisionInput{
				passivesAugments: []augments.Effect{{Name: augments.TwistOfFate}},
				roll:             3,
			},
		}

		mockRoll := func() int {
			return 6
		}

		res := ruleTwistOfFate(input, mockRoll)

		odize.AssertEqual(t, 6, res.teamA.roll) // max(2, 6) + 1
		odize.AssertEqual(t, 3, res.teamB.roll)
	})

	err := group.Run()
	odize.AssertNoError(t, err)
}
