package raceevents

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMarshallEventStartEvent(t *testing.T) {
	testTime := time.Now().UTC()
	startEvent := StartEvent{
		Source:    t.Name(),
		StartTime: testTime,
	}

	testEvent := Event{
		EventTime: testTime,
		Data:      startEvent,
	}

	data, err := json.Marshal(testEvent)
	assert.Nil(t, err)

	var loaded Event
	err = json.Unmarshal(data, &loaded)
	assert.Nil(t, err)
	assert.Equal(t, testEvent.EventTime, testTime)
	actualSe, typeOk := loaded.Data.(StartEvent)
	assert.True(t, typeOk)
	assert.Equal(t, testEvent.Data, actualSe)
	assert.Equal(t, t.Name(), startEvent.Source)
	assert.Equal(t, testTime, startEvent.StartTime)
}

func TestMarshallEventFinishEvent(t *testing.T) {
	testTime := time.Now().UTC()
	bib := 5
	finishEvent := FinishEvent{
		Source:     t.Name(),
		FinishTime: testTime,
		Bib:        bib,
	}

	testEvent := Event{
		EventTime: testTime,
		Data:      finishEvent,
	}

	data, err := json.Marshal(testEvent)
	assert.Nil(t, err)

	var loaded Event
	err = json.Unmarshal(data, &loaded)
	assert.Nil(t, err)
	assert.Equal(t, testEvent.EventTime, testTime)
	actualFe, typeOk := loaded.Data.(FinishEvent)
	assert.True(t, typeOk)
	assert.Equal(t, testEvent.Data, actualFe)
	assert.Equal(t, t.Name(), finishEvent.Source)
	assert.Equal(t, testTime, finishEvent.FinishTime)
}

func TestMarshallEventPlaceEvent(t *testing.T) {
	testTime := time.Now().UTC()
	bib := 5
	place := 1
	placeEvent := PlaceEvent{
		Source: t.Name(),
		Place:  place,
		Bib:    bib,
	}

	testEvent := Event{
		EventTime: testTime,
		Data:      placeEvent,
	}

	data, err := json.Marshal(testEvent)
	assert.Nil(t, err)

	var loaded Event
	err = json.Unmarshal(data, &loaded)
	assert.Nil(t, err)
	assert.Equal(t, testEvent.EventTime, testTime)
	actualPe, typeOk := loaded.Data.(PlaceEvent)
	assert.True(t, typeOk)
	assert.Equal(t, testEvent.Data, actualPe)
	assert.Equal(t, t.Name(), placeEvent.Source)
	assert.Equal(t, bib, placeEvent.Bib)
	assert.Equal(t, place, placeEvent.Place)
}
