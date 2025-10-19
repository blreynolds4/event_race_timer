package resultbuilder

import (
	"blreynolds4/event-race-timer/internal/meets"
	"blreynolds4/event-race-timer/internal/raceevents"
	"blreynolds4/event-race-timer/internal/results"
	"log/slog"

	"blreynolds4/event-race-timer/internal/stream"
	"encoding/json"
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
	testEvents := []raceevents.Event{
		{
			ID:        "1",
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     t.Name(),
				Bib:        10,
				FinishTime: finishTime10,
			},
		},
		{
			ID:        "2",
			EventTime: now,
			Data: raceevents.StartEvent{
				Source:    t.Name(),
				StartTime: now,
			},
		},
		{
			ID:        "3",
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Bib:    10,
				Place:  1,
			},
		},
	}

	mockInStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	athletes := make(meets.AthleteLookup)
	athletes[10] = meets.NewAthlete("D", "R", "WPI", "DAID", 12, "m")

	expectedResults := []meets.RaceResult{
		{
			Bib:          10,
			Athlete:      athletes[10],
			Place:        1,
			Time:         finishTime10.Sub(now),
			FinishSource: t.Name(),
			PlaceSource:  t.Name(),
		},
	}

	mockOutStream := &stream.MockStream{
		Events: make([]stream.Message, 0, 10),
	}
	actualResults := results.NewResultStream(mockOutStream)

	builder := NewResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking[t.Name()] = 1
	err := builder.BuildResults(inputEvents, athletes, actualResults, ranking)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, buildActualResults(mockOutStream))
}

func TestResultBuilderPlaceCausesResult(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()
	// Test data
	testEvents := []raceevents.Event{
		{
			ID:        "3",
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Bib:    10,
				Place:  1,
			},
		},
	}

	mockInStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	athletes := make(meets.AthleteLookup)
	athletes[10] = meets.NewAthlete("D", "R", "WPI", "DAID", 12, "m")

	expectedResults := []meets.RaceResult{
		{
			Bib:         10,
			Athlete:     athletes[10],
			Place:       1,
			PlaceSource: t.Name(),
		},
	}

	mockOutStream := &stream.MockStream{
		Events: make([]stream.Message, 0, 10),
	}
	actualResults := results.NewResultStream(mockOutStream)

	builder := NewResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking[t.Name()] = 1
	err := builder.BuildResults(inputEvents, athletes, actualResults, ranking)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, buildActualResults(mockOutStream))
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
	testEvents := []raceevents.Event{
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     t.Name(),
				FinishTime: finishTime10,
				Bib:        10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.StartEvent{
				Source:    t.Name(),
				StartTime: now,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Bib:    10,
				Place:  1,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     t.Name(),
				FinishTime: finishTime10updated,
				Bib:        10,
			},
		},
	}
	mockInStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	athletes := make(meets.AthleteLookup)
	athletes[10] = meets.NewAthlete("D", "R", "WPI", "DAID", 12, "m")

	// when the first place event comes in the builder should produce a result
	// when the updated finish time comes, it should produce a new result for the
	// same bib
	expectedResults := []meets.RaceResult{
		{
			Bib:          10,
			Athlete:      athletes[10],
			Place:        1,
			Time:         finishTime10.Sub(now),
			FinishSource: t.Name(),
			PlaceSource:  t.Name(),
		},
		{
			Bib:          10,
			Athlete:      athletes[10],
			Place:        1,
			Time:         finishTime10updated.Sub(now),
			FinishSource: t.Name(),
			PlaceSource:  t.Name(),
		},
	}

	mockOutStream := &stream.MockStream{
		Events: make([]stream.Message, 0, 10),
	}
	actualResults := results.NewResultStream(mockOutStream)

	builder := NewResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking[t.Name()] = 1
	err := builder.BuildResults(inputEvents, athletes, actualResults, ranking)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, buildActualResults(mockOutStream))
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
	testEvents := []raceevents.Event{
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     t.Name(),
				FinishTime: finishTime10,
				Bib:        10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.StartEvent{
				Source:    t.Name(),
				StartTime: now,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  1,
				Bib:    10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.StartEvent{
				Source:    t.Name(),
				StartTime: startUpdated,
			},
		},
	}
	mockInStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	athletes := make(meets.AthleteLookup)
	athletes[10] = meets.NewAthlete("D", "R", "WPI", "DAID", 12, "m")

	// when the first place event comes in the builder should produce a result
	// when the updated start time comes, it should produce a new result for the
	// same bib
	expectedResults := []meets.RaceResult{
		{
			Bib:          10,
			Athlete:      athletes[10],
			Place:        1,
			Time:         finishTime10.Sub(now),
			FinishSource: t.Name(),
			PlaceSource:  t.Name(),
		},
		{
			Bib:          10,
			Athlete:      athletes[10],
			Place:        1,
			Time:         finishTime10.Sub(startUpdated),
			FinishSource: t.Name(),
			PlaceSource:  t.Name(),
		},
	}

	mockOutStream := &stream.MockStream{
		Events: make([]stream.Message, 0, 10),
	}
	actualResults := results.NewResultStream(mockOutStream)

	builder := NewResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking[t.Name()] = 1
	err := builder.BuildResults(inputEvents, athletes, actualResults, ranking)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, buildActualResults(mockOutStream))
}

