package components

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

type Round struct {
	round     games.Round
	teamAName string
	teamBName string
}

func NewRound(round games.Round, teamAName, teamBName string) Round {
	return Round{
		round:     round,
		teamAName: teamAName,
		teamBName: teamBName,
	}
}

func (r Round) View(theme styles.IceTheme) string {
	renderPlayers := func(count int, align lipgloss.Position) string {
		dots := strings.TrimSpace(strings.Repeat("● ", count))
		content := theme.Player.Render(dots)
		return lipgloss.NewStyle().Width(10).Align(align).Render(content)
	}

	header := lipgloss.JoinHorizontal(lipgloss.Top,
		theme.TeamA.Render(r.teamAName),
		theme.Separator.Render(" | "),
		theme.TeamB.Render(r.teamBName),
	)

	divider := theme.Separator.Render(strings.Repeat("-", 10) + "-|-" + strings.Repeat("-", 10))

	renderLane := func(laneNum int) string {
		aPlayers := renderPlayers(len(r.round.TeamA.Lanes[laneNum-1]), lipgloss.Right)
		bPlayers := renderPlayers(len(r.round.TeamB.Lanes[laneNum-1]), lipgloss.Left)
		return lipgloss.JoinHorizontal(lipgloss.Top,
			aPlayers,
			theme.Separator.Render(" | "),
			bPlayers,
			theme.Label.Render(fmt.Sprintf("Lane %d", laneNum)),
		)
	}

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		divider,
		renderLane(1),
		renderLane(2),
		renderLane(3),
	)

	return theme.RoundBorder.Render(view)
}
