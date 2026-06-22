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
	CoachName  string `faker:"name"`
	TeamName   string `faker:"username"`
	Persona    string
	Formations []playbooks.Formation
}

var personas = []struct {
	Name           string
	FormationNames []string
}{
	{
		Name: "The Wall - A defensive specialist that focuses on protecting all lanes equally.",
		FormationNames: []string{
			"balanced-right", "balanced-left", "split-balanced",
			"strong-centre", "strong-right", "strong-left",
			"overload-centre-left", "overload-centre-right", "balanced-right", "balanced-left",
		},
	},
	{
		Name: "Blitzkrieg - Aggressive tactics focusing on overwhelming the opponent through the center.",
		FormationNames: []string{
			"strong-centre", "single-lane-centre", "overload-centre-left",
			"overload-centre-right", "balanced-left", "balanced-right",
			"strong-centre", "single-lane-centre", "strong-centre", "split-balanced",
		},
	},
	{
		Name: "Flank Specialist - Prefers to attack from the edges, leaving the middle open.",
		FormationNames: []string{
			"split-right", "split-left", "overload-right",
			"overload-left", "split-balanced", "balanced-right",
			"balanced-left", "split-right", "split-left", "overload-right",
		},
	},
	{
		Name: "Right Side Powerhouse - Concentrates most of the strength on the right flank.",
		FormationNames: []string{
			"strong-right", "overload-right", "single-lane-right",
			"balanced-right", "overload-centre-right", "split-right",
			"strong-right", "overload-right", "single-lane-right", "balanced-right",
		},
	},
	{
		Name: "Left Side Powerhouse - Concentrates most of the strength on the left flank.",
		FormationNames: []string{
			"strong-left", "overload-left", "single-lane-left",
			"balanced-left", "overload-centre-left", "split-left",
			"strong-left", "overload-left", "single-lane-left", "balanced-left",
		},
	},
	{
		Name: "The Juggernaut - Heavy focus on one single lane to break through.",
		FormationNames: []string{
			"single-lane-left", "single-lane-centre", "single-lane-right",
			"overload-left", "overload-right", "strong-centre",
			"single-lane-left", "single-lane-centre", "single-lane-right", "split-balanced",
		},
	},
	{
		Name: "Balanced Strategist - Adapts to the game with a mix of balanced formations.",
		FormationNames: []string{
			"balanced-right", "balanced-left", "split-balanced",
			"strong-centre", "strong-right", "strong-left",
			"split-right", "split-left", "balanced-right", "balanced-left",
		},
	},
	{
		Name: "Chaos Theory - Uses unpredictable and highly skewed formations.",
		FormationNames: []string{
			"overload-centre-left", "overload-centre-right", "split-right",
			"split-left", "overload-left", "overload-right",
			"single-lane-left", "single-lane-right", "overload-centre-left", "overload-centre-right",
		},
	},
	{
		Name: "The Shield - Focuses on heavy central defense and slight flank pressure.",
		FormationNames: []string{
			"strong-centre", "overload-centre-left", "overload-centre-right",
			"balanced-right", "balanced-left", "split-balanced",
			"strong-centre", "overload-centre-left", "overload-centre-right", "strong-centre",
		},
	},
	{
		Name: "Centrist - Always stays in the middle, daring the opponent to go around.",
		FormationNames: []string{
			"strong-centre", "single-lane-centre", "balanced-right",
			"balanced-left", "overload-centre-left", "overload-centre-right",
			"strong-centre", "single-lane-centre", "strong-centre", "single-lane-centre",
		},
	},
	{
		Name: "Dual Flanker - Strong presence on both left and right, ignoring the center.",
		FormationNames: []string{
			"split-balanced", "split-right", "split-left",
			"overload-right", "overload-left", "balanced-right",
			"balanced-left", "split-balanced", "split-right", "split-left",
		},
	},
	{
		Name: "Overload Master - Specializes in shifting maximum weight to one side or the other.",
		FormationNames: []string{
			"overload-right", "overload-left", "overload-centre-left",
			"overload-centre-right", "single-lane-left", "single-lane-right",
			"overload-right", "overload-left", "overload-centre-left", "overload-centre-right",
		},
	},
}

type AITeamService struct {
	teamsSvc     TeamCreator
	playbooksSvc PlaybookCreator
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

	aiTeam, err := s.teamsSvc.CreateTeam(ctx, team.TeamName, coach.ID, false)
	if err != nil {
		return fmt.Errorf("creating team: %w", err)
	}

	if _, err = s.playbooksSvc.CreatePlaybook(ctx, playbooks.PlaybookParams{
		TeamID:      aiTeam.ID,
		Name:        fmt.Sprintf("%s Playbook", aiTeam.Name),
		Description: team.Persona,
		Formations:  team.Formations,
	}); err != nil {
		return fmt.Errorf("creating AI playbook: %w", err)
	}

	return nil
}

func generateAITeams() ([]AITeam, error) {
	aiTeams := make([]AITeam, totalTeams)
	allFormations := playbooks.Formations()
	formationMap := make(map[string]playbooks.Formation)
	for _, f := range allFormations {
		formationMap[f.Name] = f
	}

	for i := 0; i < totalTeams; i++ {
		tmpTeam := AITeam{}

		if err := faker.FakeData(&tmpTeam); err != nil {
			return nil, fmt.Errorf("generating team: %w", err)
		}

		persona := personas[i%len(personas)]
		tmpTeam.Persona = persona.Name

		teamFormations := make([]playbooks.Formation, 0, len(persona.FormationNames))
		for _, name := range persona.FormationNames {
			if f, ok := formationMap[name]; ok {
				teamFormations = append(teamFormations, f)
			}
		}

		tmpTeam.Formations = teamFormations

		aiTeams[i] = tmpTeam
	}

	return aiTeams, nil
}