func TestResultBuilderPlaceUpdated(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()
	// Test data
	finishTime10 := now.Add(5 * time.Minute)
	finishTime11 := now.Add(5 * time.Minute)

	// 3 events minimum to build a result:  start, finish and place
	// if the builder doesn't get all 3 no result for the bib is produced
	testEvents := []raceevents.Event{
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     t.Name(),
				FinishTime: finishTime10,
				Bib:        10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.StartEvent{
				Source:    t.Name(),
				StartTime: now,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     t.Name(),
				FinishTime: finishTime11,
				Bib:        11,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  1,
				Bib:    10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  2,
				Bib:    11,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  1,
				Bib:    11,
			},
		},
	}
	mockInStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	athletes := make(meets.AthleteLookup)
	athletes[10] = meets.NewAthlete("D", "R", "WPI", "DAID", 12, "m")
	athletes[11] = meets.NewAthlete("M", "R", "WPI", "DAID2", 12, "m")

	// when the first place event comes in the builder should produce a result
	// when the updated start time comes, it should produce a new result for the
	// same bib
	expectedResults := []meets.RaceResult{
		{
			Bib:          10,
			Athlete:      athletes[10],
			Place:        1,
			Time:         finishTime10.Sub(now),
			FinishSource: t.Name(),
			PlaceSource:  t.Name(),
		},
		{
			Bib:          11,
			Athlete:      athletes[11],
			Place:        2,
			Time:         finishTime11.Sub(now),
			FinishSource: t.Name(),
			PlaceSource:  t.Name(),
		},
		{
			Bib:          11,
			Athlete:      athletes[11],
			Place:        1,
			Time:         finishTime11.Sub(now),
			FinishSource: t.Name(),
			PlaceSource:  t.Name(),
		},
		{
			Bib:          10,
			Athlete:      athletes[10],
			Place:        2,
			Time:         finishTime10.Sub(now),
			FinishSource: t.Name(),
			PlaceSource:  t.Name(),
		},
	}

	mockOutStream := &stream.MockStream{
		Events: make([]stream.Message, 0, 10),
	}
	actualResults := results.NewResultStream(mockOutStream)

	builder := NewResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking[t.Name()] = 1
	err := builder.BuildResults(inputEvents, athletes, actualResults, ranking)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, buildActualResults(mockOutStream))
}

