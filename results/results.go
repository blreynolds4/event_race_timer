package results

import (
	"blreynolds4/event-race-timer/competitors"
	"time"
)

// A RaceResult that can be filled in as events arrive and sent when IsComplete() is true
type RaceResult struct {
	Bib     int
	Athlete competitors.Competitor
	Place   int
	Time    time.Duration
}

func (rr RaceResult) IsComplete() bool {
	// the un-set, zero value for Athlete is nil because Competitor is an interface
	// if the bib is not 0, athlete is not nil, the place is not 0 and there is a duration, the result is complete
	return (rr.Bib > 0) &&
		(rr.Athlete != nil) &&
		(rr.Place > 0) &&
		(rr.Time.Milliseconds() > 0)
}

// ResultTarget is a result publisher.  It makes results available to things like scoring
// that need to look at each result so athletes can see them.
type ResultTarget interface {
	SendResult(rr RaceResult) error
}

// implmement a struct that implements ResultTarget as a Redis Stream
// TODO
