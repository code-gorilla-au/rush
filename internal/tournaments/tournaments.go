package tournaments

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

func NewService(deps ServiceDependencies) *Service {
	return &Service{
		gamesSvc: deps.GamesSvc,
		teamsSvc: deps.TeamsSvc,
		store:    deps.Store,
		DB:       deps.DB,
		txnFunc:  deps.TxnFunc,
	}
}

func (s *Service) CreateTournament(ctx context.Context, coachId int64, numberOfTeams NumberOfTeams) error {
	_, err := s.generateGames(ctx, coachId, int64(numberOfTeams))
	if err != nil {
		return fmt.Errorf("generating games for tournament: %w", err)
	}
	return nil
}

func (s *Service) generateGames(ctx context.Context, coachId int64, numberOfTeams int64) ([]games.NewGameParams, error) {
	var tournamentGames []games.NewGameParams

	var totalTeams []teams.AITeam
	h, err := s.getHumanTeam(ctx, coachId)
	if err != nil {
		return tournamentGames, fmt.Errorf("failed to get human team: %w", err)
	}

	totalTeams = append(totalTeams, h)

	aiTeams, err := s.getNonHumanTeams(ctx, numberOfTeams-1)
	if err != nil {
		return tournamentGames, fmt.Errorf("failed to get non-human teams: %w", err)
	}

	totalTeams = append(totalTeams, aiTeams...)

	tournamentGames = generateGameParamsFromTeams(totalTeams, tournamentGames)

	return tournamentGames, nil
}

func (s *Service) getNonHumanTeams(ctx context.Context, aiTeams int64) ([]teams.AITeam, error) {
	var tournamentList []teams.AITeam
	var err error

	tournamentList, err = s.teamsSvc.ListAITeams(ctx)
	if err != nil {
		return tournamentList, fmt.Errorf("failed to list teams: %w", err)
	}

	if int64(len(tournamentList)) < aiTeams {
		return tournamentList, errors.New("not enough teams")
	}

	return tournamentList[0:aiTeams], nil

}

func (s *Service) getHumanTeam(ctx context.Context, coachID int64) (teams.AITeam, error) {
	coach, err := s.teamsSvc.GetCoachByID(ctx, coachID)
	if err != nil {
		return teams.AITeam{}, fmt.Errorf("failed to get coach: %w", err)
	}

	team, err := s.teamsSvc.GetTeamAndPlaybooksByCoachID(ctx, coachID)
	if err != nil {
		return teams.AITeam{}, fmt.Errorf("failed to get team: %w", err)
	}

	return teams.AITeam{
		Coach:    coach,
		Team:     team.Team,
		Playbook: team.Playbooks[0],
	}, nil
}

func generateGameParamsFromTeams(totalTeams []teams.AITeam, tournamentGames []games.NewGameParams) []games.NewGameParams {

	for i := 0; i < len(totalTeams); i++ {

		for j := i + 1; j < len(totalTeams); j++ {
			first := totalTeams[i]
			second := totalTeams[j]

			tournamentGames = append(tournamentGames, games.NewGameParams{
				TeamA: games.TeamConfig{
					TeamID:     first.Team.ID,
					TeamName:   first.Team.Name,
					Players:    teams.GetPlayerIDs(first.Team),
					Augments:   first.Coach.AvailableAugments(),
					Formations: nil,
				},
				TeamB: games.TeamConfig{
					TeamID:     second.Team.ID,
					TeamName:   second.Team.Name,
					Players:    teams.GetPlayerIDs(second.Team),
					Augments:   second.Coach.AvailableAugments(),
					Formations: nil,
				},
			})
		}

	}

	return tournamentGames
}
