package resultbuilder

import (
	"blreynolds4/event-race-timer/internal/meets"
	"blreynolds4/event-race-timer/internal/raceevents"
	"log/slog"

	"blreynolds4/event-race-timer/internal/stream"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRaceResultBuilderNoStartEventNoResult(t *testing.T) {
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
	}

	mockInStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
	}
	inputEvents := raceevents.NewEventStream(mockInStream)

	athletes := make(meets.AthleteLookup)
	athletes[10] = meets.NewAthlete("D", "R", "WPI", "DAID", 12, "m")

	expectedResults := []meets.RaceResult{}

	builder := NewRaceResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking[t.Name()] = 1

	mockResults := meets.NewMockResultWriter()

	err := builder.BuildRaceResults(inputEvents, athletes, ranking, mockResults)
	assert.NoError(t, err)

	assert.Equal(t, len(expectedResults), len(mockResults.SavedResults))
}

func TestRaceResultBuilderFinishThenStartSaveTime(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()
	// Test data
	expectedDurationFinishTime := (5 * time.Minute)
	finishTime10 := now.Add(expectedDurationFinishTime)

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
			Place:        0,
			Time:         expectedDurationFinishTime,
			FinishSource: t.Name(),
		},
	}

	builder := NewRaceResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking[t.Name()] = 1

	mockResults := meets.NewMockResultWriter()

	err := builder.BuildRaceResults(inputEvents, athletes, ranking, mockResults)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(mockResults.SavedResults))
	assert.Equal(t, expectedResults[0], mockResults.SavedResults[0])
}

func TestRaceResultBuilderMuiltpleStartsAlwaysSaveResultFromFirstStart(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()
	// Test data
	expectedDurationFinishTime := (5 * time.Minute)
	finishTime10 := now.Add(expectedDurationFinishTime)
	startTimeUpdate := now.Add(-expectedDurationFinishTime)

	// 3 events minimum to build a result:  start, finish and place
	// if the builder doesn't get all 3 no result for the bib is produced
	testEvents := []raceevents.Event{
		{
			ID:        "1",
			EventTime: now,
			Data: raceevents.StartEvent{
				Source:    t.Name(),
				StartTime: now,
			},
		},
		{
			ID:        "2",
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     t.Name(),
				Bib:        10,
				FinishTime: finishTime10,
			},
		},
		{
			ID:        "3",
			EventTime: now,
			Data: raceevents.StartEvent{
				Source:    t.Name(),
				StartTime: startTimeUpdate,
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
			Place:        0,
			Time:         expectedDurationFinishTime,
			FinishSource: t.Name(),
		},
	}

	builder := NewRaceResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking[t.Name()] = 1

	mockResults := meets.NewMockResultWriter()

	err := builder.BuildRaceResults(inputEvents, athletes, ranking, mockResults)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(mockResults.SavedResults))
	assert.Equal(t, expectedResults[0], mockResults.SavedResults[0])
}

func TestRaceResultBuilderStartFinishPlace(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()
	// Test data
	expectedDurationFinishTime := (5 * time.Minute)
	finishTime10 := now.Add(expectedDurationFinishTime)

	// 3 events minimum to build a result:  start, finish and place
	// if the builder doesn't get all 3 no result for the bib is produced
	testEvents := []raceevents.Event{
		{
			ID:        "1",
			EventTime: now,
			Data: raceevents.StartEvent{
				Source:    t.Name(),
				StartTime: now,
			},
		},
		{
			ID:        "2",
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     t.Name(),
				Bib:        10,
				FinishTime: finishTime10,
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
			Place:        0,
			Time:         expectedDurationFinishTime,
			FinishSource: t.Name(),
		},
		{
			Bib:          10,
			Athlete:      athletes[10],
			Place:        1,
			Time:         expectedDurationFinishTime,
			FinishSource: t.Name(),
			PlaceSource:  t.Name(),
		},
	}

	builder := NewRaceResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking[t.Name()] = 1

	mockResults := meets.NewMockResultWriter()

	err := builder.BuildRaceResults(inputEvents, athletes, ranking, mockResults)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(mockResults.SavedResults))
	assert.Equal(t, expectedResults[0], mockResults.SavedResults[0])
	assert.Equal(t, expectedResults[1], mockResults.SavedResults[1])
}

func TestRaceResultBuilderPlaceOnly(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()

	// 3 events minimum to build a result:  start, finish and place
	// if the builder doesn't get all 3 no result for the bib is produced
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

	builder := NewRaceResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking[t.Name()] = 1

	mockResults := meets.NewMockResultWriter()

	err := builder.BuildRaceResults(inputEvents, athletes, ranking, mockResults)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(mockResults.SavedResults))
	assert.Equal(t, expectedResults[0], mockResults.SavedResults[0])
}

