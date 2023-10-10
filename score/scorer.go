package score

import (
	"blreynolds4/event-race-timer/results"
	"context"
)

// Scorer is an interface to an object that can score results on the stream.
// It is expected to be called whenever scores need to be updated and process
// as much of the stream as possible on each run.
type Scorer interface {
	ScoreResults(context.Context, *results.ResultStream) error
}
