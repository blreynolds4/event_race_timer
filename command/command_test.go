package command

import (
	"blreynolds4/event-race-timer/events"
	"blreynolds4/event-race-timer/eventstream"
	"context"
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
	mockTarget := &events.MockRaceEventStream{
		Events: make([]events.RaceEvent, 0),
	}

	start := NewStartCommand(mockTarget)
	// no seed duration arugment
	q, err := start.Run([]string{})
	assert.NoError(t, err)
	assert.False(t, q)
	assert.Equal(t, 1, len(mockTarget.Events))

	assert.Equal(t, events.StartEventType, mockTarget.Events[0].GetType())
	se, ok := mockTarget.Events[0].(events.StartEvent)
	assert.True(t, ok)
	startTime := se.GetStartTime()
	assert.False(t, startTime.IsZero())
}

func TestStartCommandWithTime(t *testing.T) {
	mockTarget := &events.MockRaceEventStream{
		Events: make([]events.RaceEvent, 0),
	}

	now := time.Now().UTC()

	start := NewStartCommand(mockTarget)
	// with duration argument
	q, err := start.Run([]string{time.Minute.String()})
	assert.NoError(t, err)
	assert.False(t, q)
	assert.Equal(t, 1, len(mockTarget.Events))

	se, ok := mockTarget.Events[0].(events.StartEvent)
	assert.True(t, ok)
	startTime := se.GetStartTime()
	assert.True(t, startTime.Before(now))
}

func TestStartCommandWithBadTime(t *testing.T) {
	mockTarget := &events.MockRaceEventStream{
		Events: make([]events.RaceEvent, 0),
	}

	start := NewStartCommand(mockTarget)
	// with duration argument
	q, err := start.Run([]string{"bad"})
	assert.Error(t, err)
	assert.False(t, q)
	assert.Equal(t, 0, len(mockTarget.Events))
}

func TestFinishCommandNoBib(t *testing.T) {
	mockTarget := &events.MockRaceEventStream{
		Events: make([]events.RaceEvent, 0),
	}

	place := NewFinishCommand(mockTarget)
	q, err := place.Run([]string{})
	assert.NoError(t, err)
	assert.False(t, q)
	assert.Equal(t, 1, len(mockTarget.Events))

	pe, ok := mockTarget.Events[0].(events.FinishEvent)
	assert.True(t, ok)
	assert.Equal(t, events.NoBib, pe.GetBib())
}

func TestFinishCommandWithBib(t *testing.T) {
	mockTarget := &events.MockRaceEventStream{
		Events: make([]events.RaceEvent, 0),
	}

	place := NewFinishCommand(mockTarget)
	q, err := place.Run([]string{"1"})
	assert.NoError(t, err)
	assert.False(t, q)
	assert.Equal(t, 1, len(mockTarget.Events))

	pe, ok := mockTarget.Events[0].(events.FinishEvent)
	assert.True(t, ok)
	assert.Equal(t, 1, pe.GetBib())
}

func TestFinishCommandWithBadBib(t *testing.T) {
	mockTarget := &events.MockRaceEventStream{
		Events: make([]events.RaceEvent, 0),
	}

	place := NewFinishCommand(mockTarget)
	q, err := place.Run([]string{"x"})
	assert.Error(t, err)
	assert.False(t, q)
	assert.Equal(t, 0, len(mockTarget.Events))
}

func TestPlacetCommandNoBib(t *testing.T) {
	mockTarget := &events.MockRaceEventStream{
		Events: make([]events.RaceEvent, 0),
	}

	place := NewPlaceCommand(mockTarget)
	q, err := place.Run([]string{"1", "1"})
	assert.NoError(t, err)
	assert.False(t, q)
	assert.Equal(t, 1, len(mockTarget.Events))

	pe, ok := mockTarget.Events[0].(events.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 1, pe.GetPlace())
	assert.Equal(t, 1, pe.GetBib())
}

func TestPlacetCommandMissingArg(t *testing.T) {
	mockTarget := &events.MockRaceEventStream{
		Events: make([]events.RaceEvent, 0),
	}

	place := NewPlaceCommand(mockTarget)
	q, err := place.Run([]string{"1"})
	assert.Error(t, err)
	assert.False(t, q)
}

func TestPlacetCommandBadBib(t *testing.T) {
	mockTarget := &events.MockRaceEventStream{
		Events: make([]events.RaceEvent, 0),
	}

	place := NewPlaceCommand(mockTarget)
	q, err := place.Run([]string{"x", "1"})
	assert.Error(t, err)
	assert.False(t, q)
}

func TestPlacetCommandBadPlace(t *testing.T) {
	mockTarget := &events.MockRaceEventStream{
		Events: make([]events.RaceEvent, 0),
	}

	place := NewPlaceCommand(mockTarget)
	q, err := place.Run([]string{"1", "x"})
	assert.Error(t, err)
	assert.False(t, q)
}

