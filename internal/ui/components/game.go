package components

import (
	"fmt"
	"math/rand/v2"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/games"
)

// Game component handles the UI for a single round of a game.
type Game struct {
	game      *games.Game
	teamAName string
	teamBName string
	resolved  bool
	result    games.Result
	rollFn    games.RollFn
	roundComp Round
}

// MsgResolveRound is sent when the round should be resolved.
type MsgResolveRound struct{}

// MsgNextRound is sent when the user wants to proceed to the next round.
type MsgNextRound struct{}

// NewGame creates a new Game component.
func NewGame(game *games.Game, teamAName, teamBName string, rollFn games.RollFn) Game {
	if rollFn == nil {
		rollFn = func() int {
			return rand.IntN(10) + 1 // 1-10
		}
	}

	currentRoundIdx := game.CurrentRound()
	rounds := game.Rounds()
	var currentRound games.Round
	if currentRoundIdx < int64(len(rounds)) {
		currentRound = rounds[currentRoundIdx]
	}

	return Game{
		game:      game,
		teamAName: teamAName,
		teamBName: teamBName,
		rollFn:    rollFn,
		roundComp: NewRound(currentRound, teamAName, teamBName),
	}
}

// Init initializes the Game component with a 1-second pause.
func (g *Game) Init() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return MsgResolveRound{}
	})
}

// Update handles messages for the Game component.
func (g *Game) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case MsgResolveRound:
		g.handleRound()
	case tea.KeyMsg:
		if g.resolved {
			switch msg.String() {
			case "enter":
				return func() tea.Msg {
					return MsgNextRound{}
				}
			}
		}
	}
	return nil
}

func (g *Game) handleRound() {
	if g.resolved {
		return
	}

	res, err := g.game.ResolveRound(g.rollFn)
	if err == nil {
		g.result = res
		g.resolved = true
		// Update the round component with the final state
		currentRoundIdx := g.game.CurrentRound() - 1
		g.roundComp = NewRound(g.game.Rounds()[currentRoundIdx], g.teamAName, g.teamBName)
	}
}

// View renders the Game component.
func (g *Game) View() string {
	roundView := g.roundComp.View()

	roundNum := g.game.CurrentRound()
	if !g.resolved {
		roundNum++
	}

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A5F2F3")).
		Bold(true).
		MarginBottom(1)

	roundInfo := headerStyle.Render(fmt.Sprintf("ROUND %d", roundNum))

	var footer string
	if g.resolved {
		winner := g.teamAName
		if g.result.Outcome == games.ResultTeamB {
			winner = g.teamBName
		} else if g.result.Outcome == games.ResultDraw {
			winner = "Draw"
		}

		winnerStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true).
			MarginTop(1)

		winnerInfo := winnerStyle.Render(fmt.Sprintf("WINNER: %s! (%d players remaining)", winner, g.result.RemainingPlayers))
		prompt := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			MarginTop(1).
			Render("Press Enter for next round...")

		footer = lipgloss.JoinVertical(lipgloss.Center, winnerInfo, prompt)
	} else {
		footer = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			MarginTop(1).
			Render("Resolving...")
	}

	content := lipgloss.JoinVertical(lipgloss.Center,
		roundInfo,
		roundView,
		footer,
	)

	return lipgloss.NewStyle().
		Padding(1).
		Render(content)
}
