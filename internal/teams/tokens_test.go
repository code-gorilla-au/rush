package teams

import (
	"testing"

	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/augments"
)

func TestCoachPersona_AvailableTokens(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("returns expected tokens for known persona", func(t *testing.T) {
			tokens := CoachPersonaBastion.Augments()

			expected := []augments.Name{
				augments.SecondChance,
				augments.Hamstring,
				augments.LastStand,
				augments.JammingSignal,
				augments.IceInVeins,
			}

			odize.AssertEqual(t, expected, tokens)
		}).
		Test("returns cloned slice to prevent mutation leaks", func(t *testing.T) {
			first := CoachPersonaVanguard.Augments()
			second := CoachPersonaVanguard.Augments()

			first[0] = augments.SecondChance

			odize.AssertEqual(t, augments.TwistOfFate, second[0])
		}).
		Test("returns empty slice for unknown persona", func(t *testing.T) {
			tokens := CoachPersona("Unknown Coach").Augments()
			odize.AssertTrue(t, len(tokens) == 0)
		}).
		Run()

	odize.AssertNoError(t, err)
}

func TestCoach_AvailableTokens(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("returns tokens from coach persona", func(t *testing.T) {
			coach := Coach{Persona: CoachPersonaTrickster}

			expected := CoachPersonaTrickster.Augments()
			actual := coach.AvailableAugments()

			odize.AssertEqual(t, expected, actual)
		}).
		Run()

	odize.AssertNoError(t, err)
}
