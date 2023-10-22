package places

import (
	"blreynolds4/event-race-timer/internal/competitors"
	"blreynolds4/event-race-timer/internal/raceevents"
	"blreynolds4/event-race-timer/internal/stream"
	"context"
	"encoding/json"
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

	athletes := make(competitors.CompetitorLookup)
	athletes[10] = &competitors.Competitor{
		Name: "bib 10",
	}
	athletes[11] = &competitors.Competitor{
		Name: "bib 11",
	}
	athletes[12] = &competitors.Competitor{
		Name: "bib 12",
	}
	athletes[13] = &competitors.Competitor{
		Name: "bib 13",
	}

	sourceRanks := make(map[string]int)
	bestSource := t.Name()
	slowSource := t.Name() + "-slow"
	sourceRanks[bestSource] = 1
	sourceRanks[slowSource] = 2

	testEvents := []raceevents.Event{
		{
			EventTime: now,
			Data: raceevents.StartEvent{
				Source:    bestSource,
				StartTime: now,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     bestSource,
				FinishTime: finishTime10,
				Bib:        10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     bestSource,
				FinishTime: finishTime12,
				Bib:        raceevents.NoBib,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     bestSource,
				FinishTime: finishTime11,
				Bib:        11,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     bestSource,
				FinishTime: finishTime13,
				Bib:        13,
			},
		},
	}

	placesSent := make([]raceevents.Event, 0)
	mockEventStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
		Send: func(ctx context.Context, sm stream.Message) error {
			var e raceevents.Event
			err := json.Unmarshal(sm.Data, &e)
			if err != nil {
				panic(err)
			}
			placesSent = append(placesSent, e)
			return nil
		},
	}
	inputEvents := raceevents.NewEventStream(mockEventStream)

	builder := NewPlaceGenerator(inputEvents)
	err := builder.GeneratePlaces(athletes, sourceRanks)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(placesSent))

	// verify the bibs and places match what we expect
	pe, ok := placesSent[0].Data.(raceevents.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 10, pe.Bib)
	assert.Equal(t, 1, pe.Place)

	pe, ok = placesSent[1].Data.(raceevents.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 11, pe.Bib)
	assert.Equal(t, 2, pe.Place)

	pe, ok = placesSent[2].Data.(raceevents.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 13, pe.Bib)
	assert.Equal(t, 3, pe.Place)
}

func TestNormalPlacingInOrderSkipUnknownBib(t *testing.T) {
	// given a set of events on the source
	// produce the set of events on the target
	now := time.Now().UTC()
	// Test data
	finishTime10 := now.Add(5 * time.Minute)
	finishTime12 := now.Add(5*time.Minute + (time.Millisecond * 1))
	finishTime11 := now.Add(5*time.Minute + (time.Second * 5))
	finishTime13 := now.Add(5*time.Minute + (time.Second * 29))

	athletes := make(competitors.CompetitorLookup)
	athletes[10] = &competitors.Competitor{
		Name: "bib 10",
	}
	athletes[11] = &competitors.Competitor{
		Name: "bib 11",
	}
	athletes[12] = &competitors.Competitor{
		Name: "bib 12",
	}
	athletes[13] = &competitors.Competitor{
		Name: "bib 13",
	}

	sourceRanks := make(map[string]int)
	bestSource := t.Name()
	slowSource := t.Name() + "-slow"
	sourceRanks[bestSource] = 1
	sourceRanks[slowSource] = 2

	testEvents := []raceevents.Event{
		{
			EventTime: now,
			Data: raceevents.StartEvent{
				Source:    bestSource,
				StartTime: now,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     bestSource,
				FinishTime: finishTime10,
				Bib:        10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     bestSource,
				FinishTime: finishTime12,
				Bib:        999,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     bestSource,
				FinishTime: finishTime11,
				Bib:        11,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     bestSource,
				FinishTime: finishTime13,
				Bib:        13,
			},
		},
	}

	placesSent := make([]raceevents.Event, 0)
	mockEventStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
		Send: func(ctx context.Context, sm stream.Message) error {
			var e raceevents.Event
			err := json.Unmarshal(sm.Data, &e)
			if err != nil {
				panic(err)
			}
			placesSent = append(placesSent, e)
			return nil
		},
	}
	inputEvents := raceevents.NewEventStream(mockEventStream)

	builder := NewPlaceGenerator(inputEvents)
	err := builder.GeneratePlaces(athletes, sourceRanks)
	assert.NoError(t, err)
	// slice off the beginning of the event stream to get to what places were sent
	assert.Equal(t, 3, len(placesSent))

	// verify the bibs and places match what we expect
	pe, ok := placesSent[0].Data.(raceevents.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 10, pe.Bib)
	assert.Equal(t, 1, pe.Place)

	pe, ok = placesSent[1].Data.(raceevents.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 11, pe.Bib)
	assert.Equal(t, 2, pe.Place)

	pe, ok = placesSent[2].Data.(raceevents.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 13, pe.Bib)
	assert.Equal(t, 3, pe.Place)
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

	athletes := make(competitors.CompetitorLookup)
	athletes[10] = &competitors.Competitor{
		Name: "bib 10",
	}
	athletes[11] = &competitors.Competitor{
		Name: "bib 11",
	}
	athletes[12] = &competitors.Competitor{
		Name: "bib 12",
	}
	athletes[13] = &competitors.Competitor{
		Name: "bib 13",
	}

	sourceRanks := make(map[string]int)
	bestSource := t.Name()
	slowSource := t.Name() + "-slow"
	sourceRanks[bestSource] = 1
	sourceRanks[slowSource] = 2

	// multiple events per finish, should only use the best time
	// test getting best time first or second
	testEvents := []raceevents.Event{
		{
			EventTime: now,
			Data: raceevents.StartEvent{
				Source:    bestSource,
				StartTime: now,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     slowSource,
				FinishTime: finishTime10.Add(time.Second),
				Bib:        10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     bestSource,
				FinishTime: finishTime10,
				Bib:        10,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     bestSource,
				FinishTime: finishTime12,
				Bib:        raceevents.NoBib,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     bestSource,
				FinishTime: finishTime11,
				Bib:        11,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     slowSource,
				FinishTime: finishTime11.Add(time.Minute),
				Bib:        11,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     bestSource,
				FinishTime: finishTime13,
				Bib:        13,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     slowSource,
				FinishTime: finishTime13.Add(-time.Minute),
				Bib:        13,
			},
		},
	}
	placesSent := make([]raceevents.Event, 0)
	mockEventStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
		Send: func(ctx context.Context, sm stream.Message) error {
			var e raceevents.Event
			err := json.Unmarshal(sm.Data, &e)
			if err != nil {
				panic(err)
			}
			placesSent = append(placesSent, e)
			return nil
		},
	}
	inputEvents := raceevents.NewEventStream(mockEventStream)

	builder := NewPlaceGenerator(inputEvents)
	err := builder.GeneratePlaces(athletes, sourceRanks)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(placesSent))

	// verify the bibs and places match what we expect
	pe, ok := placesSent[0].Data.(raceevents.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 10, pe.Bib)
	assert.Equal(t, 1, pe.Place)

	pe, ok = placesSent[1].Data.(raceevents.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 10, pe.Bib)
	assert.Equal(t, 1, pe.Place)

	pe, ok = placesSent[2].Data.(raceevents.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 11, pe.Bib)
	assert.Equal(t, 2, pe.Place)

	pe, ok = placesSent[3].Data.(raceevents.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 13, pe.Bib)
	assert.Equal(t, 3, pe.Place)
}

