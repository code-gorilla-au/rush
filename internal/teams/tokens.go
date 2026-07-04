package teams

import (
	"slices"

	"github.com/code-gorilla-au/rush/internal/augments"
)

var _coachPersonaTokens = map[CoachPersona][]augments.Name{
	CoachPersonaVanguard: {
		augments.TwistOfFate,
		augments.Overpower,
		augments.MomentumSurge,
		augments.PrecisionStrike,
	},
	CoachPersonaBastion: {
		augments.SecondChance,
		augments.Hamstring,
		augments.LastStand,
		augments.IceInVeins,
	},
	CoachPersonaTrickster: {
		augments.JammingSignal,
		augments.SecondChance,
		augments.Overpower,
		augments.IceInVeins,
	},
	CoachPersonaWildcard: {
		augments.TwistOfFate,
		augments.Overpower,
		augments.Hamstring,
		augments.PrecisionStrike,
	},
}

func (p CoachPersona) Augments() []augments.Name {
	tokens, ok := _coachPersonaTokens[p]
	if !ok {
		return []augments.Name{}
	}

	return slices.Clone(tokens)
}

func (c *Coach) AvailableAugments() []augments.Name {
	return c.Persona.Augments()
}
