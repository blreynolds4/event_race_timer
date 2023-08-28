package results

import (
	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/events"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResultBuilderSimplest(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()
	// Test data
	finishTime10 := now.Add(5 * time.Minute)

	// 3 events minimum to build a result:  start, finish and place
	// if the builder doesn't get all 3 no result for the bib is produced
	testEvents := []events.RaceEvent{
		events.NewFinishEvent(t.Name(), finishTime10, 10),
		events.NewStartEvent(t.Name(), now),
		events.NewPlaceEvent(t.Name(), 10, 1),
	}
	inputEvents := NewMockRaceEventSource(testEvents)

	athletes := make(competitors.CompetitorLookup)
	athletes[10] = competitors.NewCompetitor("DJR", "WPI", 22, 17)

	expectedResults := []RaceResult{
		{
			Bib:     10,
			Athlete: athletes[10],
			Place:   1,
			Time:    finishTime10.Sub(now),
		},
	}

	actualResults := &mockResultTarget{
		Results: make([]RaceResult, 0),
	}

	builder := NewResultBuilder()
	err := builder.BuildResults(inputEvents, athletes, actualResults)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, actualResults.Results)
}

func TestResultBuilderFinishUpdated(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()
	// Test data
	finishTime10 := now.Add(5 * time.Minute)
	finishTime10updated := now.Add(6 * time.Minute)

	// 3 events minimum to build a result:  start, finish and place
	// if the builder doesn't get all 3 no result for the bib is produced
	testEvents := []events.RaceEvent{
		events.NewFinishEvent(t.Name(), finishTime10, 10),
		events.NewStartEvent(t.Name(), now),
		events.NewPlaceEvent(t.Name(), 10, 1),
		events.NewFinishEvent(t.Name(), finishTime10updated, 10),
	}
	inputEvents := NewMockRaceEventSource(testEvents)

	athletes := make(competitors.CompetitorLookup)
	athletes[10] = competitors.NewCompetitor("DJR", "WPI", 22, 17)

	// when the first place event comes in the builder should produce a result
	// when the updated finish time comes, it should produce a new result for the
	// same bib
	expectedResults := []RaceResult{
		{
			Bib:     10,
			Athlete: athletes[10],
			Place:   1,
			Time:    finishTime10.Sub(now),
		},
		{
			Bib:     10,
			Athlete: athletes[10],
			Place:   1,
			Time:    finishTime10updated.Sub(now),
		},
	}

	actualResults := &mockResultTarget{
		Results: make([]RaceResult, 0),
	}

	builder := NewResultBuilder()
	err := builder.BuildResults(inputEvents, athletes, actualResults)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, actualResults.Results)
}

func TestResultBuilderStartUpdated(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()
	startUpdated := now.Add(time.Second)
	// Test data
	finishTime10 := now.Add(5 * time.Minute)

	// 3 events minimum to build a result:  start, finish and place
	// if the builder doesn't get all 3 no result for the bib is produced
	testEvents := []events.RaceEvent{
		events.NewFinishEvent(t.Name(), finishTime10, 10),
		events.NewStartEvent(t.Name(), now),
		events.NewPlaceEvent(t.Name(), 10, 1),
		events.NewStartEvent(t.Name(), startUpdated),
	}
	inputEvents := NewMockRaceEventSource(testEvents)

	athletes := make(competitors.CompetitorLookup)
	athletes[10] = competitors.NewCompetitor("DJR", "WPI", 22, 17)

	// when the first place event comes in the builder should produce a result
	// when the updated start time comes, it should produce a new result for the
	// same bib
	expectedResults := []RaceResult{
		{
			Bib:     10,
			Athlete: athletes[10],
			Place:   1,
			Time:    finishTime10.Sub(now),
		},
		{
			Bib:     10,
			Athlete: athletes[10],
			Place:   1,
			Time:    finishTime10.Sub(startUpdated),
		},
	}

	actualResults := &mockResultTarget{
		Results: make([]RaceResult, 0),
	}

	builder := NewResultBuilder()
	err := builder.BuildResults(inputEvents, athletes, actualResults)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, actualResults.Results)
}

func TestResultBuilderPlaceUpdated(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()
	// Test data
	finishTime10 := now.Add(5 * time.Minute)

	// 3 events minimum to build a result:  start, finish and place
	// if the builder doesn't get all 3 no result for the bib is produced
	testEvents := []events.RaceEvent{
		events.NewFinishEvent(t.Name(), finishTime10, 10),
		events.NewStartEvent(t.Name(), now),
		events.NewPlaceEvent(t.Name(), 10, 1),
		events.NewPlaceEvent(t.Name(), 10, 2),
	}
	inputEvents := NewMockRaceEventSource(testEvents)

	athletes := make(competitors.CompetitorLookup)
	athletes[10] = competitors.NewCompetitor("DJR", "WPI", 22, 17)

	// when the first place event comes in the builder should produce a result
	// when the updated start time comes, it should produce a new result for the
	// same bib
	expectedResults := []RaceResult{
		{
			Bib:     10,
			Athlete: athletes[10],
			Place:   1,
			Time:    finishTime10.Sub(now),
		},
		{
			Bib:     10,
			Athlete: athletes[10],
			Place:   2,
			Time:    finishTime10.Sub(now),
		},
	}

	actualResults := &mockResultTarget{
		Results: make([]RaceResult, 0),
	}

	builder := NewResultBuilder()
	err := builder.BuildResults(inputEvents, athletes, actualResults)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, actualResults.Results)
}

