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
			Type:     TypePassive,
			Trigger:  TriggerConditionAfterRoll,
			Effect:   "If last round was a loss, roll 2d6 keep highest.",
			Intent:   "Strong catch-up pressure in must-win rounds.",
			Action:   ActionAddDie,
			Target:   TargetSelf,
			Amount:   0,
		},
		Overpower: {
			Name:     Overpower,
			Category: CategoryOffense,
			Type:     TypePassive,
			Trigger:  TriggerConditionAfterRoll,
			Effect:   "Gain +1 to your roll total this duel.",
			Intent:   "Reliable low-variance baseline.",
			Action:   ActionIncrease,
			Target:   TargetSelf,
			Amount:   1,
		},
		PrecisionStrike: {
			Name:     PrecisionStrike,
			Category: CategoryOffense,
			Type:     TypePassive,
			Trigger:  TriggerConditionAfterRoll,
			Effect:   "If you roll a 4 or higher, gain +1 to your total.",
			Intent:   "High-variance reward for strong rolls.",
			Action:   ActionIncrease,
			Target:   TargetSelf,
			Amount:   1,
		},
		MomentumSurge: {
			Name:     MomentumSurge,
			Category: CategoryOffense,
			Type:     TypePassive,
			Trigger:  TriggerConditionAfterRoll,
			Effect:   "If last duel was a win, gain +1 this duel.",
			Intent:   "Snowball option to maintain tempo.",
			Action:   ActionIncrease,
			Target:   TargetSelf,
			Amount:   1,
		},
	}

	_defenseRepo = map[Name]Effect{
		Brace: {
			Name:     Brace,
			Category: CategoryDefense,
			Type:     TypeActive,
			Trigger:  TriggerConditionAfterAugments,
			Effect:   "Losing roll by 2 or less, convert result to a tie.",
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
			Effect:   "If last duel was a tie, gain +2 to roll",
			Intent:   "Turn the tides on opponent momentum",
			Action:   ActionIncrease,
			Target:   TargetSelf,
			Amount:   2,
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
			Amount:   0,
		},
		LastStand: {
			Name:     LastStand,
			Category: CategoryDefense,
			Type:     TypeActive,
			Trigger:  TriggerConditionAfterRoll,
			Effect:   "If last duel was a loss, add +2 to roll",
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
			Effect:   "If last duel was a loss, opponent has -2 to roll",
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
		CanceledAugment: {
			Name:     CanceledAugment,
			Category: CategoryNoOp,
			Type:     TypePassive,
			Trigger:  TriggerConditionBeforeRoll,
			Effect:   "Canceled augment.",
			Intent:   "Augment was affected by an opponent.",
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
