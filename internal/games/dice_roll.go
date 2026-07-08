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
	player  int64
	roll    int
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
	if input.teamB.augment.Name == augments.TwistOfFate {
		secondRoll := roll()
		if input.teamB.roll < secondRoll {
			input.teamB.roll = secondRoll
		}
	}

	if input.teamA.augment.Name == augments.TwistOfFate {
		secondRoll := roll()
		if input.teamA.roll < secondRoll {
			input.teamA.roll = secondRoll
		}
	}

	return input
}

type Engine struct {
	beforeRole    []DecisionEngineFunc
	afterRole     []DecisionEngineFunc
	afterAugments []DecisionEngineFunc
	rollFn        RollFn
}

func NewDecisionEngine() *Engine {
	return &Engine{
		beforeRole: []DecisionEngineFunc{},
		afterRole: []DecisionEngineFunc{
			RuleTwistOfFate,
		},
		afterAugments: []DecisionEngineFunc{},
		rollFn:        DiceRoll,
	}
}

func (e *Engine) Run(input DecisionInput) DuelResult {
	for _, ruleFn := range e.beforeRole {
		input = ruleFn(input)
	}

	input.teamB.roll = e.rollFn()
	input.teamA.roll = e.rollFn()

	for _, ruleFn := range e.afterRole {
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
			Player:    0,
			Outcome:   Draw,
			Roll:      0,
			RollDelta: 0,
		}
	}

	if input.teamA.roll > input.teamB.roll {
		return DuelResult{
			Player:    input.teamA.player,
			Outcome:   TeamA,
			Roll:      input.teamA.roll,
			RollDelta: input.teamA.roll - input.teamB.roll,
		}
	}

	return DuelResult{
		Player:    input.teamB.player,
		Outcome:   TeamB,
		Roll:      input.teamB.roll,
		RollDelta: input.teamB.roll - input.teamA.roll,
	}

}