func TestStartListFinishes(t *testing.T) {
	mockSource := &events.MockRaceEventStream{
		Events: make([]events.RaceEvent, 0),
	}

	// seed a start and finish event
	mockSource.Events = append(mockSource.Events, eventstream.NewStartEvent(t.Name(), time.Now().UTC()))
	mockSource.Events = append(mockSource.Events, eventstream.NewFinishEvent(t.Name(), time.Now().UTC(), events.NoBib))

	list := NewListFinishCommand(mockSource)
	// no seed duration arugment
	q, err := list.Run([]string{})
	assert.NoError(t, err)
	assert.False(t, q)
}

func TestStartListFinishesFailFirstGet(t *testing.T) {
	expErr := fmt.Errorf("fail")
	mockSource := &events.MockRaceEventStream{
		Get: func(ctx context.Context, t time.Duration) (events.RaceEvent, error) {
			fmt.Println("returnning error")
			return nil, expErr
		},
		Events: make([]events.RaceEvent, 0),
	}

	list := NewListFinishCommand(mockSource)
	// no seed duration arugment
	q, err := list.Run([]string{})
	assert.Equal(t, expErr, err)
	assert.False(t, q)
}

func TestStartListFinishesFailSecondGet(t *testing.T) {
	raceEvents := []events.RaceEvent{
		eventstream.NewStartEvent(t.Name(), time.Now().UTC()),
	}

	expErr := fmt.Errorf("fail")
	mockSource := &events.MockRaceEventStream{
		Get: func(ctx context.Context, t time.Duration) (events.RaceEvent, error) {
			if len(raceEvents) > 0 {
				result := raceEvents[0]
				raceEvents = raceEvents[1:]
				return result, nil
			}
			return nil, expErr
		},
	}

	list := NewListFinishCommand(mockSource)
	// no seed duration arugment
	q, err := list.Run([]string{})
	assert.Equal(t, expErr, err)
	assert.False(t, q)
}

func TestAddBibMissingArgs(t *testing.T) {
	mockEvents := &events.MockRaceEventStream{
		Events: make([]events.RaceEvent, 0),
	}

	list := NewAddBibCommand(mockEvents, mockEvents)
	// missing
	q, err := list.Run([]string{})
	assert.Error(t, err)
	assert.False(t, q)
}

func TestAddBibMissingBadBib(t *testing.T) {
	mockEvents := &events.MockRaceEventStream{
		Events: make([]events.RaceEvent, 0),
	}

	list := NewAddBibCommand(mockEvents, mockEvents)
	q, err := list.Run([]string{"x", "y"})
	assert.Error(t, err)
	assert.False(t, q)
}

func TestAddBibMissingBibRangeFails(t *testing.T) {
	expErr := fmt.Errorf("fail")
	mockEvents := &events.MockRaceEventStream{
		Range: func(ctx context.Context, start, end string) ([]events.RaceEvent, error) {
			return []events.RaceEvent{}, expErr
		},
	}

	list := NewAddBibCommand(mockEvents, mockEvents)
	q, err := list.Run([]string{"msgid", "1"})
	assert.Equal(t, expErr, err)
	assert.False(t, q)
}

func TestAddBibMissingEvent(t *testing.T) {
	expErr := fmt.Errorf("event id did not return 1 event")
	mockEvents := &events.MockRaceEventStream{
		Events: make([]events.RaceEvent, 0),
	}

	list := NewAddBibCommand(mockEvents, mockEvents)
	q, err := list.Run([]string{"msgid", "1"})
	assert.Equal(t, expErr, err)
	assert.False(t, q)
}

func TestAddBibWrongEventType(t *testing.T) {
	expErr := fmt.Errorf("expected event id to be for finish event, skipping")
	mockEvents := &events.MockRaceEventStream{
		Events: []events.RaceEvent{eventstream.NewStartEvent(t.Name(), time.Now().UTC())},
	}

	list := NewAddBibCommand(mockEvents, mockEvents)
	q, err := list.Run([]string{"msgid", "1"})
	assert.Equal(t, expErr, err)
	assert.False(t, q)
}

func TestAddBib(t *testing.T) {
	mockEvents := &events.MockRaceEventStream{
		Events: []events.RaceEvent{eventstream.NewFinishEvent(t.Name(), time.Now().UTC(), events.NoBib)},
	}

	list := NewAddBibCommand(mockEvents, mockEvents)
	q, err := list.Run([]string{"msgid", "1"})
	assert.NoError(t, err)
	assert.False(t, q)
	assert.Equal(t, 1, len(mockEvents.Events))
	assert.Equal(t, events.FinishEventType, mockEvents.Events[0].GetType())
	fe := mockEvents.Events[0].(events.FinishEvent)
	assert.Equal(t, 1, fe.GetBib())
}
