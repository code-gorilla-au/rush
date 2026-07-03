package teams

import "slices"

var coachPersonaTokens = map[CoachPersona][]TokenName{
	CoachPersonaVanguard: {
		TokenTwistOfFate,
		TokenPowerPlay,
		TokenMomentumSurge,
		TokenPrecisionStrike,
		TokenIceInVeins,
	},
	CoachPersonaBastion: {
		TokenSecondChance,
		TokenBrace,
		TokenLastStand,
		TokenJammingSignal,
		TokenIceInVeins,
	},
	CoachPersonaTrickster: {
		TokenJammingSignal,
		TokenSmokeScreen,
		TokenPrecisionStrike,
		TokenSecondChance,
		TokenPowerPlay,
	},
	CoachPersonaWildcard: {
		TokenTwistOfFate,
		TokenSecondChance,
		TokenPowerPlay,
		TokenBrace,
		TokenPrecisionStrike,
		TokenJammingSignal,
		TokenIceInVeins,
		TokenSmokeScreen,
	},
}

func (p CoachPersona) AvailableTokens() []TokenName {
	tokens, ok := coachPersonaTokens[p]
	if !ok {
		return []TokenName{}
	}

	return slices.Clone(tokens)
}

func (c *Coach) AvailableTokens() []TokenName {
	return c.Persona.AvailableTokens()
}
