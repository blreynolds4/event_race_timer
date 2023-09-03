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
		Place: 1,
		Time:  time.Second,
	}

	assert.True(t, rr.IsComplete())
}

func TestIsCompleteFalse(t *testing.T) {
	results := []RaceResult{
		{
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
			Bib:   1,
			Place: 1,
			Time:  time.Second,
		},
		{
			Bib: 1,
			Athlete: &competitors.Competitor{
				Name:  t.Name(),
				Team:  t.Name(),
				Age:   1,
				Grade: 1,
			},
			Time: time.Second,
		},
		{
			Bib: 1,
			Athlete: &competitors.Competitor{
				Name:  t.Name(),
				Team:  t.Name(),
				Age:   1,
				Grade: 1,
			},
			Place: 1,
		},
	}

	for _, rr := range results {
		assert.False(t, rr.IsComplete(), "Result should have been incomplete", rr)
	}
}

func TestStreamConstructors(t *testing.T) {
	actualW := NewResultTarget(&stream.MockStream{})
	_, ok := actualW.(*resultTargetStream)
	assert.True(t, ok)
}

func TestSendResult(t *testing.T) {
	mock := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	evStream := NewResultTarget(mock)

	rr := RaceResult{
		Bib:     1,
		Athlete: competitors.NewCompetitor(t.Name(), t.Name(), 1, 1),
		Place:   1,
		Time:    time.Second,
	}

	err := evStream.SendResult(context.TODO(), rr)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mock.Events))

	actual := RaceResult{}
	err = actual.FromStreamMessage(mock.Events[0])
	assert.NoError(t, err)
	assert.Equal(t, rr, actual)
}

func TestSendResultFails(t *testing.T) {
	expErr := fmt.Errorf("fail")
	mock := &stream.MockStream{
		Send: func(ctx context.Context, sm stream.Message) error {
			return expErr
		},
	}
	evStream := NewResultTarget(mock)

	rr := RaceResult{
		Bib:     1,
		Athlete: competitors.NewCompetitor(t.Name(), t.Name(), 1, 1),
		Place:   1,
		Time:    time.Second,
	}

	err := evStream.SendResult(context.TODO(), rr)
	assert.Equal(t, expErr, err)
}
