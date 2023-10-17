package results

import (
	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/stream"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsComplete(t *testing.T) {
	rr := RaceResult{
		Bib: 1,
		Athlete: &competitors.Competitor{
			Name:  t.Name(),
			Team:  t.Name(),
			Age:   1,
			Grade: 1,
		},
		Place:       1,
		PlaceSource: "y",
	}
	assert.True(t, rr.IsComplete())
}

func TestIsCompleteFalse(t *testing.T) {
	results := []RaceResult{
		{ //Bib is 0
			Athlete: &competitors.Competitor{
				Name:  t.Name(),
				Team:  t.Name(),
				Age:   1,
				Grade: 1,
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
			Athlete: &competitors.Competitor{
				Name:  t.Name(),
				Team:  t.Name(),
				Age:   1,
				Grade: 1,
			},
			Time: time.Second,
		},
		// {
		// 	// Time is zero
		// 	Bib: 1,
		// 	Athlete: &competitors.Competitor{
		// 		Name:  t.Name(),
		// 		Team:  t.Name(),
		// 		Age:   1,
		// 		Grade: 1,
		// 	},
		// 	Place: 1,
		// },
		{
			// place source is ""
			Bib: 1,
			Athlete: &competitors.Competitor{
				Name:  t.Name(),
				Team:  t.Name(),
				Age:   1,
				Grade: 1,
			},
			Place:        1,
			FinishSource: "y",
		},
		// {
		// 	// finish source is ""
		// 	Bib: 1,
		// 	Athlete: &competitors.Competitor{
		// 		Name:  t.Name(),
		// 		Team:  t.Name(),
		// 		Age:   1,
		// 		Grade: 1,
		// 	},
		// 	Place:       1,
		// 	PlaceSource: "y",
		// },
	}

	for _, rr := range results {
		assert.False(t, rr.IsComplete(), "Result should have been incomplete", rr)
	}
}

func TestSendResult(t *testing.T) {
	mock := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	evStream := NewResultStream(mock)

	rr := RaceResult{
		Bib:     1,
		Athlete: competitors.NewCompetitor(t.Name(), t.Name(), 1, 1),
		Place:   1,
		Time:    time.Second,
	}

	err := evStream.SendResult(context.TODO(), rr)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mock.Events))

	actual := make([]RaceResult, 2)
	resultCount, err := evStream.GetResults(context.TODO(), actual)
	assert.NoError(t, err)
	assert.Equal(t, 1, resultCount)
	assert.Equal(t, rr, actual[0])
}

func TestSendResultFails(t *testing.T) {
	expErr := fmt.Errorf("fail")
	mock := &stream.MockStream{
		Send: func(ctx context.Context, sm stream.Message) error {
			return expErr
		},
	}
	evStream := NewResultStream(mock)

	rr := RaceResult{
		Bib:     1,
		Athlete: competitors.NewCompetitor(t.Name(), t.Name(), 1, 1),
		Place:   1,
		Time:    time.Second,
	}

	err := evStream.SendResult(context.TODO(), rr)
	assert.Equal(t, expErr, err)
}

func TestGetResultsEmptyBuffer(t *testing.T) {
	mock := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	evStream := NewResultStream(mock)

	buf := make([]RaceResult, 0)
	count, err := evStream.GetResults(context.TODO(), buf)
	assert.Equal(t, fmt.Errorf("can't get results with zero length buffer"), err)
	assert.Zero(t, count)
}

func TestGetResultsReadFailure(t *testing.T) {
	expErr := fmt.Errorf("boom")

	mock := &stream.MockStream{
		Range: func(ctx context.Context, startId, endId string, msgs []stream.Message) (int, error) {
			return 0, expErr
		},
	}
	evStream := NewResultStream(mock)

	buf := make([]RaceResult, 1)
	count, err := evStream.GetResults(context.TODO(), buf)
	assert.Equal(t, expErr, err)
	assert.Zero(t, count)
}
