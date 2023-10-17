package command

import (
	"blreynolds4/event-race-timer/raceevents"
	"blreynolds4/event-race-timer/stream"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestQuitCommand(t *testing.T) {
	quit := NewQuitCommand()
	q, err := quit.Run([]string{})
	assert.NoError(t, err)
	assert.True(t, q)
}

func TestPingCommand(t *testing.T) {
	db, mock := redismock.NewClientMock()

	// set up expectations
	mock.ExpectPing().SetVal("pong")

	ping := NewPingCommand(db)
	q, err := ping.Run([]string{})
	assert.NoError(t, err)
	assert.False(t, q)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

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

func TestPlacetCommandNoBib(t *testing.T) {
	mockInStream := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	eventSource := t.Name()
	place := NewPlaceCommand(eventSource, inputEvents)
	q, err := place.Run([]string{"1", "1"})
	assert.NoError(t, err)
	assert.False(t, q)
	assert.Equal(t, 1, len(mockInStream.Events))
	actualEvents := buildActualResults(mockInStream)

	pe, ok := actualEvents[0].Data.(raceevents.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 1, pe.Place)
	assert.Equal(t, 1, pe.Bib)
	assert.Equal(t, eventSource, pe.Source)
}

func TestPlacetCommandMissingArg(t *testing.T) {
	mockInStream := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	place := NewPlaceCommand(t.Name(), inputEvents)
	q, err := place.Run([]string{"1"})
	assert.Error(t, err)
	assert.False(t, q)
}

func TestPlacetCommandBadBib(t *testing.T) {
	mockInStream := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	place := NewPlaceCommand(t.Name(), inputEvents)
	q, err := place.Run([]string{"x", "1"})
	assert.Error(t, err)
	assert.False(t, q)
}

func TestPlacetCommandBadPlace(t *testing.T) {
	mockInStream := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	place := NewPlaceCommand(t.Name(), inputEvents)
	q, err := place.Run([]string{"1", "x"})
	assert.Error(t, err)
	assert.False(t, q)
}

func TestStartListFinishes(t *testing.T) {
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

func TestStartListFinishesFailFirstGet(t *testing.T) {
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

func TestStartListFinishesFailSecondGet(t *testing.T) {
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

func buildActualResults(rawOutput *stream.MockStream) []raceevents.Event {
	result := make([]raceevents.Event, len(rawOutput.Events))
	for i, msg := range rawOutput.Events {
		err := json.Unmarshal(msg.Data, &result[i])
		if err != nil {
			panic(err)
		}
	}

	return result
}

func toMsg(r raceevents.Event) stream.Message {
	var msg stream.Message
	msgData, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	msg.Data = msgData

	return msg
}

func buildEventMessages(testEvents []raceevents.Event) []stream.Message {
	result := make([]stream.Message, len(testEvents))
	for i, e := range testEvents {
		eData, err := json.Marshal(e)
		if err != nil {
			panic(err)
		}
		result[i] = stream.Message{
			ID:   e.ID,
			Data: eData,
		}
	}

	return result
}
