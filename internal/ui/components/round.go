package components

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/games"
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

func (r Round) View() string {
	teamAStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#A5F2F3")).Bold(true).Width(10).Align(lipgloss.Right)
	teamBStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Width(10).Align(lipgloss.Left)
	playerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#87CEEB"))
	separatorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
	laneLabelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).PaddingLeft(2)

	renderPlayers := func(count int, align lipgloss.Position) string {
		dots := strings.TrimSpace(strings.Repeat("● ", count))
		content := playerStyle.Render(dots)
		return lipgloss.NewStyle().Width(10).Align(align).Render(content)
	}

	header := lipgloss.JoinHorizontal(lipgloss.Top,
		teamAStyle.Render(r.teamAName),
		separatorStyle.Render(" | "),
		teamBStyle.Render(r.teamBName),
	)

	divider := separatorStyle.Render(strings.Repeat("-", 10) + "-|-" + strings.Repeat("-", 10))

	renderLane := func(laneNum int) string {
		aPlayers := renderPlayers(len(r.round.TeamA.Lanes[laneNum-1]), lipgloss.Right)
		bPlayers := renderPlayers(len(r.round.TeamB.Lanes[laneNum-1]), lipgloss.Left)
		return lipgloss.JoinHorizontal(lipgloss.Top,
			aPlayers,
			separatorStyle.Render(" | "),
			bPlayers,
			laneLabelStyle.Render(fmt.Sprintf("Lane %d", laneNum)),
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

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#87CEEB")).
		Padding(1).
		Render(view)
}
