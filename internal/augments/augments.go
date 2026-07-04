package augments

import (
	"maps"
	"slices"
)

var _repository = map[Name]Effect{
	TwistOfFate: {
		Name:     TwistOfFate,
		Category: CategoryOffense,
		Effect:   "Roll 2d6 and keep the highest.",
		Intent:   "Front-load pressure in must-win lanes.",
		Action:   ActionAddDie,
		Target:   TargetSelf,
		Amount:   1,
	},
	SecondChance: {
		Name:     SecondChance,
		Category: CategoryDefense,
		Effect:   "Re-roll your own die once; second result replaces the first.",
		Intent:   "Stabilize critical moments after a miss or tie.",
		Action:   ActionReRoll,
		Target:   TargetSelf,
		Amount:   1,
	},
	Overpower: {
		Name:     Overpower,
		Category: CategoryOffense,
		Effect:   "Gain +1 to your roll total this duel.",
		Intent:   "Reliable low-variance push.",
		Action:   ActionIncrease,
		Target:   TargetSelf,
	},
	Hamstring: {
		Name:     Hamstring,
		Category: CategorySabotage,
		Effect:   "Opponent gets -1 to their roll total this duel.",
		Intent:   "Defensive denial and tempo slowdown.",
		Action:   ActionDecrease,
		Target:   TargetOpponent,
	},
	PrecisionStrike: {
		Name:     PrecisionStrike,
		Category: CategoryOffense,
		Effect:   "Add +1 to your revealed total.",
		Intent:   "Skill-expression token for close reads.",
		Action:   ActionIncrease,
		Target:   TargetSelf,
	},
	JammingSignal: {
		Name:     JammingSignal,
		Category: CategorySabotage,
		Effect:   "Cancel the opponent's declared Pre-roll token.",
		Intent:   "Anti-pattern counterplay.",
		Action:   ActionCancel,
		Target:   TargetOpponent,
	},
	LastStand: {
		Name:     LastStand,
		Category: CategoryDefense,
		Effect:   "Prevent your elimination this duel; lane remains unresolved.",
		Intent:   "Comeback insurance for high-value lanes.",
		Action:   "",
		Target:   TargetSelf,
	},
	MomentumSurge: {
		Name:     MomentumSurge,
		Category: CategoryOffense,
		Effect:   "If last round was a win, gain +2 this duel.",
		Intent:   "Snowball option with explicit condition gate.",
		Action:   ActionIncrease,
		Target:   TargetSelf,
		Amount:   2,
	},
	IceInVeins: {
		Name:     IceInVeins,
		Category: CategoryOffense,
		Effect:   "Convert tie into a win for your side.",
		Intent:   "Tie-state control and clutch finish potential.",
		Action:   ActionIncrease,
		Target:   TargetSelf,
		Amount:   1,
	},
	SmokeScreen: {
		Name:     SmokeScreen,
		Category: CategorySabotage,
		Effect:   "Your token declaration remains hidden until after reveal.",
		Intent:   "Mind-game tool to punish reactive opponents.",
		Action:   "",
		Target:   TargetSelf,
	},
}

func Get(name Name) (Effect, bool) {
	item, ok := _repository[name]
	return item, ok
}

func List() []Name {
	return slices.Collect(maps.Keys(_repository))
}