func TestResultBuilderNoStartNoResult(t *testing.T) {
	// test a missing start event so no finish time can be calculated
	// no results should be sent
	now := time.Now().UTC()
	// Test data
	finishTime10 := now.Add(5 * time.Minute)

	testEvents := []events.RaceEvent{
		events.NewFinishEvent(t.Name(), finishTime10, 10),
		events.NewPlaceEvent(t.Name(), 10, 1),
	}
	inputEvents := NewMockRaceEventSource(testEvents)

	athletes := make(competitors.CompetitorLookup)
	athletes[10] = competitors.NewCompetitor("DJR", "WPI", 22, 17)

	expectedResults := []RaceResult{}

	actualResults := &mockResultTarget{
		Results: make([]RaceResult, 0),
	}

	builder := NewResultBuilder()
	err := builder.BuildResults(inputEvents, athletes, actualResults)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, actualResults.Results)
}

func TestResultBuilderNoFinishNoResult(t *testing.T) {
	// test not getting a finish event
	// no result should be produced
	now := time.Now().UTC()

	testEvents := []events.RaceEvent{
		events.NewStartEvent(t.Name(), now),
		events.NewPlaceEvent(t.Name(), 10, 1),
	}
	inputEvents := NewMockRaceEventSource(testEvents)

	athletes := make(competitors.CompetitorLookup)
	athletes[10] = competitors.NewCompetitor("DJR", "WPI", 22, 17)

	expectedResults := []RaceResult{}

	actualResults := &mockResultTarget{
		Results: make([]RaceResult, 0),
	}

	builder := NewResultBuilder()
	err := builder.BuildResults(inputEvents, athletes, actualResults)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, actualResults.Results)
}

func TestResultBuilderNoPlaceNoResult(t *testing.T) {
	// test a missing place event
	// no result should be produced
	now := time.Now().UTC()
	// Test data
	finishTime10 := now.Add(5 * time.Minute)

	testEvents := []events.RaceEvent{
		events.NewFinishEvent(t.Name(), finishTime10, 10),
		events.NewStartEvent(t.Name(), now),
	}
	inputEvents := NewMockRaceEventSource(testEvents)

	athletes := make(competitors.CompetitorLookup)
	athletes[10] = competitors.NewCompetitor("DJR", "WPI", 22, 17)

	expectedResults := []RaceResult{}

	actualResults := &mockResultTarget{
		Results: make([]RaceResult, 0),
	}

	builder := NewResultBuilder()
	err := builder.BuildResults(inputEvents, athletes, actualResults)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, actualResults.Results)
}

func TestResultBuilder(t *testing.T) {
	// read events off a stream and and
	// produce a result event as enough info becomes available for each bib
	now := time.Now().UTC()
	// Test data
	finishTime10 := now.Add(5 * time.Minute)
	finishTime11 := now.Add(5*time.Minute + (time.Second * 5))
	finishTime12 := now.Add(5*time.Minute + (time.Millisecond * 1))
	finishTime13 := now.Add(5*time.Minute + (time.Second * 29))
	finishTime14 := now.Add(5*time.Minute + (time.Second * 30))

	testEvents := []events.RaceEvent{
		events.NewStartEvent(t.Name(), now),
		events.NewFinishEvent(t.Name(), finishTime10, 10),
		events.NewFinishEvent(t.Name(), finishTime12, 12),
		events.NewFinishEvent(t.Name(), finishTime11, 11),
		events.NewFinishEvent(t.Name(), finishTime14, 14),
		events.NewFinishEvent(t.Name(), finishTime13, 13),
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

	expectedResults := []RaceResult{
		{
			Bib:     12,
			Athlete: athletes[12],
			Place:   1,
			Time:    finishTime12.Sub(now),
		},
		{
			Bib:     10,
			Athlete: athletes[10],
			Place:   2,
			Time:    finishTime10.Sub(now),
		},
		{
			Bib:     11,
			Athlete: athletes[11],
			Place:   3,
			Time:    finishTime11.Sub(now),
		},
		{
			Bib:     13,
			Athlete: athletes[13],
			Place:   4,
			Time:    finishTime13.Sub(now),
		},
		{
			Bib:     14,
			Athlete: athletes[14],
			Place:   5,
			Time:    finishTime14.Sub(now),
		},
	}

	actualResults := &mockResultTarget{
		Results: make([]RaceResult, 0),
	}

	builder := NewResultBuilder()
	err := builder.BuildResults(inputEvents, athletes, actualResults)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, actualResults.Results)
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

type mockResultTarget struct {
	Results []RaceResult
}

func (mrt *mockResultTarget) SendResult(rr RaceResult) error {
	mrt.Results = append(mrt.Results, rr)
	return nil
}