func TestResultBuilderPlaceSwap(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()
	testEvents := []raceevents.Event{
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  1,
				Bib:    10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  2,
				Bib:    11,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  3,
				Bib:    12,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  1,
				Bib:    11,
			},
		},
	}
	mockInStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	athletes := make(meets.AthleteLookup)
	athletes[10] = meets.NewAthlete("D", "R", "WPI", "DAID", 12, "m")
	athletes[11] = meets.NewAthlete("M", "R", "WPI", "DAID2", 12, "m")
	athletes[12] = meets.NewAthlete("MG", "R", "MVHS", "DAID3", 12, "m")

	expectedResults := []meets.RaceResult{
		{
			Bib:         10,
			Athlete:     athletes[10],
			Place:       1,
			PlaceSource: t.Name(),
		},
		{
			Bib:         11,
			Athlete:     athletes[11],
			Place:       2,
			PlaceSource: t.Name(),
		},
		{
			Bib:         12,
			Athlete:     athletes[12],
			Place:       3,
			PlaceSource: t.Name(),
		},
		{
			Bib:         11,
			Athlete:     athletes[11],
			Place:       1,
			PlaceSource: t.Name(),
		},
		{
			Bib:         10,
			Athlete:     athletes[10],
			Place:       2,
			PlaceSource: t.Name(),
		},
	}

	mockOutStream := &stream.MockStream{
		Events: make([]stream.Message, 0, 10),
	}
	actualResults := results.NewResultStream(mockOutStream)

	builder := NewResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking[t.Name()] = 1
	err := builder.BuildResults(inputEvents, athletes, actualResults, ranking)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, buildActualResults(mockOutStream))
}

func TestResultBuilderUnknownPlaceBib(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()
	testEvents := []raceevents.Event{
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  1,
				Bib:    999,
			},
		},
	}
	mockInStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	athletes := make(meets.AthleteLookup)

	expectedResults := []meets.RaceResult{}

	mockOutStream := &stream.MockStream{
		Events: make([]stream.Message, 0, 10),
	}
	actualResults := results.NewResultStream(mockOutStream)

	builder := NewResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking[t.Name()] = 1
	err := builder.BuildResults(inputEvents, athletes, actualResults, ranking)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, buildActualResults(mockOutStream))
}

func TestResultBuilderNoPlaceNoResult(t *testing.T) {
	// test a missing place event
	// no result should be produced
	now := time.Now().UTC()
	// Test data
	finishTime10 := now.Add(5 * time.Minute)

	testEvents := []raceevents.Event{
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     t.Name(),
				FinishTime: finishTime10,
				Bib:        10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.StartEvent{
				Source:    t.Name(),
				StartTime: now,
			},
		},
	}
	mockInStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	athletes := make(meets.AthleteLookup)
	athletes[10] = meets.NewAthlete("D", "R", "WPI", "DAID", 12, "m")

	expectedResults := []meets.RaceResult{}

	mockOutStream := &stream.MockStream{
		Events: make([]stream.Message, 0, 10),
	}
	actualResults := results.NewResultStream(mockOutStream)

	builder := NewResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking[t.Name()] = 1
	err := builder.BuildResults(inputEvents, athletes, actualResults, ranking)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, buildActualResults(mockOutStream))
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

	testEvents := []raceevents.Event{
		{
			EventTime: now,
			Data: raceevents.StartEvent{
				Source:    t.Name(),
				StartTime: now,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     t.Name(),
				FinishTime: finishTime10,
				Bib:        10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     t.Name(),
				FinishTime: finishTime12,
				Bib:        12,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     t.Name(),
				FinishTime: finishTime11,
				Bib:        11,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     t.Name(),
				FinishTime: finishTime14,
				Bib:        14,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     t.Name(),
				FinishTime: finishTime13,
				Bib:        13,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  1,
				Bib:    12,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  2,
				Bib:    10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  3,
				Bib:    11,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  4,
				Bib:    13,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  5,
				Bib:    14,
			},
		},
	}
	mockInStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	athletes := make(meets.AthleteLookup)
	athletes[10] = meets.NewAthlete("D", "R", "WPI", "DAID", 12, "m")
	athletes[11] = meets.NewAthlete("M", "R", "WPI", "DAID2", 12, "m")
	athletes[12] = meets.NewAthlete("MG", "R", "MVHS", "DAID3", 12, "m")
	athletes[13] = meets.NewAthlete("S", "R", "WPI", "DAID4", 12, "f")
	athletes[14] = meets.NewAthlete("S", "L", "CU", "DAID5", 12, "f")

	expectedResults := []meets.RaceResult{
		{
			Bib:          12,
			Athlete:      athletes[12],
			Place:        1,
			Time:         finishTime12.Sub(now),
			FinishSource: t.Name(),
			PlaceSource:  t.Name(),
		},
		{
			Bib:          10,
			Athlete:      athletes[10],
			Place:        2,
			Time:         finishTime10.Sub(now),
			FinishSource: t.Name(),
			PlaceSource:  t.Name(),
		},
		{
			Bib:          11,
			Athlete:      athletes[11],
			Place:        3,
			Time:         finishTime11.Sub(now),
			FinishSource: t.Name(),
			PlaceSource:  t.Name(),
		},
		{
			Bib:          13,
			Athlete:      athletes[13],
			Place:        4,
			Time:         finishTime13.Sub(now),
			FinishSource: t.Name(),
			PlaceSource:  t.Name(),
		},
		{
			Bib:          14,
			Athlete:      athletes[14],
			Place:        5,
			Time:         finishTime14.Sub(now),
			FinishSource: t.Name(),
			PlaceSource:  t.Name(),
		},
	}

	mockOutStream := &stream.MockStream{
		Events: make([]stream.Message, 0, 10),
	}
	actualResults := results.NewResultStream(mockOutStream)

	builder := NewResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking[t.Name()] = 1
	err := builder.BuildResults(inputEvents, athletes, actualResults, ranking)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, buildActualResults(mockOutStream))
}

