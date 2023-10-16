package xc

import (
	"blreynolds4/event-race-timer/results"
	"context"
	"fmt"
	"time"
)

const resultChunkSize = 25

type XCResult struct {
	result results.RaceResult
	Score  int16
}

type XCTeamResult struct {
	Name      string
	TeamScore int16
	Top5Avg   time.Duration
	Top7Avg   time.Duration
	Finishers []XCResult
}

func NewXCScorer() *XCScorer {
	return &XCScorer{
		Results: make([]XCTeamResult, 0),
	}
}

type XCScorer struct {
	Results []XCTeamResult
}

func (xcs *XCScorer) ScoreResults(ctx context.Context, source results.ResultStream) error {
	rawResults := make([]results.RaceResult, 0)
	teams := make(map[string]XCTeamResult)
	Results := make([]results.RaceResult, resultChunkSize)
	ResultCount, err := source.GetResults(ctx, Results)

	if err != nil {
		return err
	}

	for ResultCount > 0 {
		for i := 0; i < ResultCount; i++ {
			rawResults = append(rawResults, Results[i])
		}

		ResultCount, err = source.GetResults(ctx, Results)
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(rawResults); i++ {
		teamResult := XCTeamResult{}
		_, exists := teams[rawResults[i].Athlete.Team]
		result := rawResults[i]
		if exists {
			teamResult = teams[rawResults[i].Athlete.Team]
		} else {
			teamResult.Name = result.Athlete.Team
		}

		teamResult.Finishers = append(teamResult.Finishers, XCResult{result: result, Score: int16(result.Place)})
		if len(teamResult.Finishers) > 6 {
			//average the times using indexs
		} else if len(teamResult.Finishers) > 4 {
			//average the times
			//we can also score here since there are 5
			scoreAcumulator := 0
			for b := 0; b < 4; b++ {
				scoreAcumulator += int(teamResult.Finishers[b].Score)
			}
			teamResult.TeamScore = int16(scoreAcumulator)
		}
	}

	fmt.Printf("\n\n\n")
	xcs.Results = make([]XCTeamResult, 0)
	fmt.Println("Team name		Score")
	fmt.Println("============== ======")
	for team, teamResult := range teams {
		xcs.Results = append(xcs.Results, teamResult)
		fmt.Printf("%-32s %5d", team, teamResult.TeamScore)
	}

	return nil
}
