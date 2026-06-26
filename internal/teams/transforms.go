package teams

import (
	"github.com/code-gorilla-au/rush/internal/database"
)

func fromCoachModel(m database.Coach) Coach {
	return Coach{
		ID:        m.ID,
		Name:      m.Name,
		IsDefault: m.IsDefault.Bool,
		IsHuman:   m.IsHuman.Bool,
		CreatedAt: m.CreatedAt.Time,
		UpdatedAt: m.UpdatedAt.Time,
	}
}

func fromCoachesModel(m []database.Coach) []Coach {
	coaches := make([]Coach, len(m))

	for i, coach := range m {
		coaches[i] = fromCoachModel(coach)
	}

	return coaches
}

func fromTeamModel(m database.Team, p []database.Player) Team {
	players := make([]Player, len(p))

	for i, player := range p {
		players[i] = fromPlayerModel(player)
	}

	return Team{
		ID:        m.ID,
		Name:      m.Name,
		CoachID:   int(m.CoachID.Int64),
		Players:   players,
		CreatedAt: m.CreatedAt.Time,
		UpdatedAt: m.UpdatedAt.Time,
	}
}

func fromPlayerModel(m database.Player) Player {
	return Player{
		ID:        m.ID,
		Name:      m.Name,
		TeamID:    int(m.TeamID.Int64),
		CreatedAt: m.CreatedAt.Time,
		UpdatedAt: m.UpdatedAt.Time,
	}
}
