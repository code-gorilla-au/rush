package games

import (
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestTeamStatisticsForGames(t *testing.T) {
	group := odize.NewGroup(t, nil)

	teamID := int64(1)
	otherTeamID := int64(2)

	err := group.
		Test("should return zeroed statistics for empty games list", func(t *testing.T) {
			stats := TeamStatisticsForGames(teamID, []Game{})
			odize.AssertEqual(t, 0, stats.GamesPlayed)
			odize.AssertEqual(t, 0, stats.Wins)
			odize.AssertEqual(t, 0, stats.Losses)
			odize.AssertEqual(t, 0, stats.Draws)
			odize.AssertEqual(t, 0.0, stats.WinRate)
			odize.AssertEqual(t, 0, stats.RoundsWon)
			odize.AssertEqual(t, 0, stats.RoundsLost)
			odize.AssertEqual(t, 0, stats.RoundDifferential)
			odize.AssertEqual(t, 0.0, stats.AverageRoundsWon)
			odize.AssertEqual(t, 0.0, stats.AverageRoundsLost)
		}).
		Test("should calculate statistics correctly for a win", func(t *testing.T) {
			game := Game{
				teamA:  teamID,
				teamB:  otherTeamID,
				winner: &teamID,
				results: []RoundResult{
					{Outcome: TeamA, RemainingPlayers: 1},
				},
			}
			stats := TeamStatisticsForGames(teamID, []Game{game})
			odize.AssertEqual(t, 1, stats.GamesPlayed)
			odize.AssertEqual(t, 1, stats.Wins)
			odize.AssertEqual(t, 0, stats.Losses)
			odize.AssertEqual(t, 0, stats.Draws)
			odize.AssertEqual(t, 1.0, stats.WinRate)
			odize.AssertEqual(t, 1, stats.RoundsWon)
			odize.AssertEqual(t, 0, stats.RoundsLost)
			odize.AssertEqual(t, 1, stats.RoundDifferential)
			odize.AssertEqual(t, 1.0, stats.AverageRoundsWon)
			odize.AssertEqual(t, 0.0, stats.AverageRoundsLost)
		}).
		Test("should calculate statistics correctly for a loss", func(t *testing.T) {
			game := Game{
				teamA:  teamID,
				teamB:  otherTeamID,
				winner: &otherTeamID,
				results: []RoundResult{
					{Outcome: TeamB, RemainingPlayers: 1},
				},
			}
			stats := TeamStatisticsForGames(teamID, []Game{game})
			odize.AssertEqual(t, 1, stats.GamesPlayed)
			odize.AssertEqual(t, 0, stats.Wins)
			odize.AssertEqual(t, 1, stats.Losses)
			odize.AssertEqual(t, 0, stats.Draws)
			odize.AssertEqual(t, 0.0, stats.WinRate)
			odize.AssertEqual(t, 0, stats.RoundsWon)
			odize.AssertEqual(t, 1, stats.RoundsLost)
			odize.AssertEqual(t, -1, stats.RoundDifferential)
			odize.AssertEqual(t, 0.0, stats.AverageRoundsWon)
			odize.AssertEqual(t, 1.0, stats.AverageRoundsLost)
		}).
		Test("should calculate statistics correctly for a draw", func(t *testing.T) {
			winner := int64(0)
			game := Game{
				teamA:  teamID,
				teamB:  otherTeamID,
				winner: &winner,
				results: []RoundResult{
					{Outcome: TeamA, RemainingPlayers: 1},
				},
			}
			stats := TeamStatisticsForGames(teamID, []Game{game})
			odize.AssertEqual(t, 1, stats.GamesPlayed)
			odize.AssertEqual(t, 0, stats.Wins)
			odize.AssertEqual(t, 0, stats.Losses)
			odize.AssertEqual(t, 1, stats.Draws)
			odize.AssertEqual(t, 0.0, stats.WinRate)
			odize.AssertEqual(t, 1, stats.RoundsWon)
			odize.AssertEqual(t, 0, stats.RoundsLost)
			odize.AssertEqual(t, 1, stats.RoundDifferential)
			odize.AssertEqual(t, 1.0, stats.AverageRoundsWon)
			odize.AssertEqual(t, 0.0, stats.AverageRoundsLost)
		}).
		Test("should handle nil winner as draw", func(t *testing.T) {
			game := Game{
				teamA:  teamID,
				teamB:  otherTeamID,
				winner: nil,
				results: []RoundResult{
					{Outcome: TeamA, RemainingPlayers: 1},
				},
			}
			stats := TeamStatisticsForGames(teamID, []Game{game})
			odize.AssertEqual(t, 1, stats.GamesPlayed)
			odize.AssertEqual(t, 0, stats.Wins)
			odize.AssertEqual(t, 0, stats.Losses)
			odize.AssertEqual(t, 1, stats.Draws)
			odize.AssertEqual(t, 0.0, stats.WinRate)
			odize.AssertEqual(t, 1, stats.RoundsWon)
			odize.AssertEqual(t, 0, stats.RoundsLost)
			odize.AssertEqual(t, 1, stats.RoundDifferential)
			odize.AssertEqual(t, 1.0, stats.AverageRoundsWon)
			odize.AssertEqual(t, 0.0, stats.AverageRoundsLost)
		}).
		Run()

	odize.AssertNoError(t, err)
}
