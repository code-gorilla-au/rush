package games

import (
	"errors"
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
			player := r.TeamA.LanePeak(lane)

			r.DuelResults = append(r.DuelResults, DuelResult{
				Player:    player,
				Outcome:   TeamA,
				Roll:      aRoll,
				RollDelta: aRoll - bRoll,
			})

			_, err := r.TeamB.LanePop(lane)
			if errors.Is(err, ErrNoPlayer) {
				break
			}

		} else if bRoll > aRoll {
			player := r.TeamB.LanePeak(lane)

			r.DuelResults = append(r.DuelResults, DuelResult{
				Player:    player,
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
