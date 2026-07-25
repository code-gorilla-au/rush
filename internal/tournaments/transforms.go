package tournaments

import "github.com/code-gorilla-au/rush/internal/database"

func toTournament(tournament database.Tournament, stages []database.Stage) Tournament {
	return Tournament{
		ID:     tournament.ID,
		Name:   tournament.Name,
		Number: NumberOfTeams(tournament.NumberOfTeams),
		Stages: toStages(stages),
	}
}

func toStages(stages []database.Stage) []Stage {
	var stagesTransformed []Stage
	for _, stage := range stages {
		stagesTransformed = append(stagesTransformed, toStage(stage))
	}
	return stagesTransformed
}

func toStage(stage database.Stage) Stage {
	return Stage{
		ID:     stage.ID,
		Name:   stage.Name,
		Status: StageStatus(stage.Status),
	}
}
