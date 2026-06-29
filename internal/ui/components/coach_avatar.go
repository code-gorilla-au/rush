package components

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

// CoachAvatar component displays team name and coach name.
type CoachAvatar struct {
	Coach *teams.Coach
	Team  *teams.Team
}

// NewCoachAvatar creates a new CoachAvatar component.
func NewCoachAvatar(coach *teams.Coach, team *teams.Team) CoachAvatar {
	return CoachAvatar{
		Coach: coach,
		Team:  team,
	}
}

// View renders the CoachAvatar component.
func (c CoachAvatar) View(theme styles.IceTheme) string {
	if c.Coach == nil || c.Team == nil {
		return ""
	}

	teamName := theme.CoachTeam.Render(c.Team.Name)
	coachName := theme.CoachName.Render(fmt.Sprintf("Coach: %s", c.Coach.Name))

	return lipgloss.JoinVertical(lipgloss.Left, teamName, coachName)
}