func TestResultBuilderRankUpdates(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()
	// Test data
	finishTime10 := now.Add(6 * time.Minute)
	finishTime10better := now.Add(5 * time.Minute)

	// 3 events minimum to build a result:  start, finish and place
	// if the builder doesn't get all 3 no result for the bib is produced
	testEvents := []raceevents.Event{
		{
			EventTime: now,
			Data: raceevents.StartEvent{
				Source:    t.Name(),
				StartTime: now,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  1,
				Bib:    10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     "worse",
				FinishTime: finishTime10,
				Bib:        10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     "better",
				FinishTime: finishTime10better,
				Bib:        10,
			},
		},
	}
	mockInStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	athletes := make(meets.AthleteLookup)
	athletes[10] = meets.NewAthlete("D", "R", "WPI", "DAID", 12, "m")

	// when the first place event comes in the builder should produce a result
	// when the updated start time comes, it should produce a new result for the
	// same bib
	expectedResults := []meets.RaceResult{
		{
			Bib:          10,
			Athlete:      athletes[10],
			Place:        1,
			PlaceSource:  t.Name(),
			Time:         time.Duration(0),
			FinishSource: "",
		},
		{
			Bib:          10,
			Athlete:      athletes[10],
			Place:        1,
			Time:         finishTime10.Sub(now),
			FinishSource: "worse",
			PlaceSource:  t.Name(),
		},
		{
			Bib:          10,
			Athlete:      athletes[10],
			Place:        1,
			Time:         finishTime10better.Sub(now),
			FinishSource: "better",
			PlaceSource:  t.Name(),
		},
	}

	mockOutStream := &stream.MockStream{
		Events: make([]stream.Message, 0, 10),
	}
	actualResults := results.NewResultStream(mockOutStream)

	builder := NewResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking["better"] = 1
	ranking["worse"] = 2
	err := builder.BuildResults(inputEvents, athletes, actualResults, ranking)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, buildActualResults(mockOutStream))
}

