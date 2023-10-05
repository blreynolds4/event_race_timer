package raceevents

import (
	"blreynolds4/event-race-timer/stream"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSendStartEvent(t *testing.T) {
	mock := &stream.MockStream{}
	startTime := time.Now().UTC().Add(-time.Minute)

	es := NewEventStream(mock)

	sentEvent := StartEvent{
		Source:    t.Name(),
		StartTime: startTime,
	}

	err := es.SendStartEvent(context.TODO(), sentEvent)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mock.Events))

	var actualEvent Event
	read, err := es.GetRaceEvent(context.TODO(), 0, &actualEvent)
	assert.NoError(t, err)
	assert.True(t, read)
	actualStartEvent, isStart := actualEvent.Data.(StartEvent)
	assert.True(t, isStart)
	assert.Equal(t, sentEvent, actualStartEvent)
}

func TestSendFinishEvent(t *testing.T) {
	mock := &stream.MockStream{}
	finishTime := time.Now().UTC().Add(time.Minute)

	es := NewEventStream(mock)

	sentEvent := FinishEvent{
		Source:     t.Name(),
		FinishTime: finishTime,
	}

	err := es.SendFinishEvent(context.TODO(), sentEvent)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mock.Events))

	var actualEvent Event
	read, err := es.GetRaceEvent(context.TODO(), 0, &actualEvent)
	assert.NoError(t, err)
	assert.True(t, read)
	actualFinishEvent, isStart := actualEvent.Data.(FinishEvent)
	assert.True(t, isStart)
	assert.Equal(t, sentEvent, actualFinishEvent)
}

func TestSendPlaceEvent(t *testing.T) {
	mock := &stream.MockStream{}
	es := NewEventStream(mock)

	sentEvent := PlaceEvent{
		Source: t.Name(),
		Bib:    1,
		Place:  3,
	}

	err := es.SendPlaceEvent(context.TODO(), sentEvent)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mock.Events))

	var actualEvent Event
	read, err := es.GetRaceEvent(context.TODO(), 0, &actualEvent)
	assert.NoError(t, err)
	assert.True(t, read)
	actualPlaceEvent, isStart := actualEvent.Data.(PlaceEvent)
	assert.True(t, isStart)
	assert.Equal(t, sentEvent, actualPlaceEvent)
}

func TestGetEventFails(t *testing.T) {
	expErr := fmt.Errorf("boom")
	mock := &stream.MockStream{
		Get: func(ctx context.Context, timeout time.Duration, msg *stream.Message) (bool, error) {
			return false, expErr
		},
	}
	startTime := time.Now().UTC().Add(-time.Minute)

	es := NewEventStream(mock)

	sentEvent := StartEvent{
		Source:    t.Name(),
		StartTime: startTime,
	}

	err := es.SendStartEvent(context.TODO(), sentEvent)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mock.Events))

	var actualEvent Event
	read, err := es.GetRaceEvent(context.TODO(), 0, &actualEvent)
	assert.Equal(t, expErr, err)
	assert.False(t, read)
}

func TestGetEventNoData(t *testing.T) {
	mock := &stream.MockStream{}

	es := NewEventStream(mock)

	var actualEvent Event
	read, err := es.GetRaceEvent(context.TODO(), 0, &actualEvent)
	assert.NoError(t, err)
	assert.False(t, read)
}

func TestGetRaceEventRangeNoData(t *testing.T) {
	mock := &stream.MockStream{}

	es := NewEventStream(mock)

	events := make([]Event, 1)
	count, err := es.GetRaceEventRange(context.TODO(), "start", "end", events)
	assert.NoError(t, err)
	assert.Zero(t, count)
}

func TestGetRaceEventRangeNoCapacity(t *testing.T) {
	mock := &stream.MockStream{}

	es := NewEventStream(mock)

	events := make([]Event, 0)
	count, err := es.GetRaceEventRange(context.TODO(), "start", "end", events)
	assert.Equal(t, "can't get event range with empty buffer", err.Error())
	assert.Zero(t, count)
}

func TestGetRaceEventRange(t *testing.T) {
	mock := &stream.MockStream{}
	es := NewEventStream(mock)

	sentEvent := PlaceEvent{
		Source: t.Name(),
		Bib:    1,
		Place:  3,
	}

	err := es.SendPlaceEvent(context.TODO(), sentEvent)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mock.Events))

	events := make([]Event, 1)
	count, err := es.GetRaceEventRange(context.TODO(), "start", "end", events)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	pe, isPlaceEvent := events[0].Data.(PlaceEvent)
	assert.True(t, isPlaceEvent)
	assert.Equal(t, sentEvent, pe)
}
