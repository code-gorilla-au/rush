package games

import "github.com/code-gorilla-au/rush/internal/augments"

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
