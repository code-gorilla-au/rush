package tournament

import (
	"context"
	"fmt"

	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/go-faker/faker/v4"
)

const totalTeams = 12

type AITeam struct {
	CoachName string `faker:"name"`
	TeamName  string `faker:"username"`
}

type AITeamService struct {
	teamsSvc     *teams.Service
	playbooksSvc *playbooks.Service
}

func (s *AITeamService) GenerateTeams(ctx context.Context) error {
	aiTeams, err := generateAITeams()
	if err != nil {
		return fmt.Errorf("generating AI teams: %w", err)
	}

	for _, team := range aiTeams {
		if err = s.generateTeam(ctx, team); err != nil {
			return fmt.Errorf("generating team: %w", err)
		}
	}

	return nil
}

func (s *AITeamService) generateTeam(ctx context.Context, team AITeam) error {
	coach, err := s.teamsSvc.CreateCoach(ctx, teams.CreateCoachParams{
		Name:      team.CoachName,
		IsHuman:   false,
		IsDefault: false,
	})
	if err != nil {
		return fmt.Errorf("creating AI coach: %w", err)
	}

	aiTeam, err := s.teamsSvc.CreateTeam(nil, team.TeamName, coach.ID, false)
	if err != nil {
		return fmt.Errorf("creating team: %w", err)
	}

	if _, err = s.playbooksSvc.CreatePlaybook(ctx, playbooks.PlaybookParams{
		TeamID:      aiTeam.ID,
		Name:        fmt.Sprintf("%s Playbook", aiTeam.Name),
		Description: "AI generated playbook",
		Formations:  []playbooks.Formation{},
	}); err != nil {
		return fmt.Errorf("creating AI playbook: %w", err)
	}

	return nil
}

func generateAITeams() ([]AITeam, error) {
	aiTeams := make([]AITeam, totalTeams)

	for i := 0; i < totalTeams; i++ {
		tmpTeam := AITeam{}

		if err := faker.FakeData(&tmpTeam); err != nil {
			return nil, fmt.Errorf("generating team: %w", err)
		}

		aiTeams[i] = tmpTeam
	}

	return aiTeams, nil
}
