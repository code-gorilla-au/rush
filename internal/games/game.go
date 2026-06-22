package games

import "github.com/code-gorilla-au/rush/internal/playbooks"

func NewGame(teamAPlaybook playbooks.Playbook, teamBPlaybook playbooks.Playbook) Game {
	rounds := [10]Round{}

	for i := range rounds {
		r := NewRound()

		r.FillSquad(
			SquadConfig{
				TeamID: teamAPlaybook.TeamID,
				Lane1:  teamAPlaybook.Formations[i].Lane1,
				Lane2:  teamAPlaybook.Formations[i].Lane2,
				Lane3:  teamAPlaybook.Formations[i].Lane3,
			},
			SquadConfig{
				TeamID: teamBPlaybook.TeamID,
				Lane1:  teamBPlaybook.Formations[i].Lane1,
				Lane2:  teamBPlaybook.Formations[i].Lane2,
				Lane3:  teamBPlaybook.Formations[i].Lane3,
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