func TestEventsArriveOutOfOrder(t *testing.T) {
	now := time.Now().UTC()
	// Test data
	finishTime10 := now.Add(5 * time.Minute)
	finishTime11 := now.Add(5*time.Minute + (time.Second * 5))
	finishTime13 := now.Add(5*time.Minute + (time.Second * 29))

	athletes := make(competitors.CompetitorLookup)
	athletes[10] = &competitors.Competitor{
		Name: "bib 10",
	}
	athletes[11] = &competitors.Competitor{
		Name: "bib 11",
	}
	athletes[12] = &competitors.Competitor{
		Name: "bib 12",
	}
	athletes[13] = &competitors.Competitor{
		Name: "bib 13",
	}

	sourceRanks := make(map[string]int)

	testEvents := []raceevents.Event{
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				FinishTime: finishTime13,
				Bib:        13,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				FinishTime: finishTime11,
				Bib:        11,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				FinishTime: finishTime10,
				Bib:        10,
			},
		},
	}
	placesSent := make([]raceevents.Event, 0)
	mockEventStream := &stream.MockStream{
		Events: buildEventMessages(testEvents),
		Send: func(ctx context.Context, sm stream.Message) error {
			var e raceevents.Event
			err := json.Unmarshal(sm.Data, &e)
			if err != nil {
				panic(err)
			}
			placesSent = append(placesSent, e)
			return nil
		},
	}
	inputEvents := raceevents.NewEventStream(mockEventStream)

	builder := NewPlaceGenerator(inputEvents)
	err := builder.GeneratePlaces(athletes, sourceRanks)
	assert.NoError(t, err)
	assert.Equal(t, 6, len(placesSent))

	// verify the bibs and places match what we expect
	pe, ok := placesSent[0].Data.(raceevents.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 13, pe.Bib)
	assert.Equal(t, 1, pe.Place)

	pe, ok = placesSent[1].Data.(raceevents.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 11, pe.Bib)
	assert.Equal(t, 1, pe.Place)

	pe, ok = placesSent[2].Data.(raceevents.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 13, pe.Bib)
	assert.Equal(t, 2, pe.Place)

	pe, ok = placesSent[3].Data.(raceevents.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 10, pe.Bib)
	assert.Equal(t, 1, pe.Place)

	pe, ok = placesSent[4].Data.(raceevents.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 11, pe.Bib)
	assert.Equal(t, 2, pe.Place)

	pe, ok = placesSent[5].Data.(raceevents.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 13, pe.Bib)
	assert.Equal(t, 3, pe.Place)
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
