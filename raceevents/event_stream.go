package raceevents

import (
	"blreynolds4/event-race-timer/stream"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Define the Event Stream, it can support reading and writing
// Event Stream is the api for sending/getting events
type EventStream struct {
	raw stream.ReaderWriter
}

func NewEventStream(s stream.ReaderWriter) *EventStream {
	return &EventStream{
		raw: s,
	}
}

func (es *EventStream) SendStartEvent(ctx context.Context, se StartEvent) error {
	// wrap start event with event and send
	return es.sendMessage(ctx, Event{
		EventTime: time.Now().UTC(),
		Data:      se,
	})
}

func (es *EventStream) SendFinishEvent(ctx context.Context, fe FinishEvent) error {
	// wrap finish event with event and send
	return es.sendMessage(ctx, Event{
		EventTime: time.Now().UTC(),
		Data:      fe,
	})
}

func (es *EventStream) SendPlaceEvent(ctx context.Context, pe PlaceEvent) error {
	// wrap place event with event and send
	return es.sendMessage(ctx, Event{
		EventTime: time.Now().UTC(),
		Data:      pe,
	})
}

func (es *EventStream) sendMessage(ctx context.Context, e Event) error {
	eventData, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return es.raw.SendMessage(ctx, stream.Message{
		Data: eventData,
	})
}

func (es *EventStream) GetRaceEvent(ctx context.Context, timeout time.Duration, e *Event) (bool, error) {
	var msg stream.Message
	read, err := es.raw.GetMessage(ctx, timeout, &msg)
	if err != nil {
		return false, err
	}
	if !read {
		return false, nil
	}

	err = json.Unmarshal(msg.Data, e)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (es EventStream) GetRaceEventRange(ctx context.Context, startId, endId string, resultEvents []Event) (int, error) {
	// may not be ideal to allocate a slice of Messages in same size
	// as Events, but simplest way to do this
	if len(resultEvents) == 0 {
		return 0, fmt.Errorf("can't get event range with empty buffer")
	}

	msgs := make([]stream.Message, len(resultEvents))
	count, err := es.raw.GetMessageRange(ctx, startId, endId, msgs)
	if err != nil {
		return 0, err
	}
	if count == 0 {
		return 0, nil
	}

	read := 0
	for i := 0; i < count; i++ {
		var e Event
		err = json.Unmarshal(msgs[i].Data, &e)
		if err != nil {
			return 0, err
		}
		resultEvents[i] = e
		read++
	}

	return read, nil
}
