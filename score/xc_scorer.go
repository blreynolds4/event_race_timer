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
	// keep the current scoring place to assign to the next finisher
	// it only increments if the team has less than 7 results
	// at what point are we done and do we recalculate scores for incomplete teams?
	// do we not include teams in the team result until they have 5?
	// definitely don't show incomplete teams in team results
	// do a brute force run through not worrying about updates and incomplete teams
	// assign scores first, zeros till a team has 5

	// sort incoming results by place?  insertion sort?
	return nil
}