func TestResultBuilderRankIgnores(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()
	// Test data
	finishTime10 := now.Add(5 * time.Minute)

	// 3 events minimum to build a result:  start, finish and place
	// if the builder doesn't get all 3 no result for the bib is produced
	testEvents := []raceevents.Event{
		{
			EventTime: now,
			Data: raceevents.StartEvent{
				Source:    t.Name(),
				StartTime: now,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  1,
				Bib:    10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     "better",
				FinishTime: finishTime10,
				Bib:        10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     "worse",
				FinishTime: finishTime10,
				Bib:        10,
			},
		},
	}
	mockInStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	athletes := make(meets.AthleteLookup)
	athletes[10] = meets.NewAthlete("D", "R", "WPI", "DAID", 12, "m")

	// when the first place event comes in the builder should produce a result
	// when the updated start time comes, it should produce a new result for the
	// same bib
	expectedResults := []meets.RaceResult{
		{
			Bib:          10,
			Athlete:      athletes[10],
			Place:        1,
			PlaceSource:  t.Name(),
			Time:         time.Duration(0),
			FinishSource: "",
		},
		{
			Bib:          10,
			Athlete:      athletes[10],
			Place:        1,
			Time:         finishTime10.Sub(now),
			FinishSource: "better",
			PlaceSource:  t.Name(),
		},
	}

	mockOutStream := &stream.MockStream{
		Events: make([]stream.Message, 0, 10),
	}
	actualResults := results.NewResultStream(mockOutStream)

	builder := NewResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking["better"] = 1
	ranking["worse"] = 2
	err := builder.BuildResults(inputEvents, athletes, actualResults, ranking)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, buildActualResults(mockOutStream))
}

func TestResultBuilderRankPlaceUpdate(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()
	// Test data
	finishTime10 := now.Add(5 * time.Minute)

	testEvents := []raceevents.Event{
		{
			EventTime: now,
			Data: raceevents.StartEvent{
				Source:    t.Name(),
				StartTime: now,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     "better",
				FinishTime: finishTime10,
				Bib:        10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: "worsePlace",
				Place:  2,
				Bib:    10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: "betterPlace",
				Place:  1,
				Bib:    10,
			},
		},
	}
	mockInStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	athletes := make(meets.AthleteLookup)
	athletes[10] = meets.NewAthlete("D", "R", "WPI", "DAID", 12, "m")

	// when the first place event comes in the builder should produce a result
	// when the updated start time comes, it should produce a new result for the
	// same bib
	expectedResults := []meets.RaceResult{
		{
			Bib:          10,
			Athlete:      athletes[10],
			Place:        2,
			Time:         finishTime10.Sub(now),
			FinishSource: "better",
			PlaceSource:  "worsePlace",
		},
		{
			Bib:          10,
			Athlete:      athletes[10],
			Place:        1,
			Time:         finishTime10.Sub(now),
			FinishSource: "better",
			PlaceSource:  "betterPlace",
		},
	}

	mockOutStream := &stream.MockStream{
		Events: make([]stream.Message, 0, 10),
	}
	actualResults := results.NewResultStream(mockOutStream)

	builder := NewResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking["better"] = 1
	ranking["worse"] = 2
	ranking["betterPlace"] = 1
	ranking["worsePlace"] = 2
	err := builder.BuildResults(inputEvents, athletes, actualResults, ranking)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, buildActualResults(mockOutStream))
}

func TestResultBuilderRankPlaceIgnore(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()
	// Test data
	finishTime10 := now.Add(5 * time.Minute)

	// 3 events minimum to build a result:  start, finish and place
	// if the builder doesn't get all 3 no result for the bib is produced
	testEvents := []raceevents.Event{
		{
			EventTime: now,
			Data: raceevents.StartEvent{
				Source:    t.Name(),
				StartTime: now,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     "better",
				FinishTime: finishTime10,
				Bib:        10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: "betterPlace",
				Place:  1,
				Bib:    10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: "worsePlace",
				Place:  2,
				Bib:    10,
			},
		},
	}
	mockInStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	athletes := make(meets.AthleteLookup)
	athletes[10] = meets.NewAthlete("D", "R", "WPI", "DAID", 12, "m")

	// when the first place event comes in the builder should produce a result
	// when the updated start time comes, it should produce a new result for the
	// same bib
	expectedResults := []meets.RaceResult{
		{
			Bib:          10,
			Athlete:      athletes[10],
			Place:        1,
			Time:         finishTime10.Sub(now),
			FinishSource: "better",
			PlaceSource:  "betterPlace",
		},
	}

	mockOutStream := &stream.MockStream{
		Events: make([]stream.Message, 0, 10),
	}
	actualResults := results.NewResultStream(mockOutStream)

	builder := NewResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking["better"] = 1
	ranking["worse"] = 2
	ranking["betterPlace"] = 1
	ranking["worsePlace"] = 2
	err := builder.BuildResults(inputEvents, athletes, actualResults, ranking)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, buildActualResults(mockOutStream))
}

