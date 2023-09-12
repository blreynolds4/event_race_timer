package overall

import (
	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/results"
	"context"
	"fmt"
)

type overallResults struct {
	overallResults []overallResult
}

type overallResult struct {
	Athlete *competitors.Competitor
	Place   int16
}

func newOverallResults() *overallResults {
	return &overallResults{
		overallResults: make([]overallResult, 0),
	}
}

func (OVR *overallResults) ScoreResults(ctx context.Context, source results.ResultSource) error {
	result, err := source.GetResult(ctx)

	acc := 0

	for (result != results.RaceResult{}) && err == nil && acc < 10 {
		acc++
		fmt.Println(acc)
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", result.Athlete)
		result, err = source.GetResult(ctx)
	}
	return nil
}
