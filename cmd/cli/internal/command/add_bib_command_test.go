package command

import (
	"blreynolds4/event-race-timer/internal/raceevents"
	"blreynolds4/event-race-timer/internal/stream"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddBibMissingArgs(t *testing.T) {
	mockInStream := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	list := NewAddBibCommand(inputEvents)
	// missing
	q, err := list.Run([]string{})
	assert.Error(t, err)
	assert.False(t, q)
}

func TestAddBibMissingBadBib(t *testing.T) {
	mockInStream := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	list := NewAddBibCommand(inputEvents)
	q, err := list.Run([]string{"x", "y"})
	assert.Error(t, err)
	assert.False(t, q)
}

func TestAddBibMissingBibRangeFails(t *testing.T) {
	expErr := fmt.Errorf("fail")
	mockInStream := &stream.MockStream{
		Range: func(ctx context.Context, startId, endId string, msgs []stream.Message) (int, error) {
			return 0, expErr
		},
		Events: make([]stream.Message, 0),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	list := NewAddBibCommand(inputEvents)
	q, err := list.Run([]string{"msgid", "1"})
	assert.Equal(t, expErr, err)
	assert.False(t, q)
}

func TestAddBibMissingEvent(t *testing.T) {
	expErr := fmt.Errorf("event id did not return 1 event")
	mockInStream := &stream.MockStream{
		Range: func(ctx context.Context, startId, endId string, msgs []stream.Message) (int, error) {
			return 0, expErr
		},
		Events: make([]stream.Message, 0),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	list := NewAddBibCommand(inputEvents)
	q, err := list.Run([]string{"msgid", "1"})
	assert.Equal(t, expErr, err)
	assert.False(t, q)
}

func TestAddBibWrongEventType(t *testing.T) {
	expErr := fmt.Errorf("expected event id to be for finish event, skipping")
	mockInStream := &stream.MockStream{
		Range: func(ctx context.Context, startId, endId string, msgs []stream.Message) (int, error) {
			return 0, expErr
		},
		Events: make([]stream.Message, 0),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	list := NewAddBibCommand(inputEvents)
	q, err := list.Run([]string{"msgid", "1"})
	assert.Equal(t, expErr, err)
	assert.False(t, q)
}

func TestAddBib(t *testing.T) {
	mockInStream := &stream.MockStream{
		Events: []stream.Message{toMsg(raceevents.Event{
			Data: raceevents.FinishEvent{
				Source:     t.Name(),
				FinishTime: time.Now().UTC(),
				Bib:        raceevents.NoBib,
			},
		})},
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	list := NewAddBibCommand(inputEvents)
	q, err := list.Run([]string{"msgid", "1"})
	assert.NoError(t, err)
	assert.False(t, q)
	assert.Equal(t, 1, len(mockInStream.Events))
	actualEvents := buildActualResults(mockInStream)

	fe := actualEvents[0].Data.(raceevents.FinishEvent)
	assert.Equal(t, 1, fe.Bib)
}
