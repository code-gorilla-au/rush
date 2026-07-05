package games

import (
	"math/rand/v2"

	"github.com/code-gorilla-au/rush/internal/augments"
)

func DiceRoll() int {
	return rand.IntN(6) + 1
}

type TeamDecisionInput struct {
	augment augments.Effect
	roll    int
}

type DecisionInput struct {
	lastRound    DuelResult
	teamAAugment TeamDecisionInput
	teamBAugment TeamDecisionInput
}

type DecisionEngineFunc func(input DecisionInput) DecisionInput

func RuleTwistOfFate(input DecisionInput) DecisionInput {
	return ruleTwistOfFate(input, DiceRoll)
}

func ruleTwistOfFate(input DecisionInput, roll RollFn) DecisionInput {
	if input.teamBAugment.augment.Name == augments.TwistOfFate {
		secondRoll := roll()
		if input.teamBAugment.roll < secondRoll {
			input.teamBAugment.roll = secondRoll
		}
	}

	if input.teamAAugment.augment.Name == augments.TwistOfFate {
		secondRoll := roll()
		if input.teamAAugment.roll < secondRoll {
			input.teamAAugment.roll = secondRoll
		}
	}

	return input
}
