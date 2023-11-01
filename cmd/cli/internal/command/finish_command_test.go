package command

import (
	"blreynolds4/event-race-timer/internal/raceevents"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFinishCommandNoBib(t *testing.T) {
	inputEvents := &raceevents.MockEventStream{
		Events: make([]raceevents.Event, 0),
	}

	eventSource := t.Name()
	place := NewFinishCommand(eventSource, inputEvents)
	q, err := place.Run([]string{})
	assert.NoError(t, err)
	assert.False(t, q)
	assert.Equal(t, 1, len(inputEvents.Events))

	fe, ok := inputEvents.Events[0].Data.(raceevents.FinishEvent)
	assert.True(t, ok)
	assert.Equal(t, raceevents.NoBib, fe.Bib)
	assert.Equal(t, eventSource, fe.Source)
}

func TestFinishCommandWithBib(t *testing.T) {
	inputEvents := &raceevents.MockEventStream{
		Events: make([]raceevents.Event, 0),
	}

	eventSource := t.Name()
	place := NewFinishCommand(eventSource, inputEvents)
	q, err := place.Run([]string{"1"})
	assert.NoError(t, err)
	assert.False(t, q)
	assert.Equal(t, 1, len(inputEvents.Events))

	fe, ok := inputEvents.Events[0].Data.(raceevents.FinishEvent)
	assert.True(t, ok)
	assert.Equal(t, 1, fe.Bib)
	assert.Equal(t, eventSource, fe.Source)
}

func TestFinishCommandWithBibAndTimeString(t *testing.T) {
	inputEvents := &raceevents.MockEventStream{
		Events: make([]raceevents.Event, 0),
	}

	eventSource := t.Name()
	expTime := time.Now().UTC().Add(time.Minute)
	place := NewFinishCommand(eventSource, inputEvents)
	q, err := place.Run([]string{"1", expTime.Format(time.RFC3339Nano)})
	assert.NoError(t, err)
	assert.False(t, q)
	assert.Equal(t, 1, len(inputEvents.Events))

	fe, ok := inputEvents.Events[0].Data.(raceevents.FinishEvent)
	assert.True(t, ok)
	assert.Equal(t, 1, fe.Bib)
	assert.Equal(t, eventSource, fe.Source)
	assert.Equal(t, expTime, fe.FinishTime)
}

func TestFinishCommandWithBibAndBadTimeString(t *testing.T) {
	inputEvents := &raceevents.MockEventStream{
		Events: make([]raceevents.Event, 0),
	}

	eventSource := t.Name()
	place := NewFinishCommand(eventSource, inputEvents)
	q, err := place.Run([]string{"1", "bad time"})
	assert.Error(t, err)
	assert.False(t, q)
}

func TestFinishCommandWithBadBib(t *testing.T) {
	inputEvents := &raceevents.MockEventStream{
		Events: make([]raceevents.Event, 0),
	}

	place := NewFinishCommand(t.Name(), inputEvents)
	q, err := place.Run([]string{"x"})
	assert.Error(t, err)
	assert.False(t, q)
	assert.Equal(t, 0, len(inputEvents.Events))
}
