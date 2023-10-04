package team

import (
	"blreynolds4/event-race-timer/results"
	"context"
	"fmt"
	"time"
)

type teamScorer struct {
	resultChan  chan results.RaceResult
	teamResults []teamResult
}

type teamResult struct {
	team            string
	score           float32
	runnersFinished int
}

func NewTeamResult() teamScorer {
	return teamScorer{
		resultChan:  make(chan results.RaceResult),
		teamResults: make([]teamResult, 0),
	}
}

func (tr *teamScorer) ScoreResults(ctx context.Context, source results.ResultSource) error {
	teamMap := map[string]teamResult{}

	go func() {
		for {
			var result results.RaceResult
			count, err := source.GetResult(ctx, &result, 5*time.Second)
			if err != nil {
				fmt.Println("error reading result:", err)
			}

			if count == 1 {
				tr.resultChan <- result
			}

			select {
			case <-ctx.Done():
				close(tr.resultChan)
				fmt.Println("result reader stopping")
				return
			default:
			}
		}
	}()

	for newResult := range tr.resultChan {
		value, exists := teamMap[newResult.Athlete.Team]

		if exists {
			if value.runnersFinished < 5 {
				value.score += float32(newResult.Place)
			}
			value.runnersFinished++
		} else {
			value.team = newResult.Athlete.Team
			value.score = float32(newResult.Place)
			value.runnersFinished = 1
		}
		fmt.Println("adding team: ", newResult.Athlete.Team)
		teamMap[newResult.Athlete.Team] = value

		tr.teamResults = make([]teamResult, 0)
		for _, value := range teamMap {
			if !(value.runnersFinished < 5) {
				tr.teamResults = append(tr.teamResults, value)
			}
		}
	}

	return nil
}
