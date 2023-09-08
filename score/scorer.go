package score

import "blreynolds4/event-race-timer/results"

type Scorer interface {
	ScoreResults(results.ResultSource) error
}
