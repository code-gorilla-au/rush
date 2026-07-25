package tournaments

import (
	"errors"
	"fmt"
)

var ErrTournamentComplete = errors.New("tournament complete")

var _stageMaps = map[NumberOfTeams][]string{
	Four:  {StageNameGroup, StageNameFinals},
	Eight: {StageNameGroup, StageNameKnock, StageNameFinals},
}

func GetCurrentStage(tournament Tournament) (string, error) {
	expectedStages := getStageNames(tournament.Number)
	if len(expectedStages) == 0 {
		return "", fmt.Errorf("unsupported tournament size: %d", tournament.Number)
	}

	for _, stageName := range expectedStages {
		found := false
		var currentStage Stage
		for _, s := range tournament.Stages {
			if s.Name == stageName {
				found = true
				currentStage = s
				break
			}
		}

		if !found {
			return "", fmt.Errorf("expected stage %q not found in tournament stages", stageName)
		}

		if currentStage.Status != StageStatusComplete {
			return stageName, nil
		}
	}

	return "", ErrTournamentComplete
}

func getStageNames(numberOfTeams NumberOfTeams) []string {
	s, ok := _stageMaps[numberOfTeams]
	if !ok {
		return []string{}
	}

	return s
}
