package command

import (
	"blreynolds4/event-race-timer/raceevents"
	"blreynolds4/event-race-timer/stream"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestListFinishes(t *testing.T) {
	mockInStream := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	// seed a start and finish event

	mockInStream.Events = append(mockInStream.Events, toMsg(raceevents.Event{
		EventTime: time.Now().UTC(),
		Data: raceevents.StartEvent{
			StartTime: time.Now().UTC(),
		},
	}))
	mockInStream.Events = append(mockInStream.Events, toMsg(raceevents.Event{
		EventTime: time.Now().UTC(),
		Data: raceevents.FinishEvent{
			FinishTime: time.Now().UTC(),
			Bib:        raceevents.NoBib,
		},
	}))
	list := NewListFinishCommand(inputEvents)

	// no seed duration arugment
	q, err := list.Run([]string{})
	assert.NoError(t, err)
	assert.False(t, q)
}

func TestListFinishesFailFirstGet(t *testing.T) {
	expErr := fmt.Errorf("fail")
	mockInStream := &stream.MockStream{
		Get: func(ctx context.Context, timeout time.Duration, msg *stream.Message) (bool, error) {
			fmt.Println("returnning error")
			return false, expErr
		},
		Events: make([]stream.Message, 0),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	list := NewListFinishCommand(inputEvents)
	// no seed duration arugment
	q, err := list.Run([]string{})
	assert.Equal(t, expErr, err)
	assert.False(t, q)
}

func TestListFinishesFailSecondGet(t *testing.T) {
	raceMessages := buildEventMessages(
		[]raceevents.Event{
			{
				EventTime: time.Now().UTC(),
				Data: raceevents.StartEvent{
					Source:    t.Name(),
					StartTime: time.Now().UTC(),
				},
			},
		})
	expErr := fmt.Errorf("fail")
	mockInStream := &stream.MockStream{
		Get: func(ctx context.Context, timeout time.Duration, msg *stream.Message) (bool, error) {
			if len(raceMessages) > 0 {
				*msg = raceMessages[0]
				raceMessages = raceMessages[1:]
				return true, nil
			}
			return false, expErr
		},
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	list := NewListFinishCommand(inputEvents)
	// no seed duration arugment
	q, err := list.Run([]string{})
	assert.Equal(t, expErr, err)
	assert.False(t, q)
}
