package games

import (
	"math/rand/v2"
	"slices"

	"github.com/code-gorilla-au/rush/internal/augments"
)

func DiceRoll() int {
	//nolint:gosec
	return rand.IntN(6) + 1
}

type TeamDecisionInput struct {
	triggeredAugment augments.Name
	passivesAugments []augments.Effect
	player           int64
	roll             int
}

type DecisionInput struct {
	lastRound *DuelResult
	teamA     TeamDecisionInput
	teamB     TeamDecisionInput
}

type DecisionEngineFunc func(input DecisionInput) DecisionInput

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

type Engine struct {
	beforeRoll    []DecisionEngineFunc
	afterRoll     []DecisionEngineFunc
	afterAugments []DecisionEngineFunc
	rollFn        RollFn
}

func NewDecisionEngine() *Engine {
	return &Engine{
		beforeRoll: []DecisionEngineFunc{},
		afterRoll: []DecisionEngineFunc{
			RuleTwistOfFate,
		},
		afterAugments: []DecisionEngineFunc{},
		rollFn:        DiceRoll,
	}
}

func (e *Engine) Run(input DecisionInput) DuelResult {
	for _, ruleFn := range e.beforeRoll {
		input = ruleFn(input)
	}

	input.teamB.roll = e.rollFn()
	input.teamA.roll = e.rollFn()

	for _, ruleFn := range e.afterRoll {
		input = ruleFn(input)
	}

	for _, ruleFn := range e.afterAugments {
		input = ruleFn(input)
	}

	return makeDecision(input)
}

func makeDecision(input DecisionInput) DuelResult {
	if input.teamA.roll == input.teamB.roll {
		return DuelResult{
			Player:           0,
			Outcome:          Draw,
			Roll:             0,
			RollDelta:        0,
			TriggeredAugment: augments.NoAugment,
		}
	}

	if input.teamA.roll > input.teamB.roll {
		return DuelResult{
			Player:           input.teamA.player,
			Outcome:          TeamA,
			Roll:             input.teamA.roll,
			RollDelta:        input.teamA.roll - input.teamB.roll,
			TriggeredAugment: input.teamA.triggeredAugment,
		}
	}

	return DuelResult{
		Player:           input.teamB.player,
		Outcome:          TeamB,
		Roll:             input.teamB.roll,
		RollDelta:        input.teamB.roll - input.teamA.roll,
		TriggeredAugment: input.teamB.triggeredAugment,
	}

}

func canTriggerAugment(activeAugment augments.Name, list []augments.Effect, name augments.Name) bool {
	noActiveAugment := activeAugment == augments.NoAugment || activeAugment == ""

	hasAugmentInList := slices.ContainsFunc(list, func(e augments.Effect) bool {
		return e.Name == name
	})

	return noActiveAugment && hasAugmentInList
}
