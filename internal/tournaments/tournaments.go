package tournaments

import (
	"context"
	"database/sql"

	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/teams"
)

type ServiceDependencies struct {
	GamesSvc *games.Service
	TeamsSvc *teams.Service
	Store    Store
	DB       *sql.DB
	TxnFunc  func(db *sql.Tx) Store
}

func NewTournament(deps ServiceDependencies) *Service {
	return &Service{
		gamesSvc: deps.GamesSvc,
		teamsSvc: deps.TeamsSvc,
		store:    deps.Store,
		DB:       deps.DB,
		txnFunc:  deps.TxnFunc,
	}
}

func (s *Service) CreateTournament(ctx context.Context, coachId int64, numberOfTeams NumberOfTeams) error {
	return nil
}

func (s *Service) getTeamsForTournament(ctx context.Context, coachId int64) ([]teams.Team, error) {
	var totalTeams []teams.Team

	humanTeam, err := s.teamsSvc.GetTeamByCoachID(ctx, coachId)
	if err != nil {
		return totalTeams, err
	}

	totalTeams = append(totalTeams, humanTeam)

	return totalTeams, nil
}
