package events

import (
	"blreynolds4/event-race-timer/stream"
	"context"
	"fmt"
	"time"
)

type eventSourceStream struct {
	rawStream stream.Reader
}

type eventTargetStream struct {
	rawStream stream.Writer
}

func NewRaceEventTarget(raw stream.Writer) EventTarget {
	return &eventTargetStream{
		rawStream: raw,
	}
}

func NewRaceEventSource(raw stream.Reader) EventSource {
	return &eventSourceStream{
		rawStream: raw,
	}
}

func (ets *eventTargetStream) SendRaceEvent(ctx context.Context, re RaceEvent) error {
	msg, err := re.ToStreamMessage()

	err = ets.rawStream.SendMessage(ctx, msg)
	if err != nil {
		return err
	}

	fmt.Println("sent")
	return nil
}

func (ess *eventSourceStream) GetRaceEvent(ctx context.Context, timeout time.Duration) (RaceEvent, error) {
	msg, err := ess.rawStream.GetMessage(ctx, timeout)
	if err != nil {
		return nil, err
	}

	if msg.IsValid() {
		// create a result message and deserialize
		result := new(raceEvent)
		err := result.FromStreamMessage(msg)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	return nil, nil
}

func (ess *eventSourceStream) GetRaceEventRange(ctx context.Context, start, end string) ([]RaceEvent, error) {
	msgs, err := ess.rawStream.GetMessageRange(ctx, start, end)
	if err != nil {
		return nil, err
	}

	// convert the data to RaceEvents and return them
	result := make([]RaceEvent, 0)
	for _, msg := range msgs {
		event := new(raceEvent)
		err := event.FromStreamMessage(msg)
		if err != nil {
			return result, err
		}
		result = append(result, event)
	}

	return result, nil
}
