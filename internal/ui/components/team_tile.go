package components

import (
	"fmt"
	"slices"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

type TeamTile struct {
	TeamName     string
	CoachName    string
	PlaybookName string
	PlayerNames  []string
}

func NewTeamTile(teamName, coachName, playbookName string, playerNames []string) TeamTile {
	return TeamTile{
		TeamName:     teamName,
		CoachName:    coachName,
		PlaybookName: playbookName,
		PlayerNames:  slices.Clone(playerNames),
	}
}

func (t TeamTile) View(theme styles.IceTheme, width int) string {
	renderField := func(label, value string) string {
		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			theme.Muted.Render(fmt.Sprintf("%s: ", label)),
			theme.Base.Render(nonEmpty(value, "Unknown")),
		)
	}

	players := make([]string, 0, len(t.PlayerNames))
	for _, playerName := range t.PlayerNames {
		trimmedName := strings.TrimSpace(playerName)
		if trimmedName == "" {
			continue
		}

		players = append(players, theme.Player.Render(fmt.Sprintf("• %s", trimmedName)))
	}

	playerView := theme.Muted.Render("No players")
	if len(players) > 0 {
		playerView = lipgloss.JoinVertical(lipgloss.Left, players...)
	}

	tileBody := lipgloss.JoinVertical(
		lipgloss.Left,
		theme.SecondaryHeader.Render(nonEmpty(t.TeamName, "Unknown Team")),
		renderField("Coach", t.CoachName),
		renderField("Playbook", t.PlaybookName),
		"",
		theme.Muted.Render("Players"),
		playerView,
	)

	tileStyle := theme.RoundBorder
	tileStyle.Align(lipgloss.Left)

	if width > 0 {
		tileStyle = tileStyle.Width(width)
	}

	return tileStyle.Render(tileBody)
}

func nonEmpty(value, fallback string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fallback
	}

	return trimmed
}
