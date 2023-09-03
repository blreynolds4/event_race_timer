package events

import (
	"blreynolds4/event-race-timer/stream"
	"context"
	"encoding/json"
	"fmt"
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

func TestStreamConstructors(t *testing.T) {
	actualR := NewRaceEventSource(&stream.MockStream{})
	_, ok := actualR.(*eventSourceStream)
	assert.True(t, ok)

	actualW := NewRaceEventTarget(&stream.MockStream{})
	_, ok = actualW.(*eventTargetStream)
	assert.True(t, ok)
}

func TestSendEvent(t *testing.T) {
	mock := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	evStream := NewRaceEventTarget(mock)

	se := NewStartEvent(t.Name(), time.Now().UTC())

	err := evStream.SendRaceEvent(context.TODO(), se)
	assert.NoError(t, err)

	actual := &raceEvent{}
	data, ok := mock.Events[0].Values["event"].(string)
	assert.True(t, ok)
	err = json.Unmarshal([]byte(data), actual)

	assert.NoError(t, err)
	assert.Equal(t, se, actual)
}

func TestSendEventFails(t *testing.T) {
	expErr := fmt.Errorf("fail")
	mock := &stream.MockStream{
		Events: make([]stream.Message, 0),
		Send: func(ctx context.Context, sm stream.Message) error {
			return expErr
		},
	}
	evStream := NewRaceEventTarget(mock)

	se := NewStartEvent(t.Name(), time.Now().UTC())

	err := evStream.SendRaceEvent(context.TODO(), se)
	assert.Equal(t, expErr, err)
}

func TestGetRaceEvent(t *testing.T) {
	// create the expected event.  It needs an ID, which is normally
	// added by the stream when sent.  We will add it manually
	expEvent := NewStartEvent(t.Name(), time.Now().UTC())
	msg, err := expEvent.ToStreamMessage()
	assert.NoError(t, err)
	msg.ID = "test"
	expEvent.FromStreamMessage(msg)

	mock := &stream.MockStream{
		Events: []stream.Message{msg},
	}
	evStream := NewRaceEventSource(mock)

	actualEvent, err := evStream.GetRaceEvent(context.TODO(), time.Second)
	assert.NoError(t, err)
	assert.Equal(t, expEvent, actualEvent)
}

func TestGetRaceEventEmptyStream(t *testing.T) {
	// create the expected event.  It needs an ID, which is normally
	// added by the stream when sent.  We will add it manually

	mock := &stream.MockStream{
		Events: []stream.Message{},
	}
	evStream := NewRaceEventSource(mock)

	actualEvent, err := evStream.GetRaceEvent(context.TODO(), time.Second)
	assert.NoError(t, err)
	assert.Nil(t, actualEvent)
}

func TestGetRaceEventFails(t *testing.T) {
	// create the expected event.  It needs an ID, which is normally
	// added by the stream when sent.  We will add it manually

	expErr := fmt.Errorf("fail")
	mock := &stream.MockStream{
		Get: func(ctx context.Context, timeout time.Duration) (stream.Message, error) {
			return stream.Message{}, expErr
		},
	}
	evStream := NewRaceEventSource(mock)

	_, err := evStream.GetRaceEvent(context.TODO(), time.Second)
	assert.Equal(t, expErr, err)
}

func TestGetRaceEventRange(t *testing.T) {
	// create the expected event.  It needs an ID, which is normally
	// added by the stream when sent.  We will add it manually
	expEvent := NewStartEvent(t.Name(), time.Now().UTC())
	msg, err := expEvent.ToStreamMessage()
	assert.NoError(t, err)
	msg.ID = "test"
	expEvent.FromStreamMessage(msg)
	expEvents := []RaceEvent{expEvent}

	mock := &stream.MockStream{
		Events: []stream.Message{msg},
	}
	evStream := NewRaceEventSource(mock)

	actualEvents, err := evStream.GetRaceEventRange(context.TODO(), "0", "end")
	assert.NoError(t, err)
	assert.Equal(t, expEvents, actualEvents)
}

func TestGetRaceEventRangeBadEvent(t *testing.T) {
	badMsg := stream.Message{
		Values: map[string]interface{}{
			"event": 5,
		},
	}

	mock := &stream.MockStream{
		Events: []stream.Message{badMsg},
	}
	evStream := NewRaceEventSource(mock)

	expErr := fmt.Errorf("Values data was not a string, can't build RaceEvent")
	_, err := evStream.GetRaceEventRange(context.TODO(), "0", "end")
	assert.Equal(t, expErr, err)
}
