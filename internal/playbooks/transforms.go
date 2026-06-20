package playbooks

import (
	"encoding/json"
	"fmt"

	"github.com/code-gorilla-au/rush/internal/database"
)

func fromFormationJSON(j interface{}) ([]Formation, error) {
	var f []Formation
	if err := json.Unmarshal(j.([]byte), &f); err != nil {
		return []Formation{}, fmt.Errorf("failed to unmarshal formation: %w", err)
	}

	return f, nil
}

func fromPlaybookModel(model database.Playbook) (Playbook, error) {
	formations, err := fromFormationJSON(model.Formations)
	if err != nil {
		return Playbook{}, fmt.Errorf("failed to unmarshal formations: %w", err)
	}

	return Playbook{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description.String,
		Formations:  formations,
	}, nil
}

func fromPlaybookModels(models []database.Playbook) ([]Playbook, error) {
	playbooks := make([]Playbook, len(models))

	for i, model := range models {
		playbook, err := fromPlaybookModel(model)
		if err != nil {
			return nil, err
		}
		playbooks[i] = playbook
	}

	return playbooks, nil
}