func TestResultBuilderPlaceLastToFirst(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()
	testEvents := []raceevents.Event{
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  1,
				Bib:    10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  2,
				Bib:    11,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  3,
				Bib:    12,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  4,
				Bib:    13,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  1,
				Bib:    13,
			},
		},
	}
	mockInStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	athletes := make(meets.AthleteLookup)
	athletes[10] = meets.NewAthlete("D", "R", "WPI", "DAID", 12, "m")
	athletes[11] = meets.NewAthlete("M", "R", "WPI", "DAID2", 12, "m")
	athletes[12] = meets.NewAthlete("MG", "R", "MVHS", "DAID3", 12, "m")
	athletes[13] = meets.NewAthlete("S", "R", "WPI", "DAID4", 16, "f")

	expectedResults := []meets.RaceResult{
		{
			Bib:         10,
			Athlete:     athletes[10],
			Place:       1,
			PlaceSource: t.Name(),
		},
		{
			Bib:         11,
			Athlete:     athletes[11],
			Place:       2,
			PlaceSource: t.Name(),
		},
		{
			Bib:         12,
			Athlete:     athletes[12],
			Place:       3,
			PlaceSource: t.Name(),
		},
		{
			Bib:         13,
			Athlete:     athletes[13],
			Place:       4,
			PlaceSource: t.Name(),
		},
		{
			Bib:         13,
			Athlete:     athletes[13],
			Place:       1,
			PlaceSource: t.Name(),
		},
		{
			Bib:         10,
			Athlete:     athletes[10],
			Place:       2,
			PlaceSource: t.Name(),
		},
		{
			Bib:         11,
			Athlete:     athletes[11],
			Place:       3,
			PlaceSource: t.Name(),
		},
		{
			Bib:         12,
			Athlete:     athletes[12],
			Place:       4,
			PlaceSource: t.Name(),
		},
	}

	mockOutStream := &stream.MockStream{
		Events: make([]stream.Message, 0, 10),
	}
	actualResults := results.NewResultStream(mockOutStream)

	builder := NewResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking[t.Name()] = 1
	err := builder.BuildResults(inputEvents, athletes, actualResults, ranking)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, buildActualResults(mockOutStream))
}

func TestResultBuilderPlaceSwapDemote(t *testing.T) {
	now := time.Now().UTC()
	testEvents := []raceevents.Event{
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  1,
				Bib:    10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  2,
				Bib:    11,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  2,
				Bib:    10,
			},
		},
	}
	mockInStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	athletes := make(meets.AthleteLookup)
	athletes[10] = meets.NewAthlete("D", "R", "WPI", "DAID", 12, "m")
	athletes[11] = meets.NewAthlete("M", "R", "WPI", "DAID2", 12, "m")

	expectedResults := []meets.RaceResult{
		{
			Bib:         10,
			Athlete:     athletes[10],
			Place:       1,
			PlaceSource: t.Name(),
		},
		{
			Bib:         11,
			Athlete:     athletes[11],
			Place:       2,
			PlaceSource: t.Name(),
		},
		{
			Bib:         11,
			Athlete:     athletes[11],
			Place:       1,
			PlaceSource: t.Name(),
		},
		{
			Bib:         10,
			Athlete:     athletes[10],
			Place:       2,
			PlaceSource: t.Name(),
		},
	}

	mockOutStream := &stream.MockStream{
		Events: make([]stream.Message, 0, 10),
	}
	actualResults := results.NewResultStream(mockOutStream)

	builder := NewResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking[t.Name()] = 1
	err := builder.BuildResults(inputEvents, athletes, actualResults, ranking)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, buildActualResults(mockOutStream))
}

