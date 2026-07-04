package augments

import (
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestRepository(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should contain all tokens from proposal", func(t *testing.T) {
			expectedTokens := []Name{
				TwistOfFate,
				SecondChance,
				Overpower,
				Hamstring,
				PrecisionStrike,
				JammingSignal,
				LastStand,
				MomentumSurge,
				IceInVeins,
			}

			for _, name := range expectedTokens {
				_, ok := _repository[name]
				odize.AssertTrue(t, ok)
			}
		}).
		Test("should have correct number of tokens per category", func(t *testing.T) {
			counts := make(map[Category]int)
			for _, effect := range _repository {
				counts[effect.Category]++
			}

			odize.AssertEqual(t, 5, counts[CategoryOffense])
			odize.AssertEqual(t, 2, counts[CategoryDefense])
			odize.AssertEqual(t, 2, counts[CategorySabotage])
		}).
		Run()

	odize.AssertNoError(t, err)
}

func TestGet(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should return effect and true for valid token", func(t *testing.T) {
			effect, ok := Get(TwistOfFate)
			odize.AssertTrue(t, ok)
			odize.AssertEqual(t, TwistOfFate, effect.Name)
			odize.AssertEqual(t, ActionAddDie, effect.Action)
			odize.AssertEqual(t, 1, effect.Amount)
		}).
		Test("should return empty effect and false for non-existent token", func(t *testing.T) {
			effect, ok := Get(Name("NonExistent"))
			odize.AssertFalse(t, ok)
			odize.AssertEqual(t, Effect{}, effect)
		}).
		Test("should return empty effect and false for empty name", func(t *testing.T) {
			effect, ok := Get(Name(""))
			odize.AssertFalse(t, ok)
			odize.AssertEqual(t, Effect{}, effect)
		}).
		Test("should return false for case-insensitive mismatch", func(t *testing.T) {
			_, ok := Get(Name("twist of fate"))
			odize.AssertFalse(t, ok)
		}).
		Run()

	odize.AssertNoError(t, err)
}

func TestList(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should return all token names in repository", func(t *testing.T) {
			names := List()
			odize.AssertEqual(t, len(_repository), len(names))

			for _, name := range names {
				_, ok := _repository[name]
				odize.AssertTrue(t, ok)
			}
		}).
		Run()

	odize.AssertNoError(t, err)
}
