package overall

import (
	"blreynolds4/event-race-timer/results"
	"context"
)

type overallResult struct {
	overallResults []results.RaceResult
}

func (OVR *overallResult) ScoreResults(ctx context.Context, source results.ResultSource) error {

	return nil
}
