package components

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

// CoachWinnerHuman component displays the human coach winner details.
type CoachWinnerHuman struct {
	Team  *teams.Team
	Coach *teams.Coach
}

// NewCoachWinnerHuman creates a new CoachWinnerHuman component.
func NewCoachWinnerHuman(team *teams.Team, coach *teams.Coach) CoachWinnerHuman {
	return CoachWinnerHuman{
		Team:  team,
		Coach: coach,
	}
}

// View renders the CoachWinnerHuman component.
func (c CoachWinnerHuman) View(theme styles.IceTheme) string {
	if c.Team == nil || c.Coach == nil {
		return ""
	}

	winnerHeader := theme.SecondaryHeader.
		Render(fmt.Sprintf("Winner: %s", c.Team.Name))

	coachInfo := theme.CoachName.Render(fmt.Sprintf("%s (Human Coach)", c.Coach.Name))

	players := make([]string, len(c.Team.Players))
	for i, p := range c.Team.Players {
		players[i] = "• " + p.Name
	}
	playerList := strings.Join(players, "\n")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		winnerHeader,
		coachInfo,
		"",
		"Winning Roster:",
		playerList,
	)
}

// CoachWinnerAI component displays the AI coach winner details.
type CoachWinnerAI struct {
	Team  *teams.Team
	Coach *teams.Coach
}

// NewCoachWinnerAI creates a new CoachWinnerAI component.
func NewCoachWinnerAI(team *teams.Team, coach *teams.Coach) CoachWinnerAI {
	return CoachWinnerAI{
		Team:  team,
		Coach: coach,
	}
}

// View renders the CoachWinnerAI component.
func (c CoachWinnerAI) View(theme styles.IceTheme) string {
	if c.Team == nil || c.Coach == nil {
		return ""
	}

	winnerHeader := theme.SecondaryHeader.
		Render(fmt.Sprintf("Winner: %s", c.Team.Name))

	coachInfo := theme.CoachName.Render(fmt.Sprintf("%s (AI Coach)", c.Coach.Name))

	players := make([]string, len(c.Team.Players))
	for i, p := range c.Team.Players {
		players[i] = "• " + p.Name
	}
	playerList := strings.Join(players, "\n")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		winnerHeader,
		coachInfo,
		"",
		"Winning Roster:",
		playerList,
	)
}
