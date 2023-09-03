package events

import (
	"context"
	"time"
)

type MockRaceEventStream struct {
	Send   func(ctx context.Context, re RaceEvent) error
	Get    func(ctx context.Context, t time.Duration) (RaceEvent, error)
	Range  func(ctx context.Context, start, end string) ([]RaceEvent, error)
	Events []RaceEvent
}

func (mres *MockRaceEventStream) SendRaceEvent(ctx context.Context, re RaceEvent) error {
	if mres.Send != nil {
		return mres.Send(ctx, re)
	}

	mres.Events = append(mres.Events, re)
	return nil
}

func (mres *MockRaceEventStream) GetRaceEvent(ctx context.Context, t time.Duration) (RaceEvent, error) {
	if mres.Get != nil {
		return mres.Get(ctx, t)
	}

	if len(mres.Events) > 0 {
		result := mres.Events[0]
		mres.Events = mres.Events[1:]
		return result, nil
	}
	return nil, nil
}

func (mres *MockRaceEventStream) GetRaceEventRange(ctx context.Context, start, end string) ([]RaceEvent, error) {
	if mres.Range != nil {
		return mres.Range(ctx, start, end)
	}

	if len(mres.Events) > 0 {
		result := mres.Events[0]
		mres.Events = mres.Events[1:]
		return []RaceEvent{result}, nil
	}
	return nil, nil
}
