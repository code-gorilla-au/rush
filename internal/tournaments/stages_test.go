package tournaments

import (
	"errors"
	"fmt"
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestGetCurrentStage(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("returns Group when Group stage is active for 4 teams", func(t *testing.T) {
			tournament := Tournament{
				Number: Four,
				Stages: []Stage{
					{Name: StageNameGroup, Status: StageStatusActive},
					{Name: StageNameFinals, Status: StageStatusPending},
				},
			}

			stage, err := GetCurrentStage(tournament)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, StageNameGroup, stage)
		}).
		Test("returns Finals when Group is complete and Finals is pending", func(t *testing.T) {
			tournament := Tournament{
				Number: Four,
				Stages: []Stage{
					{Name: StageNameGroup, Status: StageStatusComplete},
					{Name: StageNameFinals, Status: StageStatusPending},
				},
			}

			stage, err := GetCurrentStage(tournament)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, StageNameFinals, stage)
		}).
		Test("returns error when all stages are complete", func(t *testing.T) {
			tournament := Tournament{
				Number: Four,
				Stages: []Stage{
					{Name: StageNameGroup, Status: StageStatusComplete},
					{Name: StageNameFinals, Status: StageStatusComplete},
				},
			}

			_, err := GetCurrentStage(tournament)
			odize.AssertTrue(t, errors.Is(err, ErrTournamentComplete))
		}).
		Test("returns Group when all stages are pending", func(t *testing.T) {
			tournament := Tournament{
				Number: Eight,
				Stages: []Stage{
					{Name: StageNameGroup, Status: StageStatusPending},
					{Name: StageNameKnock, Status: StageStatusPending},
					{Name: StageNameFinals, Status: StageStatusPending},
				},
			}

			stage, err := GetCurrentStage(tournament)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, StageNameGroup, stage)
		}).
		Test("returns Knockout when Group is complete and Knockout is active", func(t *testing.T) {
			tournament := Tournament{
				Number: Eight,
				Stages: []Stage{
					{Name: StageNameGroup, Status: StageStatusComplete},
					{Name: StageNameKnock, Status: StageStatusActive},
					{Name: StageNameFinals, Status: StageStatusPending},
				},
			}

			stage, err := GetCurrentStage(tournament)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, StageNameKnock, stage)
		}).
		Test("returns error for unknown tournament size", func(t *testing.T) {
			tournament := Tournament{
				Number: NumberOfTeams(10),
				Stages: []Stage{
					{Name: StageNameGroup, Status: StageStatusActive},
				},
			}

			_, err := GetCurrentStage(tournament)
			odize.AssertTrue(t, err != nil)
			odize.AssertEqual(t, "unsupported tournament size: 10", err.Error())
		}).
		Test("returns error when expected stage is missing", func(t *testing.T) {
			tournament := Tournament{
				Number: Four,
				Stages: []Stage{
					{Name: StageNameGroup, Status: StageStatusComplete},
					// StageNameFinals is missing
				},
			}

			_, err := GetCurrentStage(tournament)
			odize.AssertTrue(t, err != nil)
			odize.AssertEqual(t, fmt.Sprintf("expected stage %q not found in tournament stages", StageNameFinals), err.Error())
		}).
		Run()

	odize.AssertNoError(t, err)
}
