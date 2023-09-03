package results

import (
	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/events"
	"context"
	"time"
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
	// comment on the usage of each of these variables
	start := []events.StartEvent{}
	rr := map[int]RaceResult{}
	ft := map[int]time.Time{}

	event, err := inputEvents.GetRaceEvent(context.TODO(), 0)

	for event != nil && err == nil {
		switch event.GetType() {
		case events.StartEventType:
			start = append(start, event.(events.StartEvent))

			// comment on why iterating on result here
			for _, result := range rr {
				latest_start := len(start) - 1 // look up slice syntax to get shortcut to get last item in go slice
				result.Time = ft[result.Bib].Sub(start[latest_start].GetStartTime())
				rr[result.Bib] = result

				if rr[result.Bib].IsComplete() {
					results.SendResult(context.TODO(), rr[result.Bib])
				}
			}
		case events.FinishEventType:
			// add code here to always keep the best rated finish source
			// we need a way to define event source ranking
			// will require chaning the interface (function signature)
			fe := event.(events.FinishEvent)

			result := rr[fe.GetBib()]
			result.Bib = fe.GetBib()
			result.Athlete = athletes[fe.GetBib()]
			if len(start) > 0 {
				latest_start := len(start) - 1 // use go slice for last item
				result.Time = fe.GetFinishTime().Sub(start[latest_start].GetStartTime())
			} else {
				ft[fe.GetBib()] = fe.GetFinishTime()
			}

			rr[fe.GetBib()] = result

			if rr[fe.GetBib()].IsComplete() {
				results.SendResult(context.TODO(), rr[fe.GetBib()])
			}
		case events.PlaceEventType:
			pe := event.(events.PlaceEvent)

			result := rr[pe.GetBib()]
			result.Bib = pe.GetBib()
			result.Athlete = athletes[pe.GetBib()]
			result.Place = pe.GetPlace()
			rr[pe.GetBib()] = result

			if rr[pe.GetBib()].IsComplete() {
				results.SendResult(context.TODO(), rr[pe.GetBib()])
			}
		}

		event, err = inputEvents.GetRaceEvent(context.TODO(), 0)
	}
	return nil
}
