package results

import (
	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/events"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOverallScoring(t *testing.T) {
	// read events off a stream and return
	// overall results in place order
	now := time.Now().UTC()

	// Test data
	finishTime10 := now.Add(5 * time.Minute)
	finishTime11 := now.Add(5*time.Minute + (time.Second * 5))
	finishTime12 := now.Add(5*time.Minute + (time.Millisecond * 1))
	finishTime13 := now.Add(5*time.Minute + (time.Second * 29))
	finishTime14 := now.Add(5*time.Minute + (time.Second * 30))

	testEvents := []events.RaceEvent{
		events.NewFinishEvent(t.Name(), finishTime10, 10),
		events.NewFinishEvent(t.Name(), finishTime12, 12),
		events.NewFinishEvent(t.Name(), finishTime11, 11),
		events.NewFinishEvent(t.Name(), finishTime14, 14),
		events.NewFinishEvent(t.Name(), finishTime13, 13),
		events.NewStartEvent(t.Name(), now),
		events.NewPlaceEvent(t.Name(), 12, 1),
		events.NewPlaceEvent(t.Name(), 10, 2),
		events.NewPlaceEvent(t.Name(), 11, 3),
		events.NewPlaceEvent(t.Name(), 13, 4),
		events.NewPlaceEvent(t.Name(), 14, 5),
	}
	inputEvents := NewMockRaceEventSource(testEvents)

	athletes := make(competitors.CompetitorLookup)
	athletes[10] = competitors.NewCompetitor("DJR", "WPI", 22, 17)
	athletes[11] = competitors.NewCompetitor("MWR", "OGTC", 22, 16)
	athletes[12] = competitors.NewCompetitor("MGR", "MVXC", 16, 11)
	athletes[13] = competitors.NewCompetitor("SSR", "WPI", 19, 14)
	athletes[14] = competitors.NewCompetitor("SSL", "CU", 53, 20)

	expectedResults := []ScoredResult{
		{
			Athlete: athletes[12],
			Place:   1,
			Time:    finishTime12.Sub(now),
		},
		{
			Athlete: athletes[10],
			Place:   2,
			Time:    finishTime10.Sub(now),
		},
		{
			Athlete: athletes[11],
			Place:   3,
			Time:    finishTime11.Sub(now),
		},
		{
			Athlete: athletes[13],
			Place:   4,
			Time:    finishTime13.Sub(now),
		},
		{
			Athlete: athletes[14],
			Place:   5,
			Time:    finishTime14.Sub(now),
		},
	}

	overall := NewOverallScoring()
	actualResults, err := overall.Score(inputEvents, athletes)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, actualResults)
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
