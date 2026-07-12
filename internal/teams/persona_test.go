package teams

import (
	"slices"
	"testing"

	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/augments"
)

func TestCoachPersona_AvailableTokens(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("returns expected tokens for known persona", func(t *testing.T) {
			tokens := CoachPersonaBastion.Augments()

			odize.AssertTrue(t, slices.ContainsFunc(tokens, func(effect augments.Effect) bool {
				return effect.Name == augments.SecondChance
			}))
			odize.AssertTrue(t, slices.ContainsFunc(tokens, func(effect augments.Effect) bool {
				return effect.Name == augments.LastStand
			}))
			odize.AssertTrue(t, slices.ContainsFunc(tokens, func(effect augments.Effect) bool {
				return effect.Name == augments.Fortify
			}))
			odize.AssertTrue(t, slices.ContainsFunc(tokens, func(effect augments.Effect) bool {
				return effect.Name == augments.Brace
			}))

		}).
		Test("returns empty slice for unknown persona", func(t *testing.T) {
			tokens := CoachPersona("Unknown Coach").Augments()
			odize.AssertTrue(t, len(tokens) == 0)
		}).
		Run()

	odize.AssertNoError(t, err)
}
