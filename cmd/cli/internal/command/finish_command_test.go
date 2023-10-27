package command

import (
	"blreynolds4/event-race-timer/internal/raceevents"
	"blreynolds4/event-race-timer/internal/stream"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFinishCommandNoBib(t *testing.T) {
	mockInStream := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	eventSource := t.Name()
	place := NewFinishCommand(eventSource, inputEvents)
	q, err := place.Run([]string{})
	assert.NoError(t, err)
	assert.False(t, q)
	assert.Equal(t, 1, len(mockInStream.Events))
	actualEvents := buildActualResults(mockInStream)

	fe, ok := actualEvents[0].Data.(raceevents.FinishEvent)
	assert.True(t, ok)
	assert.Equal(t, raceevents.NoBib, fe.Bib)
	assert.Equal(t, eventSource, fe.Source)
}

func TestFinishCommandWithBib(t *testing.T) {
	mockInStream := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	eventSource := t.Name()
	place := NewFinishCommand(eventSource, inputEvents)
	q, err := place.Run([]string{"1"})
	assert.NoError(t, err)
	assert.False(t, q)
	assert.Equal(t, 1, len(mockInStream.Events))
	actualEvents := buildActualResults(mockInStream)

	fe, ok := actualEvents[0].Data.(raceevents.FinishEvent)
	assert.True(t, ok)
	assert.Equal(t, 1, fe.Bib)
	assert.Equal(t, eventSource, fe.Source)
}

func TestFinishCommandWithBibAndTimeString(t *testing.T) {
	mockInStream := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	eventSource := t.Name()
	expTime := time.Now().UTC().Add(time.Minute)
	place := NewFinishCommand(eventSource, inputEvents)
	q, err := place.Run([]string{"1", expTime.Format(time.RFC3339Nano)})
	assert.NoError(t, err)
	assert.False(t, q)
	assert.Equal(t, 1, len(mockInStream.Events))
	actualEvents := buildActualResults(mockInStream)

	fe, ok := actualEvents[0].Data.(raceevents.FinishEvent)
	assert.True(t, ok)
	assert.Equal(t, 1, fe.Bib)
	assert.Equal(t, eventSource, fe.Source)
	assert.Equal(t, expTime, fe.FinishTime)
}

func TestFinishCommandWithBibAndBadTimeString(t *testing.T) {
	mockInStream := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	eventSource := t.Name()
	place := NewFinishCommand(eventSource, inputEvents)
	q, err := place.Run([]string{"1", "bad time"})
	assert.Error(t, err)
	assert.False(t, q)
}

func TestFinishCommandWithBadBib(t *testing.T) {
	mockInStream := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	place := NewFinishCommand(t.Name(), inputEvents)
	q, err := place.Run([]string{"x"})
	assert.Error(t, err)
	assert.False(t, q)
	assert.Equal(t, 0, len(mockInStream.Events))
}
