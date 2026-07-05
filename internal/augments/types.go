package augments

type Name string

const (
	TwistOfFate     Name = "Twist of Fate"
	SecondChance    Name = "Second Chance"
	Overpower       Name = "Overpower"
	Hamstring       Name = "Hamstring"
	PrecisionStrike Name = "Precision Strike"
	PocketSand      Name = "Jamming Signal"
	LastStand       Name = "Last Stand"
	MomentumSurge   Name = "Momentum Surge"
	IceInVeins      Name = "Ice in Veins"
	Brace           Name = "Brace"
	Fortify         Name = "Fortify"
	CriticalHit     Name = "Critical Hit"
)

type Action string

const (
	ActionIncrease  Action = "increase"
	ActionDecrease  Action = "decrease"
	ActionAddDie    Action = "add_die"
	ActionReRoll    Action = "re_roll"
	ActionCancel    Action = "cancel"
	ActionResultTie Action = "result_tie"
)

type Target string

const (
	TargetSelf     Target = "self"
	TargetOpponent Target = "opponent"
	TargetBoth     Target = "both"
)

type Category string

const (
	CategoryOffense  Category = "offense"
	CategoryDefense  Category = "defense"
	CategorySabotage Category = "sabotage"
)

type Type string

const (
	TypeActive  Type = "active"
	TypePassive Type = "passive"
)

type TriggerCondition string

const (
	TriggerConditionBeforeRoll    TriggerCondition = "before_roll"
	TriggerConditionAfterRoll     TriggerCondition = "after_roll"
	TriggerConditionAfterAugments TriggerCondition = "after_augments"
)

type Effect struct {
	Name     Name             `json:"name,omitempty"`
	Category Category         `json:"category,omitempty"`
	Type     Type             `json:"type,omitempty"`
	Trigger  TriggerCondition `json:"trigger,omitempty"`
	Effect   string           `json:"effect,omitempty"`
	Intent   string           `json:"intent,omitempty"`
	Action   Action           `json:"action,omitempty"`
	Target   Target           `json:"target,omitempty"`
	Amount   int              `json:"amount,omitempty"`
}
