package restorer

import (
	"blreynolds4/event-race-timer/cmd/race_archiver/internal/racearchive"
	"blreynolds4/event-race-timer/internal/raceevents"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type Restorer interface {
	Restore(r io.Reader, es raceevents.EventStream) error
}

func NewRestorer() Restorer {
	return restorer{}
}

type restorer struct{}

func (r restorer) Restore(rdr io.Reader, es raceevents.EventStream) error {
	// Decode the reader
	// send the events to the stream
	decode := json.NewDecoder(rdr)

	var archive racearchive.RaceArchive
	err := decode.Decode(&archive)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	for i := 0; i < len(archive.RaceEvents); i++ {
		e := archive.RaceEvents[i]
		switch t := e.Data.(type) {
		case raceevents.StartEvent:
			es.SendStartEvent(context.TODO(), e.Data.(raceevents.StartEvent))
		case raceevents.FinishEvent:
			es.SendFinishEvent(context.TODO(), e.Data.(raceevents.FinishEvent))
		case raceevents.PlaceEvent:
			es.SendPlaceEvent(context.TODO(), e.Data.(raceevents.PlaceEvent))
		default:
			return fmt.Errorf("unknown type in Event Data %v", t)
		}
	}

	return nil
}