func TestResultBuilderDemote(t *testing.T) {
	now := time.Now().UTC()
	testEvents := []raceevents.Event{
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  1,
				Bib:    10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  2,
				Bib:    11,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  3,
				Bib:    12,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  2,
				Bib:    10,
			},
		},
	}
	mockInStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	athletes := make(meets.AthleteLookup)
	athletes[10] = meets.NewAthlete("D", "R", "WPI", "DAID", 12, "m")
	athletes[11] = meets.NewAthlete("M", "R", "WPI", "DAID2", 12, "m")
	athletes[12] = meets.NewAthlete("MG", "R", "MVHS", "DAID3", 12, "m")

	expectedResults := []meets.RaceResult{
		{
			Bib:         10,
			Athlete:     athletes[10],
			Place:       1,
			PlaceSource: t.Name(),
		},
		{
			Bib:         11,
			Athlete:     athletes[11],
			Place:       2,
			PlaceSource: t.Name(),
		},
		{
			Bib:         12,
			Athlete:     athletes[12],
			Place:       3,
			PlaceSource: t.Name(),
		},
		{
			Bib:         11,
			Athlete:     athletes[11],
			Place:       1,
			PlaceSource: t.Name(),
		},
		{
			Bib:         10,
			Athlete:     athletes[10],
			Place:       2,
			PlaceSource: t.Name(),
		},
	}

	mockOutStream := &stream.MockStream{
		Events: make([]stream.Message, 0, 10),
	}
	actualResults := results.NewResultStream(mockOutStream)

	builder := NewResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking[t.Name()] = 1
	err := builder.BuildResults(inputEvents, athletes, actualResults, ranking)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, buildActualResults(mockOutStream))
}

func TestResultBuilderPlaceFirstToLast(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()
	testEvents := []raceevents.Event{
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  1,
				Bib:    10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  2,
				Bib:    11,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  3,
				Bib:    12,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  4,
				Bib:    13,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Place:  4,
				Bib:    10,
			},
		},
	}
	mockInStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	athletes := make(meets.AthleteLookup)
	athletes[10] = meets.NewAthlete("D", "R", "WPI", "DAID", 12, "m")
	athletes[11] = meets.NewAthlete("M", "R", "WPI", "DAID2", 12, "m")
	athletes[12] = meets.NewAthlete("MG", "R", "MVHS", "DAID3", 12, "m")
	athletes[13] = meets.NewAthlete("S", "R", "WPI", "DAID4", 16, "f")

	expectedResults := []meets.RaceResult{
		{
			Bib:         10,
			Athlete:     athletes[10],
			Place:       1,
			PlaceSource: t.Name(),
		},
		{
			Bib:         11,
			Athlete:     athletes[11],
			Place:       2,
			PlaceSource: t.Name(),
		},
		{
			Bib:         12,
			Athlete:     athletes[12],
			Place:       3,
			PlaceSource: t.Name(),
		},
		{
			Bib:         13,
			Athlete:     athletes[13],
			Place:       4,
			PlaceSource: t.Name(),
		},
		{
			Bib:         11,
			Athlete:     athletes[11],
			Place:       1,
			PlaceSource: t.Name(),
		},
		{
			Bib:         12,
			Athlete:     athletes[12],
			Place:       2,
			PlaceSource: t.Name(),
		},
		{
			Bib:         13,
			Athlete:     athletes[13],
			Place:       3,
			PlaceSource: t.Name(),
		},
		{
			Bib:         10,
			Athlete:     athletes[10],
			Place:       4,
			PlaceSource: t.Name(),
		},
	}

	mockOutStream := &stream.MockStream{
		Events: make([]stream.Message, 0, 10),
	}
	actualResults := results.NewResultStream(mockOutStream)

	builder := NewResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking[t.Name()] = 1
	err := builder.BuildResults(inputEvents, athletes, actualResults, ranking)
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, buildActualResults(mockOutStream))
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

func buildActualResults(rawOutput *stream.MockStream) []meets.RaceResult {
	result := make([]meets.RaceResult, len(rawOutput.Events))
	for i, msg := range rawOutput.Events {
		err := json.Unmarshal(msg.Data, &result[i])
		if err != nil {
			panic(err)
		}
	}

	return result
}
