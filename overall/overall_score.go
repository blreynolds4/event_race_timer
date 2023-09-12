package overall

import (
	"blreynolds4/event-race-timer/results"
	"context"
)

type overallResult struct {
	overallResults []results.RaceResult
}

func (OVR *overallResult) ScoreResults(context.Context, results.ResultSource) error {
	return nil
}
