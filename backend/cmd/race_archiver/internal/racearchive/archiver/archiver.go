package archiver

import (
	"blreynolds4/event-race-timer/cmd/race_archiver/internal/racearchive"
	"blreynolds4/event-race-timer/internal/raceevents"

	"context"
	"encoding/json"
	"io"
)

const eventBufferSize = 100

type Arciver interface {
	Archive(raceevents.EventStream) error
}

type jsonFileArchiver struct {
	output io.Writer
}

func NewJsonFileArchiver(w io.Writer) Arciver {
	return jsonFileArchiver{
		output: w,
	}
}

func (jfa jsonFileArchiver) Archive(eventStream raceevents.EventStream) error {
	// read all the raceevents
	archivedEvents := make([]raceevents.Event, 0, eventBufferSize)
	raceEventsBuffer := make([]raceevents.Event, eventBufferSize)

	startId := eventStream.RangeQueryMin()
	count, err := eventStream.GetRaceEventRange(context.TODO(), startId, eventStream.RangeQueryMax(), raceEventsBuffer)
	if err != nil {
		return err
	}

	for count > 0 {
		for i := 0; i < count; i++ {
			archivedEvents = append(archivedEvents, raceEventsBuffer[i])
		}

		startId = eventStream.ExclusiveQueryStart(raceEventsBuffer[count-1].ID)
		count, err = eventStream.GetRaceEventRange(context.TODO(), startId, eventStream.RangeQueryMax(), raceEventsBuffer)
		if err != nil {
			return err
		}
	}

	encoder := json.NewEncoder(jfa.output)
	encoder.SetIndent("", "  ")
	archive := racearchive.RaceArchive{
		RaceEvents: archivedEvents,
	}
	err = encoder.Encode(archive)
	if err != nil {
		return err
	}
	return nil
}
