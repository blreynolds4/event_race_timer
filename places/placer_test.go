package places

import (
	"blreynolds4/event-race-timer/events"
	"blreynolds4/event-race-timer/eventstream"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNormalPlacingInOrderSkipNoBib(t *testing.T) {
	// given a set of events on the source
	// produce the set of events on the target
	now := time.Now().UTC()
	// Test data
	finishTime10 := now.Add(5 * time.Minute)
	finishTime12 := now.Add(5*time.Minute + (time.Millisecond * 1))
	finishTime11 := now.Add(5*time.Minute + (time.Second * 5))
	finishTime13 := now.Add(5*time.Minute + (time.Second * 29))

	sourceRanks := make(map[string]int)
	bestSource := t.Name()
	slowSource := t.Name() + "-slow"
	sourceRanks[bestSource] = 1
	sourceRanks[slowSource] = 2

	testEvents := []events.RaceEvent{
		eventstream.NewStartEvent(bestSource, now),
		eventstream.NewFinishEvent(bestSource, finishTime10, 10),
		eventstream.NewFinishEvent(bestSource, finishTime12, events.NoBib),
		eventstream.NewFinishEvent(bestSource, finishTime11, 11),
		eventstream.NewFinishEvent(bestSource, finishTime13, 13),
	}
	inputEvents := NewMockRaceEventSource(testEvents)

	actualResults := &mockEventTarget{
		Events: make([]events.RaceEvent, 0, 5),
	}

	builder := NewPlaceGenerator(inputEvents, actualResults)
	err := builder.GeneratePlaces(sourceRanks)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(actualResults.Events))

	// verify the bibs and places match what we expect
	pe, ok := actualResults.Events[0].(events.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 10, pe.GetBib())
	assert.Equal(t, 1, pe.GetPlace())

	pe, ok = actualResults.Events[1].(events.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 11, pe.GetBib())
	assert.Equal(t, 2, pe.GetPlace())

	pe, ok = actualResults.Events[2].(events.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 13, pe.GetBib())
	assert.Equal(t, 3, pe.GetPlace())
}

func TestNoisyMultiSourceEvents(t *testing.T) {
	// given a set of events on the source
	// produce the set of events on the target
	now := time.Now().UTC()
	// Test data
	finishTime10 := now.Add(5 * time.Minute)
	finishTime12 := now.Add(5*time.Minute + (time.Millisecond * 1))
	finishTime11 := now.Add(5*time.Minute + (time.Second * 5))
	finishTime13 := now.Add(5*time.Minute + (time.Second * 29))

	sourceRanks := make(map[string]int)
	bestSource := t.Name()
	slowSource := t.Name() + "-slow"
	sourceRanks[bestSource] = 1
	sourceRanks[slowSource] = 2

	// multiple events per finish, should only use the best time
	// test getting best time first or second
	testEvents := []events.RaceEvent{
		eventstream.NewStartEvent(bestSource, now),

		eventstream.NewFinishEvent(slowSource, finishTime10.Add(time.Second), 10),
		eventstream.NewFinishEvent(bestSource, finishTime10, 10),

		eventstream.NewFinishEvent(bestSource, finishTime12, events.NoBib),

		eventstream.NewFinishEvent(bestSource, finishTime11, 11),
		eventstream.NewFinishEvent(slowSource, finishTime11.Add(time.Minute), 11),

		eventstream.NewFinishEvent(bestSource, finishTime13, 13),
		eventstream.NewFinishEvent(slowSource, finishTime13.Add(-time.Minute), 13),
	}
	inputEvents := NewMockRaceEventSource(testEvents)

	actualResults := &mockEventTarget{
		Events: make([]events.RaceEvent, 0, 4),
	}

	builder := NewPlaceGenerator(inputEvents, actualResults)
	err := builder.GeneratePlaces(sourceRanks)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(actualResults.Events))

	// verify the bibs and places match what we expect
	pe, ok := actualResults.Events[0].(events.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 10, pe.GetBib())
	assert.Equal(t, 1, pe.GetPlace())

	pe, ok = actualResults.Events[1].(events.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 10, pe.GetBib())
	assert.Equal(t, 1, pe.GetPlace())

	pe, ok = actualResults.Events[2].(events.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 11, pe.GetBib())
	assert.Equal(t, 2, pe.GetPlace())

	pe, ok = actualResults.Events[3].(events.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 13, pe.GetBib())
	assert.Equal(t, 3, pe.GetPlace())
}

func TestEventsArriveOutOfOrder(t *testing.T) {
	// given a set of events on the source
	// produce the set of events on the target
	now := time.Now().UTC()
	// Test data
	finishTime10 := now.Add(5 * time.Minute)
	finishTime11 := now.Add(5*time.Minute + (time.Second * 5))
	finishTime13 := now.Add(5*time.Minute + (time.Second * 29))

	sourceRanks := make(map[string]int)

	testEvents := []events.RaceEvent{
		eventstream.NewFinishEvent(t.Name(), finishTime13, 13),
		eventstream.NewFinishEvent(t.Name(), finishTime11, 11),
		eventstream.NewFinishEvent(t.Name(), finishTime10, 10),
	}
	inputEvents := NewMockRaceEventSource(testEvents)

	actualResults := &mockEventTarget{
		Events: make([]events.RaceEvent, 0, 6),
	}

	builder := NewPlaceGenerator(inputEvents, actualResults)
	err := builder.GeneratePlaces(sourceRanks)
	assert.NoError(t, err)
	assert.Equal(t, 6, len(actualResults.Events))

	// verify the bibs and places match what we expect
	pe, ok := actualResults.Events[0].(events.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 13, pe.GetBib())
	assert.Equal(t, 1, pe.GetPlace())

	pe, ok = actualResults.Events[1].(events.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 11, pe.GetBib())
	assert.Equal(t, 1, pe.GetPlace())

	pe, ok = actualResults.Events[2].(events.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 13, pe.GetBib())
	assert.Equal(t, 2, pe.GetPlace())

	pe, ok = actualResults.Events[3].(events.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 10, pe.GetBib())
	assert.Equal(t, 1, pe.GetPlace())

	pe, ok = actualResults.Events[4].(events.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 11, pe.GetBib())
	assert.Equal(t, 2, pe.GetPlace())

	pe, ok = actualResults.Events[5].(events.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 13, pe.GetBib())
	assert.Equal(t, 3, pe.GetPlace())
}

func NewMockRaceEventSource(testEvents []events.RaceEvent) events.EventSource {
	return &mockEventSource{events: testEvents}
}

type mockEventSource struct {
	events []events.RaceEvent
}

func (mes *mockEventSource) GetRaceEvent(ctx context.Context, timeout time.Duration) (events.RaceEvent, error) {
	if len(mes.events) > 0 {
		var result events.RaceEvent
		// remove the first item in the list and shift everything else up
		result, mes.events = mes.events[0], mes.events[1:]

		return result, nil
	}

	return nil, nil
}

func (mes *mockEventSource) GetRaceEventRange(ctx context.Context, start, end string) ([]events.RaceEvent, error) {
	return []events.RaceEvent{}, nil
}

type mockEventTarget struct {
	Events []events.RaceEvent
}

func (mrt *mockEventTarget) SendRaceEvent(ctx context.Context, e events.RaceEvent) error {
	mrt.Events = append(mrt.Events, e)
	return nil
}
