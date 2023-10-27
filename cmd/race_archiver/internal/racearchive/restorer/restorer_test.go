package restorer

import (
	"blreynolds4/event-race-timer/cmd/race_archiver/internal/racearchive"
	"blreynolds4/event-race-timer/internal/raceevents"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDecodeFailure(t *testing.T) {
	badReader := strings.NewReader("not json")

	m := &raceevents.MockEventStream{}

	restore := NewRestorer()
	err := restore.Restore(badReader, m)
	assert.Error(t, err)
}

func TestEmptyFiile(t *testing.T) {
	badReader := strings.NewReader("")

	m := &raceevents.MockEventStream{}

	restore := NewRestorer()
	err := restore.Restore(badReader, m)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(m.Events))
}

func TestAllEventTypesSuccess(t *testing.T) {
	now := time.Now().UTC()
	expEvents := []raceevents.Event{
		{
			EventTime: now,
			Data: raceevents.StartEvent{
				Source:    t.Name(),
				StartTime: now,
			},
		},
		{
			EventTime: now,
			Data: raceevents.FinishEvent{
				Source:     t.Name(),
				Bib:        1,
				FinishTime: now,
			},
		},
		{
			EventTime: now,
			Data: raceevents.PlaceEvent{
				Source: t.Name(),
				Bib:    1,
				Place:  1,
			},
		},
	}
	data, err := json.Marshal(racearchive.RaceArchive{RaceEvents: expEvents})
	assert.NoError(t, err)

	eventReader := strings.NewReader(string(data))

	actualEvents := make([]raceevents.Event, 0, len(expEvents))

	m := &raceevents.MockEventStream{
		Events: actualEvents,
	}

	restore := NewRestorer()
	err = restore.Restore(eventReader, m)
	assert.NoError(t, err)
	assert.Equal(t, len(expEvents), len(m.Events))
}
