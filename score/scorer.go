package score

import (
	"blreynolds4/event-race-timer/results"
	"context"
)

type Scorer interface {
	ScoreResults(context.Context, results.ResultSource) error
}
