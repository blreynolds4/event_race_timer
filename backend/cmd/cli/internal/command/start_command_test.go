package command

import (
	"blreynolds4/event-race-timer/internal/raceevents"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartCommandNoTime(t *testing.T) {
	inputEvents := &raceevents.MockEventStream{
		Events: make([]raceevents.Event, 0),
	}

	eventSource := t.Name()
	start := NewStartCommand(eventSource, inputEvents)
	// no seed duration arugment
	q, err := start.Run([]string{})
	assert.NoError(t, err)
	assert.False(t, q)
	assert.Equal(t, 1, len(inputEvents.Events))

	se, ok := inputEvents.Events[0].Data.(raceevents.StartEvent)
	assert.True(t, ok)
	startTime := se.StartTime
	assert.False(t, startTime.IsZero())
	assert.Equal(t, eventSource, se.Source)
}

func TestStartCommandWithTime(t *testing.T) {
	inputEvents := &raceevents.MockEventStream{
		Events: make([]raceevents.Event, 0),
	}

	now := time.Now().UTC()

	eventSource := t.Name()
	start := NewStartCommand(eventSource, inputEvents)
	// with duration argument
	q, err := start.Run([]string{time.Minute.String()})
	assert.NoError(t, err)
	assert.False(t, q)
	assert.Equal(t, 1, len(inputEvents.Events))

	se, ok := inputEvents.Events[0].Data.(raceevents.StartEvent)
	assert.True(t, ok)
	assert.True(t, se.StartTime.Before(now))
	assert.Equal(t, eventSource, se.Source)
}

func TestStartCommandWithBadTime(t *testing.T) {
	inputEvents := &raceevents.MockEventStream{
		Events: make([]raceevents.Event, 0),
	}

	start := NewStartCommand(t.Name(), inputEvents)
	// with duration argument
	q, err := start.Run([]string{"bad"})
	assert.Error(t, err)
	assert.False(t, q)
	assert.Equal(t, 0, len(inputEvents.Events))
}
