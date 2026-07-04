package teams

import (
	"slices"

	"github.com/code-gorilla-au/rush/internal/augments"
)

var coachPersonaTokens = map[CoachPersona][]augments.Name{
	CoachPersonaVanguard: {
		augments.TwistOfFate,
		augments.Overpower,
		augments.MomentumSurge,
		augments.PrecisionStrike,
		augments.IceInVeins,
	},
	CoachPersonaBastion: {
		augments.SecondChance,
		augments.Hamstring,
		augments.LastStand,
		augments.JammingSignal,
		augments.IceInVeins,
	},
	CoachPersonaTrickster: {
		augments.JammingSignal,
		augments.PrecisionStrike,
		augments.SecondChance,
		augments.Overpower,
	},
	CoachPersonaWildcard: {
		augments.TwistOfFate,
		augments.SecondChance,
		augments.Overpower,
		augments.Hamstring,
		augments.PrecisionStrike,
		augments.JammingSignal,
		augments.IceInVeins,
	},
}

func (p CoachPersona) Augments() []augments.Name {
	tokens, ok := coachPersonaTokens[p]
	if !ok {
		return []augments.Name{}
	}

	return slices.Clone(tokens)
}

func (c *Coach) AvailableAugments() []augments.Name {
	return c.Persona.Augments()
}
