package games

import (
	"testing"

	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/augments"
)

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

func TestRules(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("RuleTwistOfFate should do nothing if lastRound is nil", func(t *testing.T) {
			input := DecisionInput{
				teamA: TeamDecisionInput{
					roll:             3,
					passivesAugments: []augments.Effect{{Name: augments.TwistOfFate}},
				},
			}
			result := ruleTwistOfFate(input, func() int { return 6 })
			odize.AssertEqual(t, 3, result.teamA.roll)
			odize.AssertEqual(t, augments.Name(""), result.teamA.triggeredAugment)
		}).
		Test("RuleTwistOfFate should trigger for Team B if Team A won last round", func(t *testing.T) {
			input := DecisionInput{
				lastRound: &DuelResult{Outcome: TeamA},
				teamB: TeamDecisionInput{
					roll:             2,
					passivesAugments: []augments.Effect{{Name: augments.TwistOfFate}},
				},
			}
			result := ruleTwistOfFate(input, func() int { return 5 })
			odize.AssertEqual(t, 5, result.teamB.roll)
			odize.AssertEqual(t, augments.TwistOfFate, result.teamB.triggeredAugment)
		}).
		Test("RuleTwistOfFate should not update roll if second roll is lower", func(t *testing.T) {
			input := DecisionInput{
				lastRound: &DuelResult{Outcome: TeamB},
				teamA: TeamDecisionInput{
					roll:             4,
					passivesAugments: []augments.Effect{{Name: augments.TwistOfFate}},
				},
			}
			result := ruleTwistOfFate(input, func() int { return 2 })
			odize.AssertEqual(t, 4, result.teamA.roll)
			odize.AssertEqual(t, augments.TwistOfFate, result.teamA.triggeredAugment)
		}).
		Test("RuleOverPower should increase roll by 1", func(t *testing.T) {
			input := DecisionInput{
				teamA: TeamDecisionInput{
					roll:             3,
					passivesAugments: []augments.Effect{{Name: augments.Overpower}},
				},
				teamB: TeamDecisionInput{
					roll:             3,
					passivesAugments: []augments.Effect{{Name: augments.Overpower}},
				},
			}
			result := RuleOverPower(input)
			odize.AssertEqual(t, 4, result.teamA.roll)
			odize.AssertEqual(t, 4, result.teamB.roll)
			odize.AssertEqual(t, augments.Overpower, result.teamA.triggeredAugment)
			odize.AssertEqual(t, augments.Overpower, result.teamB.triggeredAugment)
		}).
		Test("RuleMomentumSurge should trigger only if won last round", func(t *testing.T) {
			// Team A won last round
			input := DecisionInput{
				lastRound: &DuelResult{Outcome: TeamA},
				teamA: TeamDecisionInput{
					roll:             3,
					passivesAugments: []augments.Effect{{Name: augments.MomentumSurge}},
				},
				teamB: TeamDecisionInput{
					roll:             3,
					passivesAugments: []augments.Effect{{Name: augments.MomentumSurge}},
				},
			}
			result := RuleMomentumSurge(input)
			odize.AssertEqual(t, 4, result.teamA.roll)
			odize.AssertEqual(t, 3, result.teamB.roll) // No change for Team B
			odize.AssertEqual(t, augments.MomentumSurge, result.teamA.triggeredAugment)
			odize.AssertEqual(t, augments.Name(""), result.teamB.triggeredAugment)

			// Team B won last round
			input.lastRound.Outcome = TeamB
			input.teamA.triggeredAugment = ""
			input.teamA.roll = 3
			input.teamB.triggeredAugment = ""
			input.teamB.roll = 3
			result = RuleMomentumSurge(input)
			odize.AssertEqual(t, 3, result.teamA.roll)
			odize.AssertEqual(t, 4, result.teamB.roll)
			odize.AssertEqual(t, augments.Name(""), result.teamA.triggeredAugment)
			odize.AssertEqual(t, augments.MomentumSurge, result.teamB.triggeredAugment)
		}).
		Test("RulePrecisionStrike should increase roll by 2 if roll >= 4", func(t *testing.T) {
			input := DecisionInput{
				teamA: TeamDecisionInput{
					roll:             4,
					passivesAugments: []augments.Effect{{Name: augments.PrecisionStrike}},
				},
				teamB: TeamDecisionInput{
					roll:             3,
					passivesAugments: []augments.Effect{{Name: augments.PrecisionStrike}},
				},
			}
			result := RulePrecisionStrike(input)
			odize.AssertEqual(t, 5, result.teamA.roll)
			odize.AssertEqual(t, 3, result.teamB.roll)
			odize.AssertEqual(t, augments.PrecisionStrike, result.teamA.triggeredAugment)
			odize.AssertEqual(t, augments.Name(""), result.teamB.triggeredAugment)
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
