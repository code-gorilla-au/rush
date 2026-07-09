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

		odize.AssertEqual(t, 5, res.teamA.roll)
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

		odize.AssertEqual(t, 4, res.teamA.roll)
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

		odize.AssertEqual(t, 5, res.teamB.roll)
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

		odize.AssertEqual(t, 6, res.teamA.roll)
		odize.AssertEqual(t, 3, res.teamB.roll)
	})

	err := group.Run()
	odize.AssertNoError(t, err)
}

func TestEngineRun(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should trigger configured rules in run pipeline", func(t *testing.T) {
			beforeCalled := false
			afterCalled := false
			afterAugmentsCalled := false

			engine := &Engine{
				beforeRoll: []DecisionEngineFunc{
					func(input DecisionInput) DecisionInput {
						beforeCalled = true
						input.teamA.player = 11
						return input
					},
				},
				afterRoll: []DecisionEngineFunc{
					func(input DecisionInput) DecisionInput {
						afterCalled = true
						input.teamA.roll++
						return input
					},
				},
				afterAugments: []DecisionEngineFunc{
					func(input DecisionInput) DecisionInput {
						afterAugmentsCalled = true
						input.teamA.roll++
						return input
					},
				},
				rollFn: newSequentialRollFn([]int{3, 3}),
			}

			result := engine.Run(DecisionInput{})

			odize.AssertTrue(t, beforeCalled)
			odize.AssertTrue(t, afterCalled)
			odize.AssertTrue(t, afterAugmentsCalled)
			odize.AssertEqual(t, TeamA, result.Outcome)
			odize.AssertEqual(t, int64(11), result.Player)
			odize.AssertEqual(t, 5, result.Roll)
			odize.AssertEqual(t, 2, result.RollDelta)
		}).
		Test("should trigger twist of fate effect when rule matches", func(t *testing.T) {
			secondRolls := newSequentialRollFn([]int{6})

			engine := &Engine{
				beforeRoll: []DecisionEngineFunc{},
				afterRoll: []DecisionEngineFunc{
					func(input DecisionInput) DecisionInput {
						return ruleTwistOfFate(input, secondRolls)
					},
				},
				afterAugments: []DecisionEngineFunc{},
				rollFn:        newSequentialRollFn([]int{5, 2}),
			}

			result := engine.Run(DecisionInput{
				lastRound: &DuelResult{Outcome: TeamB},
				teamA: TeamDecisionInput{
					player:           101,
					passivesAugments: []augments.Effect{{Name: augments.TwistOfFate}},
				},
				teamB: TeamDecisionInput{player: 202},
			})

			odize.AssertEqual(t, TeamA, result.Outcome)
			odize.AssertEqual(t, int64(101), result.Player)
			odize.AssertEqual(t, 6, result.Roll)
			odize.AssertEqual(t, 1, result.RollDelta)
		}).
		Test("should not trigger twist of fate effect when rule does not match", func(t *testing.T) {
			secondRolls := newSequentialRollFn([]int{6})

			engine := &Engine{
				beforeRoll: []DecisionEngineFunc{},
				afterRoll: []DecisionEngineFunc{
					func(input DecisionInput) DecisionInput {
						return ruleTwistOfFate(input, secondRolls)
					},
				},
				afterAugments: []DecisionEngineFunc{},
				rollFn:        newSequentialRollFn([]int{5, 2}),
			}

			result := engine.Run(DecisionInput{
				lastRound: &DuelResult{Outcome: TeamB},
				teamA: TeamDecisionInput{
					player:           101,
					triggeredAugment: augments.Brace,
				},
				teamB: TeamDecisionInput{player: 202},
			})

			odize.AssertEqual(t, TeamB, result.Outcome)
			odize.AssertEqual(t, int64(202), result.Player)
			odize.AssertEqual(t, 5, result.Roll)
			odize.AssertEqual(t, 3, result.RollDelta)
		}).
		Test("should return draw when both final rolls are equal", func(t *testing.T) {
			engine := &Engine{
				beforeRoll:    []DecisionEngineFunc{},
				afterRoll:     []DecisionEngineFunc{},
				afterAugments: []DecisionEngineFunc{},
				rollFn:        newSequentialRollFn([]int{4, 4}),
			}

			result := engine.Run(DecisionInput{
				teamA: TeamDecisionInput{player: 101},
				teamB: TeamDecisionInput{player: 202},
			})

			odize.AssertEqual(t, Draw, result.Outcome)
			odize.AssertEqual(t, int64(0), result.Player)
			odize.AssertEqual(t, 0, result.Roll)
			odize.AssertEqual(t, 0, result.RollDelta)
		}).
		Run()

	odize.AssertNoError(t, err)
}

func newSequentialRollFn(rolls []int) RollFn {
	idx := 0

	return func() int {
		if idx >= len(rolls) {
			return rolls[len(rolls)-1]
		}

		value := rolls[idx]
		idx++
		return value
	}
}
