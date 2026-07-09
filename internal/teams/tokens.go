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

func (p CoachPersona) Augments() []augments.Name {
	aug, ok := _coachPersonaTokens[p]
	if !ok {
		return []augments.Name{}
	}

	return augments.NamesFromCategory(aug)
}

func (c *Coach) AvailableAugments() []augments.Name {
	return c.Persona.Augments()
}
