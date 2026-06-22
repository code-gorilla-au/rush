package game

import "github.com/code-gorilla-au/rush/internal/playbooks"

func NewGame(teamAPlaybook playbooks.Playbook, teamBPlaybook playbooks.Playbook) Game {
	rounds := [10]Round{}

	for i := range rounds {
		r := NewRound()

		r.FillTeams(
			SquadLanes{
				Lane1: teamAPlaybook.Formations[i].Lane1,
				Lane2: teamAPlaybook.Formations[i].Lane2,
				Lane3: teamAPlaybook.Formations[i].Lane3,
			},
			SquadLanes{
				Lane1: teamBPlaybook.Formations[i].Lane1,
				Lane2: teamBPlaybook.Formations[i].Lane2,
				Lane3: teamBPlaybook.Formations[i].Lane3,
			},
		)

		rounds[i] = r
	}

	return Game{
		Rounds: rounds,
	}
}
