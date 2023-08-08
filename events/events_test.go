package events

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// unit test the events stuff
func TestCreateStartEvent(t *testing.T) {
	source := "testSource"
	startTime := time.Now().UTC()
	startEvent := NewStartEvent(source, startTime)
	assert.NotNil(t, startEvent)

	assert.Equal(t, StartEventType, startEvent.GetType())
	assert.Equal(t, startTime, startEvent.GetTime())
	assert.Equal(t, startTime, startEvent.GetStartTime())
	assert.Equal(t, source, startEvent.GetSource())
}

func TestGetStartTimeMissing(t *testing.T) {
	startEvent := &raceEvent{
		Type: StartEventType,
		Data: make(map[string]interface{}),
	}

	assert.Panics(t, func() { startEvent.GetStartTime() }, "GetStartTime did not panic")
}

func TestGetStartTimeNotADate(t *testing.T) {
	startEvent := &raceEvent{
		Type: StartEventType,
		Data: make(map[string]interface{}),
	}

	startEvent.Data[startTimeData] = "not a date"

	assert.Panics(t, func() { startEvent.GetStartTime() }, "GetStartTime did not panic")
}

func TestGetBibMissing(t *testing.T) {
	finishEvent := &raceEvent{
		Type: FinishEventType,
		Data: make(map[string]interface{}),
	}

	assert.Panics(t, func() { finishEvent.GetBib() }, "GetBib did not panic")
}

func TestGetBibNotAnInt(t *testing.T) {
	finishEvent := &raceEvent{
		Type: FinishEventType,
		Data: make(map[string]interface{}),
	}

	finishEvent.Data[bibData] = "not an int"

	assert.Panics(t, func() { finishEvent.GetBib() }, "GetBib did not panic")
}

func TestCreateFinishEvent(t *testing.T) {
	source := "testSource"
	finishTime := time.Now().UTC()
	bib := 5
	finishEvent := NewFinishEvent(source, finishTime, bib)
	assert.NotNil(t, finishEvent)

	assert.Equal(t, FinishEventType, finishEvent.GetType())
	assert.Equal(t, finishTime, finishEvent.GetTime())
	assert.Equal(t, finishTime, finishEvent.GetFinishTime())
	assert.Equal(t, bib, finishEvent.GetBib())
	assert.Equal(t, source, finishEvent.GetSource())
}

func TestCreatePlaceEvent(t *testing.T) {
	source := "testSource"
	bib := 5
	place := 1
	placeEvent := NewPlaceEvent(source, bib, place)
	assert.NotNil(t, placeEvent)

	assert.Equal(t, PlaceEventType, placeEvent.GetType())
	assert.Equal(t, bib, placeEvent.GetBib())
	assert.Equal(t, place, placeEvent.GetPlace())
	assert.Equal(t, source, placeEvent.GetSource())
}

func TestGetPlaceMissing(t *testing.T) {
	placeEvent := &raceEvent{
		Type: PlaceEventType,
		Data: make(map[string]interface{}),
	}

	assert.Panics(t, func() { placeEvent.GetPlace() }, "GetPlace did not panic")
}

func TestGetPlaceNotAnInt(t *testing.T) {
	placeEvent := &raceEvent{
		Type: PlaceEventType,
		Data: make(map[string]interface{}),
	}

	placeEvent.Data[placeData] = "not an int"

	assert.Panics(t, func() { placeEvent.GetPlace() }, "GetPlace did not panic")
}

func TestMarshallEvent(t *testing.T) {
	source := "testSource"
	bib := 5
	place := 1
	placeEvent := NewPlaceEvent(source, bib, place)
	assert.NotNil(t, placeEvent)

	assert.Equal(t, PlaceEventType, placeEvent.GetType())
	assert.Equal(t, bib, placeEvent.GetBib())
	assert.Equal(t, place, placeEvent.GetPlace())
	assert.Equal(t, source, placeEvent.GetSource())

	data, err := json.Marshal(placeEvent)
	assert.Nil(t, err)

	var loaded raceEvent
	err = json.Unmarshal(data, &loaded)
	assert.Nil(t, err)

	assert.Equal(t, PlaceEventType, loaded.GetType())
	assert.Equal(t, bib, loaded.GetBib())
	assert.Equal(t, place, loaded.GetPlace())
	assert.Equal(t, source, loaded.GetSource())
}