func TestRaceResultBuilderFinishTimeOverride(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()
	// Test data
	expected5DurationFinishTime := (5 * time.Minute)
	expected10DurationFinishTime := (10 * time.Minute)
	finishTime5 := now.Add(expected5DurationFinishTime)
	finishTime10 := now.Add(expected10DurationFinishTime)

	// 3 events minimum to build a result:  start, finish and place
	// if the builder doesn't get all 3 no result for the bib is produced
	testEvents := []raceevents.Event{
		{
			ID:        "1",
			EventTime: now,
			Data: raceevents.StartEvent{
				Source:    t.Name(),
				StartTime: now,
			},
		},
		{
			ID:        "2",
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     "good",
				Bib:        10,
				FinishTime: finishTime10,
			},
		},
		{
			ID:        "2",
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     "better",
				Bib:        10,
				FinishTime: finishTime5,
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
			Place:        0,
			Time:         expected10DurationFinishTime,
			FinishSource: "good",
		},
		{
			Bib:          10,
			Athlete:      athletes[10],
			Place:        0,
			Time:         expected5DurationFinishTime,
			FinishSource: "better",
		},
	}

	builder := NewRaceResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking["good"] = 2
	ranking["better"] = 1

	mockResults := meets.NewMockResultWriter()

	err := builder.BuildRaceResults(inputEvents, athletes, ranking, mockResults)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(mockResults.SavedResults))
	assert.Equal(t, expectedResults[0], mockResults.SavedResults[0])
	assert.Equal(t, expectedResults[1], mockResults.SavedResults[1])
}

func TestRaceResultBuilderFinishTimeSkipUpdate(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()
	// Test data
	expected5DurationFinishTime := (5 * time.Minute)
	expected10DurationFinishTime := (10 * time.Minute)
	finishTime5 := now.Add(expected5DurationFinishTime)
	finishTime10 := now.Add(expected10DurationFinishTime)

	// 3 events minimum to build a result:  start, finish and place
	// if the builder doesn't get all 3 no result for the bib is produced
	testEvents := []raceevents.Event{
		{
			ID:        "1",
			EventTime: now,
			Data: raceevents.StartEvent{
				Source:    t.Name(),
				StartTime: now,
			},
		},
		{
			ID:        "2",
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     "better",
				Bib:        10,
				FinishTime: finishTime10,
			},
		},
		{
			ID:        "2",
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     "good",
				Bib:        10,
				FinishTime: finishTime5,
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
			Place:        0,
			Time:         expected10DurationFinishTime,
			FinishSource: "better",
		},
	}

	builder := NewRaceResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking["good"] = 2
	ranking["better"] = 1

	mockResults := meets.NewMockResultWriter()

	err := builder.BuildRaceResults(inputEvents, athletes, ranking, mockResults)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(mockResults.SavedResults))
	assert.Equal(t, expectedResults[0], mockResults.SavedResults[0])
}

func TestRaceResultBuilderPlaceOverride(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()

	// 3 events minimum to build a result:  start, finish and place
	// if the builder doesn't get all 3 no result for the bib is produced
	testEvents := []raceevents.Event{
		{
			ID:        "3",
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: "good",
				Bib:    10,
				Place:  2,
			},
		},
		{
			ID:        "4",
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: "better",
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
			Place:       2,
			PlaceSource: "good",
		},
		{
			Bib:         10,
			Athlete:     athletes[10],
			Place:       1,
			PlaceSource: "better",
		},
	}

	builder := NewRaceResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking["good"] = 2
	ranking["better"] = 1

	mockResults := meets.NewMockResultWriter()

	err := builder.BuildRaceResults(inputEvents, athletes, ranking, mockResults)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(mockResults.SavedResults))
	assert.Equal(t, expectedResults[0], mockResults.SavedResults[0])
	assert.Equal(t, expectedResults[1], mockResults.SavedResults[1])
}

func TestRaceResultBuilderPlaceSkipUpdate(t *testing.T) {
	// read events off a stream and return
	// result events when they are complete
	now := time.Now().UTC()

	// 3 events minimum to build a result:  start, finish and place
	// if the builder doesn't get all 3 no result for the bib is produced
	testEvents := []raceevents.Event{
		{
			ID:        "3",
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: "better",
				Bib:    10,
				Place:  2,
			},
		},
		{
			ID:        "4",
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: "good",
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
			Place:       2,
			PlaceSource: "better",
		},
	}

	builder := NewRaceResultBuilder(slog.Default())
	ranking := map[string]int{}
	ranking["good"] = 2
	ranking["better"] = 1

	mockResults := meets.NewMockResultWriter()

	err := builder.BuildRaceResults(inputEvents, athletes, ranking, mockResults)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(mockResults.SavedResults))
	assert.Equal(t, expectedResults[0], mockResults.SavedResults[0])
}
