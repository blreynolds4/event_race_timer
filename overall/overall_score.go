package overall

import (
	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/results"
	"context"
	"fmt"
	"time"
)

type OverallScorer struct {
	resultChan     chan results.RaceResult
	overallResults []OverallResult
}

type OverallResult struct {
	Athlete    *competitors.Competitor
	Finishtime time.Duration
	Place      int
}

func NewOverallResults() OverallScorer {
	return OverallScorer{
		resultChan:     make(chan results.RaceResult),
		overallResults: make([]OverallResult, 0),
	}
}

func (ovr *OverallScorer) ScoreResults(ctx context.Context, source results.ResultSource) error {
	resultMap := map[int]results.RaceResult{}
	placeMap := map[int]OverallResult{}

	go func() {
		// want to keep trying until told to stop via context
		for {
			var result results.RaceResult
			count, err := source.GetResult(ctx, &result, 5*time.Second)
			if err != nil {
				fmt.Println("error reading result:", err)
			}

			// if we read something send it
			if count == 1 {
				// successful message read
				ovr.resultChan <- result
			}

			// check to stop
			select {
			case <-ctx.Done():
				close(ovr.resultChan)
				fmt.Println("result reader stopping")
				return
			default:
			}
		}
	}()

	// read from the channel and build overall results for each new result
	for newResult := range ovr.resultChan {
		resultMap[newResult.Bib] = newResult

		// convert to overallResult and order into a finish map
		for _, result := range resultMap {
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
	}
	return nil
}
