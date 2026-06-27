package components

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/teams"
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
func (c CoachAvatar) View(teamStyle lipgloss.Style, coachStyle lipgloss.Style) string {
	if c.Coach == nil || c.Team == nil {
		return ""
	}

	teamName := teamStyle.Render(c.Team.Name)
	coachName := coachStyle.Render(fmt.Sprintf("Coach: %s", c.Coach.Name))

	return lipgloss.JoinVertical(lipgloss.Left, teamName, coachName)
}
