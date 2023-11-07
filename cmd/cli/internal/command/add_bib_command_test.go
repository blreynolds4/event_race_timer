package command

import (
	"blreynolds4/event-race-timer/internal/raceevents"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddBibMissingArgs(t *testing.T) {
	inputEvents := &raceevents.MockEventStream{}

	list := NewAddBibCommand(inputEvents)
	// missing
	q, err := list.Run([]string{})
	assert.Error(t, err)
	assert.False(t, q)
}

func TestAddBibMissingBadBib(t *testing.T) {
	inputEvents := &raceevents.MockEventStream{}

	list := NewAddBibCommand(inputEvents)
	q, err := list.Run([]string{"x", "y"})
	assert.Error(t, err)
	assert.False(t, q)
}

func TestAddBibMissingBibRangeFails(t *testing.T) {
	expErr := fmt.Errorf("fail")
	inputEvents := &raceevents.MockEventStream{
		Range: func(ctx context.Context, startId, endId string, msgs []raceevents.Event) (int, error) {
			return 0, expErr
		},
		Events: make([]raceevents.Event, 0),
	}

	ab := NewAddBibCommand(inputEvents)
	q, err := ab.Run([]string{"msgid", "1"})
	assert.Equal(t, expErr, err)
	assert.False(t, q)
}

func TestAddBibMissingEvent(t *testing.T) {
	expErr := fmt.Errorf("event id did not return 1 event")
	inputEvents := &raceevents.MockEventStream{
		Range: func(ctx context.Context, startId, endId string, msgs []raceevents.Event) (int, error) {
			return 0, expErr
		},
		Events: make([]raceevents.Event, 0),
	}

	ab := NewAddBibCommand(inputEvents)
	q, err := ab.Run([]string{"msgid", "1"})
	assert.Equal(t, expErr, err)
	assert.False(t, q)
}

func TestAddBibWrongEventType(t *testing.T) {
	expErr := fmt.Errorf("expected event id to be for finish event, skipping")
	inputEvents := &raceevents.MockEventStream{
		Range: func(ctx context.Context, startId, endId string, msgs []raceevents.Event) (int, error) {
			return 0, expErr
		},
		Events: make([]raceevents.Event, 0),
	}

	ab := NewAddBibCommand(inputEvents)
	q, err := ab.Run([]string{"msgid", "1"})
	assert.Equal(t, expErr, err)
	assert.False(t, q)
}

func TestAddBib(t *testing.T) {
	inputEvents := &raceevents.MockEventStream{
		Events: []raceevents.Event{
			{
				Data: raceevents.FinishEvent{
					Source:     t.Name(),
					FinishTime: time.Now().UTC(),
					Bib:        raceevents.NoBib,
				},
			},
		},
	}

	ab := NewAddBibCommand(inputEvents)
	q, err := ab.Run([]string{"msgid", "1"})
	assert.NoError(t, err)
	assert.False(t, q)
	assert.Equal(t, 1, len(inputEvents.Events))

	fe := inputEvents.Events[0].Data.(raceevents.FinishEvent)
	assert.Equal(t, 1, fe.Bib)
}
