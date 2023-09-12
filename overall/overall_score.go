package overall

import (
	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/results"
	"context"
	"time"
)

type overallResults struct {
	overallResults []overallResult
}

type overallResult struct {
	Athlete *competitors.Competitor
	Ftime   time.Duration
	Place   int16
}

func newOverallResults() *overallResults {
	return &overallResults{
		overallResults: make([]overallResult, 0),
	}
}

func (OVR *overallResults) ScoreResults(ctx context.Context, source results.ResultSource) error {
	result, err := source.GetResult(ctx)

	resultMap := map[int]results.RaceResult{}
	placeMap := map[int16]overallResult{}
	//store all results to bib number, this collects the most recent results
	for (result != results.RaceResult{}) && err == nil {
		resultMap[result.Bib] = result

		result, err = source.GetResult(ctx)
	}
	//convert to overallResult and order into a finish map
	for _, result := range resultMap {
		finish := overallResult{}
		finish.Athlete = result.Athlete
		finish.Place = int16(result.Place)
		finish.Ftime = result.Time

		placeMap[finish.Place] = finish
	}
	//output in the correct order
	for i := int16(1); i <= int16(len(placeMap)); i++ {
		OVR.overallResults = append(OVR.overallResults, placeMap[i])
	}
	return nil
}
