package command

import (
	"blreynolds4/event-race-timer/raceevents"
	"blreynolds4/event-race-timer/stream"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartCommandNoTime(t *testing.T) {
	mockInStream := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	eventSource := t.Name()
	start := NewStartCommand(eventSource, inputEvents)
	// no seed duration arugment
	q, err := start.Run([]string{})
	assert.NoError(t, err)
	assert.False(t, q)
	assert.Equal(t, 1, len(mockInStream.Events))
	actualEvents := buildActualResults(mockInStream)

	se, ok := actualEvents[0].Data.(raceevents.StartEvent)
	assert.True(t, ok)
	startTime := se.StartTime
	assert.False(t, startTime.IsZero())
	assert.Equal(t, eventSource, se.Source)
}

func TestStartCommandWithTime(t *testing.T) {
	mockInStream := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	now := time.Now().UTC()

	eventSource := t.Name()
	start := NewStartCommand(eventSource, inputEvents)
	// with duration argument
	q, err := start.Run([]string{time.Minute.String()})
	assert.NoError(t, err)
	assert.False(t, q)
	assert.Equal(t, 1, len(mockInStream.Events))
	actuaEvents := buildActualResults(mockInStream)

	se, ok := actuaEvents[0].Data.(raceevents.StartEvent)
	assert.True(t, ok)
	assert.True(t, se.StartTime.Before(now))
	assert.Equal(t, eventSource, se.Source)
}

func TestStartCommandWithBadTime(t *testing.T) {
	mockInStream := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	start := NewStartCommand(t.Name(), inputEvents)
	// with duration argument
	q, err := start.Run([]string{"bad"})
	assert.Error(t, err)
	assert.False(t, q)
	assert.Equal(t, 0, len(mockInStream.Events))
}
