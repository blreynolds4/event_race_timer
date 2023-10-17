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
	// the un-set, zero value for Athlete is nil
	// if the bib is not 0, athlete is not nil, the place is not 0 and has a source
	// then the result can be published
	return (rr.Bib > 0) &&
		(rr.Athlete != nil) &&
		(rr.Place > 0) &&
		(rr.PlaceSource != "")
}

type ResultStream struct {
	rawStream       stream.ReaderWriter
	rangeQueryStart string
}

func NewResultStream(raw stream.ReaderWriter) *ResultStream {
	return &ResultStream{
		rawStream:       raw,
		rangeQueryStart: raw.RangeQueryMin(), // default to earliest stream msg
	}
}

func (rs *ResultStream) GetResults(ctx context.Context, results []RaceResult) (int, error) {
	if len(results) == 0 {
		return 0, fmt.Errorf("can't get results with zero length buffer")
	}

	msgBuffer := make([]stream.Message, len(results))
	count, err := rs.rawStream.GetMessageRange(ctx, rs.rangeQueryStart, rs.rawStream.RangeQueryMax(), msgBuffer)
	if err != nil {
		return 0, err
	}

	for i := 0; i < count; i++ {
		// create a result message and deserialize
		// using the temp copy here means that the athlete pointer
		// doesn't get shared with each read
		// (assigning to results[i] directly only kept last competitor for all)
		// using rr creates new space that escapes into the result
		var rr RaceResult
		err = json.Unmarshal(msgBuffer[i].Data, &rr)
		if err != nil {
			return 0, err
		}
		results[i] = rr

		// update query start id to start with last read msg but not include it in next result
		rs.rangeQueryStart = rs.rawStream.ExclusiveQueryStart(msgBuffer[i].ID)
	}

	return count, nil
}

func (rts *ResultStream) SendResult(ctx context.Context, rr RaceResult) error {
	resData, err := json.Marshal(rr)
	if err != nil {
		return err
	}

	err = rts.rawStream.SendMessage(ctx, stream.Message{
		Data: resData,
	})
	if err != nil {
		return err
	}

	return nil
}
