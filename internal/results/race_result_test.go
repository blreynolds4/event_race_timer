package results

import (
	"blreynolds4/event-race-timer/internal/meets"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsComplete(t *testing.T) {
	rr := meets.RaceResult{
		Bib: 1,
		Athlete: &meets.Athlete{
			FirstName: t.Name(),
			LastName:  t.Name(),
			Team:      t.Name(),
			Grade:     1,
			Gender:    "m",
		},
		Place:       1,
		PlaceSource: "y",
	}
	assert.True(t, rr.IsComplete())
}

func TestIsCompleteFalse(t *testing.T) {
	results := []meets.RaceResult{
		{ //Bib is 0
			Athlete: &meets.Athlete{
				FirstName: t.Name(),
				LastName:  t.Name(),
				Team:      t.Name(),
				Grade:     1,
				Gender:    "m",
			},
			Place: 1,
			Time:  time.Second,
		},
		{
			// Athlete is nil
			Bib:   1,
			Place: 1,
			Time:  time.Second,
		},
		{
			// place is 0
			Bib: 1,
			Athlete: &meets.Athlete{
				FirstName: t.Name(),
				LastName:  t.Name(),
				Team:      t.Name(),
				Grade:     1,
				Gender:    "m",
			},
			Time: time.Second,
		},
		{
			// place source is ""
			Bib: 1,
			Athlete: &meets.Athlete{
				FirstName: t.Name(),
				LastName:  t.Name(),
				Team:      t.Name(),
				Grade:     1,
				Gender:    "m",
			},
			Place:        1,
			FinishSource: "y",
		},
	}

	for _, rr := range results {
		assert.False(t, rr.IsComplete(), "Result should have been incomplete", rr)
	}
}
