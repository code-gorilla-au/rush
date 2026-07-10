package augments

import (
	"slices"
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestRepository(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should contain all tokens from proposal", func(t *testing.T) {
			expectedTokens := []Name{
				NoAugment,
				TwistOfFate,
				SecondChance,
				Overpower,
				Hamstring,
				PrecisionStrike,
				PocketSand,
				LastStand,
				MomentumSurge,
				IceInVeins,
				Brace,
				Fortify,
				PoisonEdge,
			}

			for _, name := range expectedTokens {
				_, ok := Get(name)
				odize.AssertTrue(t, ok)
			}
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
			odize.AssertEqual(t, 0, effect.Amount)
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

func TestGetByCategory(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should return offense effects for offense category", func(t *testing.T) {
			effects, ok := GetByCategory(CategoryOffense)
			odize.AssertTrue(t, ok)
			odize.AssertEqual(t, 4, len(effects))

			for _, effect := range effects {
				odize.AssertEqual(t, CategoryOffense, effect.Category)
			}
		}).
		Test("should return defense effects for defense category", func(t *testing.T) {
			effects, ok := GetByCategory(CategoryDefense)
			odize.AssertTrue(t, ok)
			odize.AssertEqual(t, 4, len(effects))

			for _, effect := range effects {
				odize.AssertEqual(t, CategoryDefense, effect.Category)
			}
		}).
		Test("should return sabotage effects for sabotage category", func(t *testing.T) {
			effects, ok := GetByCategory(CategorySabotage)
			odize.AssertTrue(t, ok)
			odize.AssertEqual(t, 4, len(effects))

			for _, effect := range effects {
				odize.AssertEqual(t, CategorySabotage, effect.Category)
			}
		}).
		Test("should return empty effects and false for unknown category", func(t *testing.T) {
			effects, ok := GetByCategory(Category("unknown"))
			odize.AssertFalse(t, ok)
			odize.AssertEqual(t, 0, len(effects))
		}).
		Run()

	odize.AssertNoError(t, err)
}

func TestNamesFromCategory(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should return offense names", func(t *testing.T) {
			names := NamesFromCategory(CategoryOffense)
			odize.AssertEqual(t, 4, len(names))
			odize.AssertTrue(t, slices.Contains(names, TwistOfFate))
			odize.AssertTrue(t, slices.Contains(names, Overpower))
			odize.AssertTrue(t, slices.Contains(names, PrecisionStrike))
			odize.AssertTrue(t, slices.Contains(names, MomentumSurge))
		}).
		Test("should return defense names", func(t *testing.T) {
			names := NamesFromCategory(CategoryDefense)
			odize.AssertEqual(t, 4, len(names))
			odize.AssertTrue(t, slices.Contains(names, Brace))
			odize.AssertTrue(t, slices.Contains(names, Fortify))
			odize.AssertTrue(t, slices.Contains(names, SecondChance))
			odize.AssertTrue(t, slices.Contains(names, LastStand))
		}).
		Test("should return sabotage names", func(t *testing.T) {
			names := NamesFromCategory(CategorySabotage)
			odize.AssertEqual(t, 4, len(names))
			odize.AssertTrue(t, slices.Contains(names, Hamstring))
			odize.AssertTrue(t, slices.Contains(names, PocketSand))
			odize.AssertTrue(t, slices.Contains(names, PoisonEdge))
			odize.AssertTrue(t, slices.Contains(names, IceInVeins))
		}).
		Test("should return empty slice for NoOp category", func(t *testing.T) {
			names := NamesFromCategory(CategoryNoOp)
			odize.AssertEqual(t, 0, len(names))
		}).
		Test("should return empty slice for unknown category", func(t *testing.T) {
			names := NamesFromCategory(Category("unknown"))
			odize.AssertEqual(t, 0, len(names))
		}).
		Run()

	odize.AssertNoError(t, err)
}
