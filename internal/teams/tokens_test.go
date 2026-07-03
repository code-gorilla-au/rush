package teams

import (
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestCoachPersona_AvailableTokens(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("returns expected tokens for known persona", func(t *testing.T) {
			tokens := CoachPersonaBastion.AvailableTokens()

			expected := []TokenName{
				TokenSecondChance,
				TokenBrace,
				TokenLastStand,
				TokenJammingSignal,
				TokenIceInVeins,
			}

			odize.AssertEqual(t, expected, tokens)
		}).
		Test("returns cloned slice to prevent mutation leaks", func(t *testing.T) {
			first := CoachPersonaVanguard.AvailableTokens()
			second := CoachPersonaVanguard.AvailableTokens()

			first[0] = TokenSmokeScreen

			odize.AssertEqual(t, TokenTwistOfFate, second[0])
		}).
		Test("returns nil for unknown persona", func(t *testing.T) {
			tokens := CoachPersona("Unknown Coach").AvailableTokens()
			odize.AssertTrue(t, tokens == nil)
		}).
		Run()

	odize.AssertNoError(t, err)
}

func TestCoach_AvailableTokens(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("returns tokens from coach persona", func(t *testing.T) {
			coach := Coach{Persona: CoachPersonaTrickster}

			expected := CoachPersonaTrickster.AvailableTokens()
			actual := coach.AvailableTokens()

			odize.AssertEqual(t, expected, actual)
		}).
		Run()

	odize.AssertNoError(t, err)
}
