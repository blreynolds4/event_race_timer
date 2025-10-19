package command

import (
	"blreynolds4/event-race-timer/internal/raceevents"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestListFinishes(t *testing.T) {
	inputEvents := &raceevents.MockEventStream{
		Events: make([]raceevents.Event, 0),
	}

	inputEvents.Events = append(inputEvents.Events, raceevents.Event{
		EventTime: time.Now().UTC(),
		Data: raceevents.StartEvent{
			StartTime: time.Now().UTC(),
		},
	})
	inputEvents.Events = append(inputEvents.Events, raceevents.Event{
		EventTime: time.Now().UTC(),
		Data: raceevents.FinishEvent{
			FinishTime: time.Now().UTC(),
			Bib:        raceevents.NoBib,
		},
	})
	list := NewListFinishCommand(inputEvents)

	// no seed duration arugment
	q, err := list.Run([]string{})
	assert.NoError(t, err)
	assert.False(t, q)
}

func TestListFinishesFailFirstGet(t *testing.T) {
	expErr := fmt.Errorf("fail")
	inputEvents := &raceevents.MockEventStream{
		Get: func(ctx context.Context, timeout time.Duration, msg *raceevents.Event) (bool, error) {
			return false, expErr
		},
		Events: make([]raceevents.Event, 0),
	}

	list := NewListFinishCommand(inputEvents)
	// no seed duration arugment
	q, err := list.Run([]string{})
	assert.Equal(t, expErr, err)
	assert.False(t, q)
}

func TestListFinishesFailSecondGet(t *testing.T) {
	raceMessages := []raceevents.Event{
		{
			EventTime: time.Now().UTC(),
			Data: raceevents.StartEvent{
				Source:    t.Name(),
				StartTime: time.Now().UTC(),
			},
		},
	}
	expErr := fmt.Errorf("fail")
	inputEvents := &raceevents.MockEventStream{
		Get: func(ctx context.Context, timeout time.Duration, msg *raceevents.Event) (bool, error) {
			if len(raceMessages) > 0 {
				*msg = raceMessages[0]
				raceMessages = raceMessages[1:]
				return true, nil
			}
			return false, expErr
		},
		Events: make([]raceevents.Event, 0),
	}

	list := NewListFinishCommand(inputEvents)
	// no seed duration arugment
	q, err := list.Run([]string{})
	assert.Equal(t, expErr, err)
	assert.False(t, q)
}
