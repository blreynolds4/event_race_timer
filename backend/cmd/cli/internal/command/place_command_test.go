package command

import (
	"blreynolds4/event-race-timer/internal/raceevents"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlaceCommandNoBib(t *testing.T) {
	inputEvents := &raceevents.MockEventStream{
		Events: make([]raceevents.Event, 0),
	}

	eventSource := t.Name()
	place := NewPlaceCommand(eventSource, inputEvents)
	q, err := place.Run([]string{"1", "1"})
	assert.NoError(t, err)
	assert.False(t, q)
	assert.Equal(t, 1, len(inputEvents.Events))

	pe, ok := inputEvents.Events[0].Data.(raceevents.PlaceEvent)
	assert.True(t, ok)
	assert.Equal(t, 1, pe.Place)
	assert.Equal(t, 1, pe.Bib)
	assert.Equal(t, eventSource, pe.Source)
}

func TestPlaceCommandMissingArg(t *testing.T) {
	inputEvents := &raceevents.MockEventStream{
		Events: make([]raceevents.Event, 0),
	}

	place := NewPlaceCommand(t.Name(), inputEvents)
	q, err := place.Run([]string{"1"})
	assert.Error(t, err)
	assert.False(t, q)
}

func TestPlaceCommandBadBib(t *testing.T) {
	inputEvents := &raceevents.MockEventStream{
		Events: make([]raceevents.Event, 0),
	}

	place := NewPlaceCommand(t.Name(), inputEvents)
	q, err := place.Run([]string{"x", "1"})
	assert.Error(t, err)
	assert.False(t, q)
}

func TestPlaceCommandBadPlace(t *testing.T) {
	inputEvents := &raceevents.MockEventStream{
		Events: make([]raceevents.Event, 0),
	}

	place := NewPlaceCommand(t.Name(), inputEvents)
	q, err := place.Run([]string{"1", "x"})
	assert.Error(t, err)
	assert.False(t, q)
}
