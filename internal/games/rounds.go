package games

import (
	"errors"

	"github.com/code-gorilla-au/rush/internal/augments"
)

func NewRound() Round {
	return Round{
		TeamA:       TeamFormation{Lanes: [3][]int64{}},
		TeamB:       TeamFormation{Lanes: [3][]int64{}},
		DuelResults: []DuelResult{},
	}
}

func (r *Round) FillSquad(a LanesConfig, b LanesConfig) {
	r.TeamA.FillLanes(a)
	r.TeamB.FillLanes(b)
}

func (r *Round) ResolveLanes(rollFn RollStrategy) RoundResult {
	var result []RoundResult

	for lane := 0; lane < len(r.TeamA.Lanes); lane++ {
		laneResult := r.ResolveLane(lane, rollFn)
		result = append(result, laneResult)
	}

	return r.calculateWinner(result)

}

func (r *Round) calculateWinner(result []RoundResult) RoundResult {

	teamAPlayers := 0
	teamBPlayers := 0

	for _, laneResult := range result {

		switch laneResult.Outcome {
		case TeamA:
			teamAPlayers += laneResult.RemainingPlayers
		case TeamB:
			teamBPlayers += laneResult.RemainingPlayers
		}

	}

	if teamAPlayers > teamBPlayers {
		return RoundResult{
			Outcome:          TeamA,
			RemainingPlayers: teamAPlayers,
		}
	}

	if teamAPlayers < teamBPlayers {
		return RoundResult{
			Outcome:          TeamB,
			RemainingPlayers: teamBPlayers,
		}
	}

	return RoundResult{
		Outcome:          Draw,
		RemainingPlayers: 0,
	}
}

func (r *Round) ResolveLane(lane int, rollFn RollStrategy) RoundResult {

	for r.TeamA.LaneHasPlayers(lane) && r.TeamB.LaneHasPlayers(lane) {
		playerA := r.TeamA.LanePeak(lane)
		playerB := r.TeamB.LanePeak(lane)

		var lastRound *DuelResult
		if len(r.DuelResults) > 0 {
			lastRound = new(r.DuelResults[len(r.DuelResults)-1])
		}

		rollInput := DecisionInput{
			lastRound: lastRound,
			teamA: TeamDecisionInput{
				triggeredAugment: augments.NoAugment,
				passivesAugments: r.TeamA.Augments,
				player:           playerA,
				roll:             0,
			},
			teamB: TeamDecisionInput{
				triggeredAugment: augments.NoAugment,
				passivesAugments: r.TeamB.Augments,
				player:           playerB,
				roll:             0,
			},
		}

		re := rollFn.Run(rollInput)
		r.TeamB.Augments = popAugment(r.TeamB.Augments, re.TriggeredAugment)
		r.TeamA.Augments = popAugment(r.TeamA.Augments, re.TriggeredAugment)

		for re.Outcome == Draw {
			re = rollFn.Run(rollInput)
		}

		r.DuelResults = append(r.DuelResults, re)
		if re.Outcome == TeamA {
			_, err := r.TeamB.LanePop(lane)
			if errors.Is(err, ErrNoPlayer) {
				break
			}
		} else {
			_, err := r.TeamA.LanePop(lane)
			if errors.Is(err, ErrNoPlayer) {
				break
			}
		}

	}

	if r.TeamA.LaneHasPlayers(lane) {
		return RoundResult{
			Outcome:          TeamA,
			RemainingPlayers: r.TeamA.LaneCount(lane),
		}
	}

	return RoundResult{
		Outcome:          TeamB,
		RemainingPlayers: r.TeamB.LaneCount(lane),
	}

}

func (s *TeamFormation) LaneCount(lane int) int {
	return len(s.Lanes[lane])
}

func (s *TeamFormation) LaneHasPlayers(lane int) bool {
	return len(s.Lanes[lane]) > 0
}

func (s *TeamFormation) LanePeak(lane int) int64 {
	return s.Lanes[lane][len(s.Lanes[lane])-1]
}

func (s *TeamFormation) LanePop(lane int) (int64, error) {
	tmpLane := s.Lanes[lane]
	if len(tmpLane) == 0 {
		return 0, ErrNoPlayer
	}

	item := tmpLane[len(tmpLane)-1]
	s.Lanes[lane] = tmpLane[:len(tmpLane)-1]

	return item, nil
}

func (s *TeamFormation) FillLanes(f LanesConfig) {
	s.TeamID = f.TeamID
	s.Augments = f.Augments

	remainder := s.LaneFill(0, f.Lane1, f.Players)
	remainder = s.LaneFill(1, f.Lane2, remainder)
	s.LaneFill(2, f.Lane3, remainder)
}

func (s *TeamFormation) LaneFill(lane int, players int, teamPlayers []int64) []int64 {
	remainder := teamPlayers
	for i := 0; i < players; i++ {
		if len(remainder) == 0 {
			break
		}
		player := remainder[0]
		remainder = remainder[1:]
		s.Lanes[lane] = append(s.Lanes[lane], player)
	}

	return remainder
}

func popAugment(list []augments.Effect, name augments.Name) []augments.Effect {

	newList := list

	for i, e := range newList {
		if e.Name == name {
			newList = append(newList[:i], newList[i+1:]...)
			return newList
		}
	}

	return newList
}
