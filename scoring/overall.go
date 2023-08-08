package scoring

import (
	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/events"
	"time"
)

type ScoredResult struct {
	Athlete competitors.Competitor
	Place   int
	Time    time.Duration
}

type Scorer interface {
	Score(inputEvents events.EventSource, athletes competitors.CompetitorLookup) ([]ScoredResult, error)
}

func NewOverallScoring() Scorer {
	return &overallScorer{}
}

type overallScorer struct {
	eventSource events.EventSource
	athletes    competitors.CompetitorLookup
}

func (os *overallScorer) Score(inputEvents events.EventSource, athletes competitors.CompetitorLookup) ([]ScoredResult, error) {
	return []ScoredResult{}, nil
}
