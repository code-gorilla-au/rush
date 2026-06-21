package tournament

import (
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestGenerateAITeams(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should generate exactly 12 teams with personas and formations", func(t *testing.T) {
			teams, err := generateAITeams()
			odize.AssertNoError(t, err)
			odize.AssertTrue(t, len(teams) == 12)

			for i, team := range teams {
				odize.AssertTrue(t, team.Persona != "")
				odize.AssertTrue(t, len(team.Formations) == 10)

				// Verify personas are assigned in order
				expectedPersona := personas[i%len(personas)].Name
				odize.AssertTrue(t, team.Persona == expectedPersona)
			}
		}).
		Run()

	odize.AssertNoError(t, err)
}
