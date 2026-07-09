package teams

import (
	"github.com/code-gorilla-au/rush/internal/augments"
)

var _coachPersonaTokens = map[CoachPersona]augments.Category{
	CoachPersonaVanguard:  augments.CategoryOffense,
	CoachPersonaBastion:   augments.CategoryDefense,
	CoachPersonaTrickster: augments.CategorySabotage,
	CoachPersonaWildcard:  augments.CategoryOffense,
}

func (p CoachPersona) Augments() []augments.Effect {
	aug, ok := _coachPersonaTokens[p]
	if !ok {
		return []augments.Effect{}
	}

	list, ok := augments.GetByCategory(aug)
	if !ok {
		return []augments.Effect{}
	}

	return list
}

func (c *Coach) AvailableAugments() []augments.Effect {
	return c.Persona.Augments()
}
