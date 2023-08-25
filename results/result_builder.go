package results

import (
	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/events"
)

type ResultBuilder interface {
	BuildResults(inputEvents events.EventSource, athletes competitors.CompetitorLookup, results ResultTarget) error
}

func NewResultBuilder() ResultBuilder {
	return &resultBuilder{}
}

type resultBuilder struct {
	eventSource events.EventSource
	results     ResultTarget
	athletes    competitors.CompetitorLookup
}

func (os *resultBuilder) BuildResults(inputEvents events.EventSource, athletes competitors.CompetitorLookup, results ResultTarget) error {
	start := []events.StartEvent{}
	finishes := map[int]events.FinishEvent{} //key is bib number
	places := []events.PlaceEvent{}

	incomplete := []events.PlaceEvent{}

	event, err := inputEvents.GetRaceEvent()

	for event != nil && err == nil {
		switch event.GetType() {
		case events.StartEventType:
			start = append(start, event.(events.StartEvent))
		case events.FinishEventType:
			finishes[event.(events.FinishEvent).GetBib()] = event.(events.FinishEvent)
		case events.PlaceEventType:
			places = append(places, event.(events.PlaceEvent))
		}

		event, err = inputEvents.GetRaceEvent()
	}

	for _, event := range places {
		rr := RaceResult{event.GetBib(), athletes[event.GetBib()], event.GetPlace(), finishes[event.GetBib()].GetFinishTime().Sub(start[0].GetStartTime())}
		if rr.IsComplete() {
			results.SendResult(rr)
			delete(finishes, event.GetBib())
		} else {
			incomplete = append(incomplete, event)
		}
	}
	places = incomplete

	return nil
}
