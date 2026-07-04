package augments

type Name string

const (
	TwistOfFate     Name = "Twist of Fate"
	SecondChance    Name = "Second Chance"
	Overpower       Name = "Power Play"
	Hamstring       Name = "Brace"
	PrecisionStrike Name = "Precision Strike"
	JammingSignal   Name = "Jamming Signal"
	LastStand       Name = "Last Stand"
	MomentumSurge   Name = "Momentum Surge"
	IceInVeins      Name = "Ice in Veins"
	SmokeScreen     Name = "Smoke Screen"
)

type Action string

const (
	ActionIncrease Action = "increase"
	ActionDecrease Action = "decrease"
	ActionAddDie   Action = "add_die"
	ActionReRoll   Action = "re_roll"
	ActionCancel   Action = "cancel"
)

type Target string

const (
	TargetSelf     Target = "self"
	TargetOpponent Target = "opponent"
)

type Category string

const (
	CategoryOffense  Category = "offense"
	CategoryDefense  Category = "defense"
	CategorySabotage Category = "sabotage"
)

type Effect struct {
	Name     Name     `json:"name,omitempty"`
	Category Category `json:"category,omitempty"`
	Effect   string   `json:"effect,omitempty"`
	Intent   string   `json:"intent,omitempty"`
	Action   Action   `json:"action,omitempty"`
	Target   Target   `json:"target,omitempty"`
	Amount   int      `json:"amount,omitempty"`
}
