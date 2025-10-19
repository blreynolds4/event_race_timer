package archiver

import (
	"blreynolds4/event-race-timer/cmd/race_archiver/internal/racearchive"
	"blreynolds4/event-race-timer/internal/raceevents"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFirstRangeFails(t *testing.T) {
	expErr := fmt.Errorf("boom")
	// create a mock EventStream
	// where getrange fails
	mock := &raceevents.MockEventStream{
		Range: func(ctx context.Context, startId, endId string, msgs []raceevents.Event) (int, error) {
			return 0, expErr
		},
	}

	w := &strings.Builder{}
	archiver := NewJsonFileArchiver(w)
	err := archiver.Archive(mock)
	assert.Error(t, err)
	assert.Equal(t, 0, w.Len())
}

func TestSecondRangeFails(t *testing.T) {
	expErr := fmt.Errorf("boom")
	callCount := 0
	// create a mock EventStream
	// where getrange fails after first call
	mock := &raceevents.MockEventStream{
		Range: func(ctx context.Context, startId, endId string, events []raceevents.Event) (int, error) {
			if callCount < 1 {
				callCount++
				events[0] = raceevents.Event{}
				return 1, nil
			}
			return 0, expErr
		},
	}

	w := &strings.Builder{}
	archiver := NewJsonFileArchiver(w)
	err := archiver.Archive(mock)
	assert.Error(t, err)
	assert.Equal(t, 0, w.Len())
}

func TestFirstRangeReturnsZero(t *testing.T) {
	mock := &raceevents.MockEventStream{}

	w := &strings.Builder{}
	archiver := NewJsonFileArchiver(w)
	err := archiver.Archive(mock)
	assert.NoError(t, err)

	assert.True(t, w.Len() > 0)

	actual := racearchive.RaceArchive{}
	err = json.Unmarshal([]byte(w.String()), &actual)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(actual.RaceEvents))
}

func TestEncodeFails(t *testing.T) {
	expErr := fmt.Errorf("boom")
	mock := &raceevents.MockEventStream{}

	w := badWriter{e: expErr}
	archiver := NewJsonFileArchiver(w)
	err := archiver.Archive(mock)
	assert.Error(t, err)
}

func TestSuccess(t *testing.T) {
	expectedArchive := racearchive.RaceArchive{
		RaceEvents: []raceevents.Event{
			{
				ID:        "id",
				EventTime: time.Now().UTC(),
				Data: raceevents.StartEvent{
					Source:    t.Name(),
					StartTime: time.Now().UTC(),
				},
			},
		},
	}

	mock := &raceevents.MockEventStream{
		Events: expectedArchive.RaceEvents,
	}

	w := &strings.Builder{}
	archiver := NewJsonFileArchiver(w)
	err := archiver.Archive(mock)
	assert.NoError(t, err)

	assert.True(t, w.Len() > 0)

	actualArchive := racearchive.RaceArchive{}
	err = json.Unmarshal([]byte(w.String()), &actualArchive)
	assert.NoError(t, err)

	assert.Equal(t, expectedArchive, actualArchive)
}

type badWriter struct {
	e error
}

func (bw badWriter) Write(p []byte) (n int, err error) {
	if bw.e != nil {
		return 0, bw.e
	}

	return 0, nil
}
