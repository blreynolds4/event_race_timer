package score

import (
	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/results"
	"context"
	"time"
)

type XCResult struct {
	Athlete competitors.Competitor
	Score   int16
}

type XCTeamResult struct {
	Name      string
	TeamScore int16
	Top5Avg   time.Duration
	Top7Avg   time.Duration
	Finishers []XCResult
}

func NewXCScorer() Scorer {
	return &xcScorer{}
}

type xcScorer struct {
	Results map[string]XCTeamResult
}

func (xcs *xcScorer) ScoreResults(ctx context.Context, source results.ResultSource) error {
	return nil
}
