package teams

import (
	"github.com/code-gorilla-au/rush/internal/database"
)

func fromCoachModel(m database.Coach) Coach {
	return Coach{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt.Time,
		UpdatedAt: m.UpdatedAt.Time,
	}
}
