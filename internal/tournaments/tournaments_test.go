package tournaments

import (
	"fmt"
	"testing"

	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/teams"
)

func TestGenerateGameParamsFromTeams(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.Test("generates correct number of unique games and ensures no duplicate pairs for 4 teams", func(t *testing.T) {
		n := 4
		var totalTeams []teams.AITeam
		for i := range n {
			totalTeams = append(totalTeams, teams.AITeam{
				Team:  teams.Team{ID: int64(i + 1), Name: fmt.Sprintf("Team%d", i+1)},
				Coach: teams.Coach{},
			})
		}

		tournamentGames := generateGameParamsFromTeams(totalTeams, nil)

		expectedGames := (n * (n - 1)) / 2
		odize.AssertEqual(t, expectedGames, len(tournamentGames))

		seenPairs := make(map[string]bool)
		for _, game := range tournamentGames {
			pairKey := fmt.Sprintf("%d-%d", game.TeamA.TeamID, game.TeamB.TeamID)
			reverseKey := fmt.Sprintf("%d-%d", game.TeamB.TeamID, game.TeamA.TeamID)

			odize.AssertFalse(t, seenPairs[pairKey])
			odize.AssertFalse(t, seenPairs[reverseKey])

			seenPairs[pairKey] = true
		}
	}).
		Test("handles 0 teams", func(t *testing.T) {
			tournamentGames := generateGameParamsFromTeams([]teams.AITeam{}, nil)
			odize.AssertEqual(t, 0, len(tournamentGames))
		}).
		Test("handles 1 team", func(t *testing.T) {
			totalTeams := []teams.AITeam{
				{Team: teams.Team{ID: 1, Name: "Team1"}, Coach: teams.Coach{}},
			}
			tournamentGames := generateGameParamsFromTeams(totalTeams, nil)
			odize.AssertEqual(t, 0, len(tournamentGames))
		}).
		Test("handles 2 teams", func(t *testing.T) {
			totalTeams := []teams.AITeam{
				{Team: teams.Team{ID: 1, Name: "Team1"}, Coach: teams.Coach{}},
				{Team: teams.Team{ID: 2, Name: "Team2"}, Coach: teams.Coach{}},
			}
			tournamentGames := generateGameParamsFromTeams(totalTeams, nil)
			odize.AssertEqual(t, 1, len(tournamentGames))
		}).
		Run()

	odize.AssertNoError(t, err)
}
