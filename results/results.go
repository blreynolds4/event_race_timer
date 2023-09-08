package results

import (
	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/stream"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const resultValueKey = "result"

// A RaceResult that can be filled in as events arrive and sent when IsComplete() is true
type RaceResult struct {
	Bib          int
	Athlete      *competitors.Competitor
	Place        int
	Time         time.Duration
	FinishSource string
	PlaceSource  string
}

func (rr RaceResult) IsComplete() bool {
	// the un-set, zero value for Athlete is nil because Competitor is an interface
	// if the bib is not 0, athlete is not nil, the place is not 0 and there is a duration, the result is complete
	return (rr.Bib > 0) &&
		(rr.Athlete != nil) &&
		(rr.Place > 0) &&
		(rr.Time.Milliseconds() > 0)
}

func (rr *RaceResult) ToStreamMessage() (stream.Message, error) {
	resultData, err := json.Marshal(rr)
	if err != nil {
		return stream.Message{}, err
	}

	msg := stream.Message{
		Values: map[string]interface{}{
			resultValueKey: string(resultData),
		},
	}

	return msg, nil
}

func (rr *RaceResult) FromStreamMessage(msg stream.Message) error {
	data, ok := msg.Values[resultValueKey].(string)
	if ok {
		err := json.Unmarshal([]byte(data), &rr)
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("Values data was not a string, can't build RaceResult")

}

// ResultTarget is a result publisher.  It makes results available to things like scoring
// that need to look at each result so athletes can see them.
type ResultTarget interface {
	SendResult(ctx context.Context, rr RaceResult) error
}

type ResultSource interface {
	GetResult(ctx context.Context) (RaceResult, error)
}

type resultTargetStream struct {
	rawStream stream.Writer
}

func NewResultTarget(raw stream.Writer) ResultTarget {
	return &resultTargetStream{
		rawStream: raw,
	}
}

func (rts *resultTargetStream) SendResult(ctx context.Context, rr RaceResult) error {
	msg, err := rr.ToStreamMessage()
	if err != nil {
		return err
	}

	err = rts.rawStream.SendMessage(ctx, msg)
	if err != nil {
		return err
	}

	fmt.Println("sent")
	return nil
}
