package games

func TeamStatisticsForGames(teamID int64, ga []Game) TeamStatistics {
	var ts TeamStatistics

	for _, g := range ga {
		tmp := TeamStatisticsForGame(teamID, g)

		ts.Wins += tmp.Wins
		ts.Losses += tmp.Losses
		ts.Draws += tmp.Draws
		ts.RoundsWon += tmp.RoundsWon
		ts.RoundsLost += tmp.RoundsLost
	}

	ts.GamesPlayed = len(ga)

	if ts.GamesPlayed == 0 {
		return ts
	}

	ts.WinRate = float64(ts.Wins) / float64(ts.GamesPlayed) * 100
	ts.RoundDifferential = ts.RoundsWon - ts.RoundsLost
	ts.AverageRoundsWon = float64(ts.RoundsWon) / float64(ts.GamesPlayed)
	ts.AverageRoundsLost = float64(ts.RoundsLost) / float64(ts.GamesPlayed)

	return ts
}

func TeamStatisticsForGame(teamID int64, ga Game) TeamStatistics {
	var ts TeamStatistics

	if ga.winner == nil {
		ga.winner = new(int64)
	}

	sideAWins := len(filterResultsByTeam(TeamA, ga.results))
	sideBWins := len(filterResultsByTeam(TeamB, ga.results))

	switch *ga.winner {
	case 0:
		ts.Draws += 1
	case teamID:
		ts.Wins += 1
	default:
		ts.Losses += 1
	}

	teamSide := TeamA
	if ga.teamB == teamID {
		teamSide = TeamB
	}

	switch teamSide {
	case TeamA:
		ts.RoundsWon += sideAWins
		ts.RoundsLost += sideBWins
	default:
		ts.RoundsWon += sideBWins
		ts.RoundsLost += sideAWins
	}

	return ts
}
