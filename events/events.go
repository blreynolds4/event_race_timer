package events

import (
	"context"
	"time"
)

type EventType string

const (
	StartEventType  EventType = "StartEvent"
	FinishEventType EventType = "FinishEvent"
	PlaceEventType  EventType = "PlaceEvent"

	NoBib = -1
)

type RaceEvent interface {
	GetID() string
	GetSource() string
	GetType() EventType
	GetTime() time.Time
}

// Start event will have type StartEvent and Data:
// StartTime
type StartEvent interface {
	RaceEvent
	GetStartTime() time.Time
}

// Finish event will have the type FinishEvent and Data:
// Bib and Finish Time
type FinishEvent interface {
	RaceEvent
	GetFinishTime() time.Time
	GetBib() int
}

type PlaceEvent interface {
	RaceEvent
	GetBib() int
	GetPlace() int
}

// all the data for all event types is the same underneath
// so all can be sent and read as Race Events
type EventTarget interface {
	SendRaceEvent(ctx context.Context, re RaceEvent) error
}

type EventSource interface {
	GetRaceEvent(ctx context.Context, t time.Duration) (RaceEvent, error)
	GetRaceEventRange(ctx context.Context, start, end string) ([]RaceEvent, error)
}
