package games

import (
	"math"

	"github.com/code-gorilla-au/rush/internal/augments"
)

func RuleTwistOfFate(input DecisionInput) DecisionInput {
	return ruleTwistOfFate(input, DiceRoll)
}

func ruleTwistOfFate(input DecisionInput, roll RollFn) DecisionInput {
	if input.lastRound == nil {
		return input
	}

	switch input.lastRound.Outcome {
	case TeamB:
		if canTriggerAugment(input.teamA.triggeredAugment, input.teamA.passivesAugments, augments.TwistOfFate) {
			input.teamA.triggeredAugment = augments.TwistOfFate

			secondRoll := roll()
			if input.teamA.roll < secondRoll {
				input.teamA.roll = secondRoll
			}
		}

	case TeamA:
		if canTriggerAugment(input.teamB.triggeredAugment, input.teamB.passivesAugments, augments.TwistOfFate) {
			input.teamB.triggeredAugment = augments.TwistOfFate

			secondRoll := roll()
			if input.teamB.roll < secondRoll {
				input.teamB.roll = secondRoll
			}
		}
	}

	return input
}

func RulePrecisionStrike(input DecisionInput) DecisionInput {
	effect, ok := augments.Get(augments.PrecisionStrike)
	if !ok {
		return input
	}

	if input.teamA.roll >= 4 && canTriggerAugment(input.teamA.triggeredAugment, input.teamA.passivesAugments, augments.PrecisionStrike) {
		input.teamA.triggeredAugment = augments.PrecisionStrike

		input.teamA.roll += effect.Amount
	}

	if input.teamB.roll >= 4 && canTriggerAugment(input.teamB.triggeredAugment, input.teamB.passivesAugments, augments.PrecisionStrike) {
		input.teamB.triggeredAugment = augments.PrecisionStrike

		input.teamB.roll += effect.Amount
	}

	return input
}

func RuleOverPower(input DecisionInput) DecisionInput {
	effect, ok := augments.Get(augments.Overpower)
	if !ok {
		return input
	}

	if canTriggerAugment(input.teamA.triggeredAugment, input.teamA.passivesAugments, augments.Overpower) {
		input.teamA.triggeredAugment = augments.Overpower

		input.teamA.roll += effect.Amount
	}

	if canTriggerAugment(input.teamB.triggeredAugment, input.teamB.passivesAugments, augments.Overpower) {
		input.teamB.triggeredAugment = augments.Overpower

		input.teamB.roll += effect.Amount
	}

	return input
}

func RuleMomentumSurge(input DecisionInput) DecisionInput {
	effect, ok := augments.Get(augments.MomentumSurge)
	if !ok {
		return input
	}

	if input.lastRound == nil {
		return input
	}

	if input.lastRound.Outcome == TeamA && canTriggerAugment(input.teamA.triggeredAugment, input.teamA.passivesAugments, augments.MomentumSurge) {
		input.teamA.triggeredAugment = augments.MomentumSurge

		input.teamA.roll += effect.Amount
	}

	if input.lastRound.Outcome == TeamB && canTriggerAugment(input.teamB.triggeredAugment, input.teamB.passivesAugments, augments.MomentumSurge) {
		input.teamB.triggeredAugment = augments.MomentumSurge

		input.teamB.roll += effect.Amount
	}

	return input
}

func RuleBrace(input DecisionInput) DecisionInput {
	effect, ok := augments.Get(augments.Brace)
	if !ok {
		return input
	}

	rollDelta := math.Abs(float64(input.teamA.roll) - float64(input.teamB.roll))
	lostByOne := rollDelta <= 2

	if !lostByOne {
		return input
	}

	if canTriggerAugment(
		input.teamA.triggeredAugment,
		input.teamA.passivesAugments,
		augments.MomentumSurge,
	) {
		if input.teamA.roll < input.teamB.roll {
			input.teamA.triggeredAugment = augments.Brace
			input.teamA.roll = effect.Amount
			input.teamB.roll = effect.Amount
		}

	}

	if canTriggerAugment(
		input.teamB.triggeredAugment,
		input.teamB.passivesAugments,
		augments.MomentumSurge,
	) {
		if input.teamB.roll < input.teamA.roll {
			input.teamB.triggeredAugment = augments.Brace
			input.teamA.roll = effect.Amount
			input.teamB.roll = effect.Amount
		}
	}

	return input
}

func RuleFortify(input DecisionInput) DecisionInput {
	effect, ok := augments.Get(augments.Fortify)
	if !ok {
		return input
	}

	if input.lastRound == nil {
		return input
	}

	if input.lastRound.Outcome == Draw {
		if canTriggerAugment(input.teamA.triggeredAugment, input.teamA.passivesAugments, augments.Fortify) {
			input.teamA.triggeredAugment = augments.Fortify
			input.teamA.roll += effect.Amount
		}

		if canTriggerAugment(input.teamB.triggeredAugment, input.teamB.passivesAugments, augments.Fortify) {
			input.teamB.triggeredAugment = augments.Fortify
			input.teamB.roll += effect.Amount
		}
	}

	return input
}

func RuleSecondChance(input DecisionInput) DecisionInput {
	return ruleSecondChance(input, DiceRoll)
}

func ruleSecondChance(input DecisionInput, rollFn RollFn) DecisionInput {

	aRoll := input.teamA.roll
	bRoll := input.teamB.roll

	if aRoll <= bRoll {
		if canTriggerAugment(input.teamA.triggeredAugment, input.teamA.passivesAugments, augments.SecondChance) {
			input.teamA.triggeredAugment = augments.SecondChance
			input.teamA.roll = rollFn()
		}
	} else {
		if canTriggerAugment(input.teamB.triggeredAugment, input.teamB.passivesAugments, augments.SecondChance) {
			input.teamB.triggeredAugment = augments.SecondChance
			input.teamB.roll = rollFn()
		}
	}

	return input
}

func LastStand(input DecisionInput) DecisionInput {
	return lastStand(input, DiceRoll)
}

func lastStand(input DecisionInput, roll RollFn) DecisionInput {
	effect, ok := augments.Get(augments.LastStand)
	if !ok {
		return input
	}

	if input.lastRound == nil {
		return input
	}

	switch input.lastRound.Outcome {
	case TeamB:
		if canTriggerAugment(input.teamA.triggeredAugment, input.teamA.passivesAugments, augments.LastStand) {
			input.teamA.triggeredAugment = augments.LastStand
			input.teamA.roll += effect.Amount
		}

	case TeamA:
		if canTriggerAugment(input.teamB.triggeredAugment, input.teamB.passivesAugments, augments.LastStand) {
			input.teamB.triggeredAugment = augments.LastStand
			input.teamB.roll += effect.Amount
		}
	}

	return input
}
