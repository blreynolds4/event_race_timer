package places

import (
	"blreynolds4/event-race-timer/events"
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

	testEvents := []events.RaceEvent{
		events.NewStartEvent(t.Name(), now),
		events.NewFinishEvent(t.Name(), finishTime10, 10),
		events.NewFinishEvent(t.Name(), finishTime12, events.NoBib),
		events.NewFinishEvent(t.Name(), finishTime11, 11),
		events.NewFinishEvent(t.Name(), finishTime13, 13),
	}
	inputEvents := NewMockRaceEventSource(testEvents)

	actualResults := &mockEventTarget{
		Events: make([]events.RaceEvent, 0, 5),
	}

	builder := NewPlaceGenerator(inputEvents, actualResults)
	err := builder.GeneratePlaces()
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

func TestEventsArriveOutOfOrder(t *testing.T) {
	// given a set of events on the source
	// produce the set of events on the target
	now := time.Now().UTC()
	// Test data
	finishTime10 := now.Add(5 * time.Minute)
	finishTime11 := now.Add(5*time.Minute + (time.Second * 5))
	finishTime13 := now.Add(5*time.Minute + (time.Second * 29))

	testEvents := []events.RaceEvent{
		events.NewFinishEvent(t.Name(), finishTime13, 13),
		events.NewFinishEvent(t.Name(), finishTime11, 11),
		events.NewFinishEvent(t.Name(), finishTime10, 10),
	}
	inputEvents := NewMockRaceEventSource(testEvents)

	actualResults := &mockEventTarget{
		Events: make([]events.RaceEvent, 0, 6),
	}

	builder := NewPlaceGenerator(inputEvents, actualResults)
	err := builder.GeneratePlaces()
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

func (mes *mockEventSource) GetRaceEvent() (events.RaceEvent, error) {
	if len(mes.events) > 0 {
		var result events.RaceEvent
		// remove the first item in the list and shift everything else up
		result, mes.events = mes.events[0], mes.events[1:]

		return result, nil
	}

	return nil, nil
}

type mockEventTarget struct {
	Events []events.RaceEvent
}

func (mrt *mockEventTarget) SendRaceEvent(e events.RaceEvent) error {
	mrt.Events = append(mrt.Events, e)
	return nil
}
