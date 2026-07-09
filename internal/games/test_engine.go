package games

type TestEngine struct {
	RollFn RollFn
}

func (e *TestEngine) Run(input DecisionInput) DuelResult {
	input.teamB.roll = e.RollFn()
	input.teamA.roll = e.RollFn()
	return makeDecision(input)
}
