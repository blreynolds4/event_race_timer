package raceevents

import (
	"context"
	"time"
)

type MockEventStream struct {
	SendWorkout func(ctx context.Context, se WorkoutEvent) error
	SendStart   func(ctx context.Context, se StartEvent) error
	SendFinish  func(ctx context.Context, fe FinishEvent) error
	SendPlace   func(ctx context.Context, pe PlaceEvent) error
	Get         func(ctx context.Context, timeout time.Duration, msg *Event) (bool, error)
	Range       func(ctx context.Context, startId, endId string, msgs []Event) (int, error)
	Events      []Event
}

func (mes *MockEventStream) SendStartEvent(ctx context.Context, se StartEvent) error {
	if mes.SendStart != nil {
		return mes.SendStart(ctx, se)
	}

	mes.Events = append(mes.Events, Event{
		EventTime: se.StartTime,
		Data:      se,
	})

	return nil
}

func (mes *MockEventStream) SendFinishEvent(ctx context.Context, fe FinishEvent) error {
	if mes.SendStart != nil {
		return mes.SendFinish(ctx, fe)
	}

	mes.Events = append(mes.Events, Event{
		EventTime: fe.FinishTime,
		Data:      fe,
	})
	return nil
}

func (mes *MockEventStream) SendPlaceEvent(ctx context.Context, pe PlaceEvent) error {
	if mes.SendStart != nil {
		return mes.SendPlace(ctx, pe)
	}

	mes.Events = append(mes.Events, Event{
		EventTime: time.Now().UTC(),
		Data:      pe,
	})
	return nil
}

func (mes *MockEventStream) SendWorkoutEvent(ctx context.Context, we WorkoutEvent) error {
	if mes.SendStart != nil {
		return mes.SendWorkout(ctx, we)
	}

	mes.Events = append(mes.Events, Event{
		EventTime: time.Now().UTC(),
		Data:      we,
	})
	return nil
}

func (mes *MockEventStream) GetRaceEvent(ctx context.Context, timeout time.Duration, e *Event) (bool, error) {
	if mes.Get != nil {
		return mes.Get(ctx, timeout, e)
	}
	if len(mes.Events) > 0 {
		*e = mes.Events[0]
		mes.Events = mes.Events[1:]
		return true, nil
	}

	return false, nil
}

func (mes *MockEventStream) RangeQueryMin() string {
	return "-"
}

func (mes *MockEventStream) ExclusiveQueryStart(string) string {
	return "("
}

func (mes *MockEventStream) RangeQueryMax() string {
	return "+"
}

func (mes *MockEventStream) GetRaceEventRange(ctx context.Context, startId, endId string, events []Event) (int, error) {
	if mes.Range != nil {
		return mes.Range(ctx, startId, endId, events)
	}
	if len(mes.Events) > 0 {
		events[0] = mes.Events[0]
		mes.Events = mes.Events[1:]
		// returning 1 because it should be count returned, not buffer size
		// mock only returns 1 at a time
		return 1, nil
	}

	return 0, nil
}
