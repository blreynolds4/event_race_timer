package meets

import (
	"time"
)

// A RaceResult that can be filled in as events arrive and sent when IsComplete() is true
type RaceResult struct {
	Bib          int
	Athlete      *Athlete
	Place        int
	XcPlace      int
	Time         time.Duration
	FinishSource string
	PlaceSource  string
}

func (rr RaceResult) IsComplete() bool {
	// the un-set, zero value for Athlete is nil
	// if the bib is not 0, athlete is not nil, the place is not 0 and has a source
	// then the result can be published
	return (rr.Bib > 0) &&
		(rr.Athlete != nil) &&
		(rr.Place > 0) &&
		(rr.PlaceSource != "")
}
