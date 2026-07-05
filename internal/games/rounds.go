package games

import (
	"errors"
)

func NewRound() Round {
	return Round{
		TeamA:       TeamFormation{Lanes: [3][]int{}},
		TeamB:       TeamFormation{Lanes: [3][]int{}},
		DuelResults: []DuelResult{},
	}
}

func (r *Round) FillSquad(a LanesConfig, b LanesConfig) {
	r.TeamA.FillLanes(a)
	r.TeamB.FillLanes(b)
}

func (r *Round) ResolveLanes(rollFn RollFn) RoundResult {
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

func (r *Round) ResolveLane(lane int, rollFn RollFn) RoundResult {
	for r.TeamA.LaneHasPlayers(lane) && r.TeamB.LaneHasPlayers(lane) {
		aRoll := rollFn()
		bRoll := rollFn()

		for aRoll == bRoll {
			bRoll = rollFn()
			aRoll = rollFn()
		}

		if aRoll > bRoll {
			r.DuelResults = append(r.DuelResults, DuelResult{
				Outcome:   TeamA,
				Roll:      aRoll,
				RollDelta: aRoll - bRoll,
			})

			_, err := r.TeamB.LanePop(lane)
			if errors.Is(err, ErrNoPlayer) {
				break
			}

		} else if bRoll > aRoll {
			r.DuelResults = append(r.DuelResults, DuelResult{
				Outcome:   TeamB,
				Roll:      bRoll,
				RollDelta: bRoll - aRoll,
			})

			_, err := r.TeamA.LanePop(lane)
			if errors.Is(err, ErrNoPlayer) {
				break
			}
		}

		r.DuelResults = append(r.DuelResults, DuelResult{
			Outcome:   Draw,
			Roll:      0,
			RollDelta: 0,
		})

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

func (s *TeamFormation) LanePop(lane int) (int, error) {
	tmpLane := s.Lanes[lane]
	if len(tmpLane) == 0 {
		return 0, ErrNoPlayer
	}

	s.Lanes[lane] = tmpLane[:len(tmpLane)-1]

	return 1, nil
}

func (s *TeamFormation) FillLanes(f LanesConfig) {
	s.TeamID = f.TeamID

	s.LaneFill(0, f.Lane1)
	s.LaneFill(1, f.Lane2)
	s.LaneFill(2, f.Lane3)
}

func (s *TeamFormation) LaneFill(lane int, players int) {
	for i := 0; i < players; i++ {
		s.Lanes[lane] = append(s.Lanes[lane], i)
	}
}
