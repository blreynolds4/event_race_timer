package results

import (
	"blreynolds4/event-race-timer/internal/meets"
	"blreynolds4/event-race-timer/internal/stream"
	"context"
	"encoding/json"
	"fmt"
)

type ResultReader interface {
	GetResults(ctx context.Context, results []meets.RaceResult) (int, error)
}

type ResultWriter interface {
	SendResult(ctx context.Context, rr meets.RaceResult) error
}

type ResultStream interface {
	ResultReader
	ResultWriter
}

type resultStream struct {
	rawStream       stream.ReaderWriter
	rangeQueryStart string
}

func NewResultStream(raw stream.ReaderWriter) ResultStream {
	return &resultStream{
		rawStream:       raw,
		rangeQueryStart: raw.RangeQueryMin(), // default to earliest stream msg
	}
}

func (rs *resultStream) GetResults(ctx context.Context, results []meets.RaceResult) (int, error) {
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
		var rr meets.RaceResult
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

func (rts *resultStream) SendResult(ctx context.Context, rr meets.RaceResult) error {
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
