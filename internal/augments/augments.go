package augments

import (
	"maps"
	"slices"
)

var (
	_offenseRepo = map[Name]Effect{
		TwistOfFate: {
			Name:     TwistOfFate,
			Category: CategoryOffense,
			Type:     TypeActive,
			Trigger:  TriggerConditionAfterRoll,
			Effect:   "Roll 2d6 and keep the highest.",
			Intent:   "Front-load pressure in must-win lanes.",
			Action:   ActionAddDie,
			Target:   TargetSelf,
			Amount:   1,
		},
		Overpower: {
			Name:     Overpower,
			Category: CategoryOffense,
			Type:     TypeActive,
			Trigger:  TriggerConditionBeforeRoll,
			Effect:   "Gain +1 to your roll total this duel.",
			Intent:   "Reliable low-variance push.",
			Action:   ActionIncrease,
			Target:   TargetSelf,
			Amount:   1,
		},
		PrecisionStrike: {
			Name:     PrecisionStrike,
			Category: CategoryOffense,
			Type:     TypeActive,
			Effect:   "Add +1 to your revealed total.",
			Intent:   "Skill-expression token for close reads.",
			Action:   ActionIncrease,
			Target:   TargetSelf,
		},
		MomentumSurge: {
			Name:     MomentumSurge,
			Category: CategoryOffense,
			Type:     TypePassive,
			Trigger:  TriggerConditionAfterRoll,
			Effect:   "If last duel was a win, gain +2 this duel.",
			Intent:   "Snowball option with explicit condition gate.",
			Action:   ActionIncrease,
			Target:   TargetSelf,
			Amount:   2,
		},
	}

	_defenseRepo = map[Name]Effect{
		Brace: {
			Name:     Brace,
			Category: CategoryDefense,
			Type:     TypeActive,
			Trigger:  TriggerConditionAfterRoll,
			Effect:   "Losing roll by 1, convert result to a tie.",
			Intent:   "Stabilize critical moments after a close loss",
			Action:   ActionResultTie,
			Target:   TargetBoth,
			Amount:   0,
		},
		Fortify: {
			Name:     Fortify,
			Category: CategoryDefense,
			Type:     TypePassive,
			Trigger:  TriggerConditionAfterRoll,
			Effect:   "If last duel was a tie, gain +1 to roll",
			Intent:   "Turn the tides on opponent momentum",
			Action:   ActionIncrease,
			Target:   TargetSelf,
			Amount:   1,
		},
		SecondChance: {
			Name:     SecondChance,
			Category: CategoryDefense,
			Type:     TypePassive,
			Trigger:  TriggerConditionAfterRoll,
			Effect:   "Re-roll your own die once; second result replaces the first.",
			Intent:   "Stabilize critical moments after a miss or tie.",
			Action:   ActionReRoll,
			Target:   TargetSelf,
			Amount:   1,
		},
		LastStand: {
			Name:     LastStand,
			Category: CategoryDefense,
			Type:     TypeActive,
			Trigger:  TriggerConditionAfterRoll,
			Effect:   "If last duel was a loss, and only 1 player remains in the lane, +2 to roll",
			Intent:   "Slow down snowball effects.",
			Action:   ActionIncrease,
			Target:   TargetSelf,
			Amount:   2,
		},
	}

	_sabotageRepo = map[Name]Effect{
		Hamstring: {
			Name:     Hamstring,
			Category: CategorySabotage,
			Type:     TypeActive,
			Trigger:  TriggerConditionAfterRoll,
			Effect:   "Opponent gets -1 to their roll total this duel.",
			Intent:   "Defensive denial and tempo slowdown.",
			Action:   ActionDecrease,
			Target:   TargetOpponent,
			Amount:   1,
		},
		PocketSand: {
			Name:     PocketSand,
			Category: CategorySabotage,
			Type:     TypeActive,
			Trigger:  TriggerConditionBeforeRoll,
			Effect:   "Cancel the opponent's declared Pre-roll augment.",
			Intent:   "Anti-pattern counterplay.",
			Action:   ActionCancel,
			Target:   TargetOpponent,
			Amount:   0,
		},
		PoisonEdge: {
			Name:     PoisonEdge,
			Category: CategorySabotage,
			Type:     TypeActive,
			Trigger:  TriggerConditionAfterRoll,
			Effect:   "If last duel was a loss, and only 1 player remains in the lane, +2 to roll",
			Intent:   "Slow down snowball effects.",
			Action:   ActionIncrease,
			Target:   TargetSelf,
			Amount:   2,
		},
		IceInVeins: {
			Name:     IceInVeins,
			Category: CategorySabotage,
			Type:     TypePassive,
			Trigger:  TriggerConditionAfterAugments,
			Effect:   "Convert tie into a win for your side.",
			Intent:   "Tie-state control and clutch finish potential.",
			Action:   ActionIncrease,
			Target:   TargetSelf,
			Amount:   1,
		},
	}

	_noop = map[Name]Effect{
		NoAugment: {
			Name:     NoAugment,
			Category: CategoryNoOp,
			Type:     TypePassive,
			Trigger:  TriggerConditionBeforeRoll,
			Effect:   "No augment.",
			Intent:   "Standard attack",
			Action:   ActionNoOp,
			Target:   TargetSelf,
			Amount:   0,
		},
	}

	_repositories = []map[Name]Effect{
		_offenseRepo,
		_defenseRepo,
		_sabotageRepo,
		_noop,
	}
)

func Get(name Name) (Effect, bool) {
	for _, repository := range _repositories {
		item, ok := repository[name]
		if ok {
			return item, true
		}
	}

	return Effect{}, false
}

func GetByCategory(category Category) ([]Effect, bool) {
	var repo map[Name]Effect

	switch category {
	case CategoryOffense:
		repo = _offenseRepo
	case CategoryDefense:
		repo = _defenseRepo
	case CategorySabotage:
		repo = _sabotageRepo
	default:
		return []Effect{}, false
	}

	return slices.Collect(maps.Values(repo)), true
}

func NamesFromCategory(category Category) []Name {
	repo, ok := GetByCategory(category)
	if !ok {
		return []Name{}
	}

	var result []Name
	for _, effect := range repo {
		result = append(result, effect.Name)
	}

	return result
}
