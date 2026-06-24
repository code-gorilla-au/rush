package games

type NewGameParams struct {
	TeamA TeamConfig
	TeamB TeamConfig
}

func NewGame(params NewGameParams) Game {
	rounds := [10]Round{}

	for i := range rounds {
		r := NewRound()

		r.FillSquad(
			LanesConfig{
				TeamID: params.TeamA.TeamID,
				Lane1:  params.TeamA.Formations[i].Lane1,
				Lane2:  params.TeamA.Formations[i].Lane2,
				Lane3:  params.TeamA.Formations[i].Lane3,
			},
			LanesConfig{
				TeamID: params.TeamB.TeamID,
				Lane1:  params.TeamB.Formations[i].Lane1,
				Lane2:  params.TeamB.Formations[i].Lane2,
				Lane3:  params.TeamB.Formations[i].Lane3,
			},
		)

		rounds[i] = r
	}

	return Game{
		rounds:       rounds,
		currentRound: 0,
		results:      []Result{},
	}
}

func (g *Game) ResolveRound(roll RollFn) (Result, error) {
	if g.currentRound >= len(g.rounds) {
		return Result{}, ErrNoRounds
	}

	round := g.rounds[g.currentRound]
	result := round.ResolveLanes(roll)

	g.results = append(g.results, result)
	g.currentRound++

	return result, nil
}
