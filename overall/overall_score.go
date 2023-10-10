package overall

import (
	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/results"
	"context"
	"fmt"
	"time"
)

const resultChunkSize = 25

type OverallScorer struct {
	resultChan     chan results.RaceResult
	overallResults []OverallResult
	rawResults     map[int]results.RaceResult // bib to result, keep latest result
}

type OverallResult struct {
	Athlete    *competitors.Competitor
	Finishtime time.Duration
	Place      int
}

func NewOverallResults() OverallScorer {
	return OverallScorer{
		resultChan:     make(chan results.RaceResult, resultChunkSize),
		overallResults: make([]OverallResult, 0),
		rawResults:     make(map[int]results.RaceResult),
	}
}

func (ovr *OverallScorer) ScoreResults(ctx context.Context, source *results.ResultStream) error {
	placeMap := make(map[int]OverallResult)

	// want to keep trying until told to stop via context
	results := make([]results.RaceResult, resultChunkSize)
	resultCount, err := source.GetResults(ctx, results)
	if err != nil {
		return err
	}

	// get new results until the stream is empty
	for resultCount > 0 {
		// add any new results read to the raw storage
		for i := 0; i < resultCount; i++ {
			newResult := results[i]
			ovr.rawResults[newResult.Bib] = newResult
		}

		resultCount, err = source.GetResults(ctx, results)
		if err != nil {
			return err
		}
	}

	// the stream is empty
	// build the output
	for _, result := range ovr.rawResults {
		placeMap[result.Place] = OverallResult{Athlete: result.Athlete, Place: result.Place, Finishtime: result.Time}
	}

	//output in the correct order
	fmt.Printf("\n\n\n")
	ovr.overallResults = make([]OverallResult, 0)
	fmt.Println("Place Name                             Grade Team                             Time")
	fmt.Println("===== ================================ ===== ================================ ========")
	for i := 1; i <= len(placeMap); i++ {
		ovr.overallResults = append(ovr.overallResults, placeMap[i])
		r := placeMap[i]
		fmt.Printf("%5d %-32s %-5d %-32s %-8s\n", r.Place, r.Athlete.Name, r.Athlete.Grade, r.Athlete.Team, r.Finishtime)
	}
	return nil
}
