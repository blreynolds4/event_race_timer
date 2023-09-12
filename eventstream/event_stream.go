package eventstream

import (
	"blreynolds4/event-race-timer/events"
	"blreynolds4/event-race-timer/stream"
	"context"
	"fmt"
	"time"
)

type StreamMessageToRaceEventFunc func(stream.Message) (events.RaceEvent, error)
type RaceEventToStreamMessageFunc func(events.RaceEvent) (stream.Message, error)

type eventSourceStream struct {
	rawStream  stream.Reader
	conversion StreamMessageToRaceEventFunc
}

type eventTargetStream struct {
	rawStream  stream.Writer
	conversion RaceEventToStreamMessageFunc
}

func NewRaceEventTarget(raw stream.Writer, conv RaceEventToStreamMessageFunc) events.EventTarget {
	return &eventTargetStream{
		rawStream:  raw,
		conversion: conv,
	}
}

func NewRaceEventSource(raw stream.Reader, conv StreamMessageToRaceEventFunc) events.EventSource {
	return &eventSourceStream{
		rawStream:  raw,
		conversion: conv,
	}
}

func (ets *eventTargetStream) SendRaceEvent(ctx context.Context, re events.RaceEvent) error {
	msg, err := ets.conversion(re)
	if err != nil {
		return err
	}

	err = ets.rawStream.SendMessage(ctx, msg)
	if err != nil {
		return err
	}

	fmt.Println("sent")
	return nil
}

func (ess *eventSourceStream) GetRaceEvent(ctx context.Context, timeout time.Duration) (events.RaceEvent, error) {
	msg, err := ess.rawStream.GetMessage(ctx, timeout)
	if err != nil {
		return nil, err
	}

	if msg.IsValid() {
		// create a result message and deserialize
		result, err := ess.conversion(msg)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	return nil, nil
}

func (ess *eventSourceStream) GetRaceEventRange(ctx context.Context, start, end string) ([]events.RaceEvent, error) {
	msgs, err := ess.rawStream.GetMessageRange(ctx, start, end)
	if err != nil {
		return nil, err
	}

	// convert the data to RaceEvents and return them
	result := make([]events.RaceEvent, 0)
	for _, msg := range msgs {
		event, err := ess.conversion(msg)
		if err != nil {
			return result, err
		}
		result = append(result, event)
	}

	return result, nil
}
