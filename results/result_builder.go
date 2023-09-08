package results

import (
	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/events"
	"context"
	"time"
)

type ResultBuilder interface {
	BuildResults(inputEvents events.EventSource, athletes competitors.CompetitorLookup, results ResultTarget, ranking map[string]int) error
}

func NewResultBuilder() ResultBuilder {
	return &resultBuilder{}
}

type resultBuilder struct {
	eventSource events.EventSource
	results     ResultTarget
	athletes    competitors.CompetitorLookup
	ranking     map[string]int
}

func (os *resultBuilder) BuildResults(inputEvents events.EventSource, athletes competitors.CompetitorLookup, results ResultTarget, ranking map[string]int) error {

	start := []events.StartEvent{} //array to store all of the start events
	rr := map[int]RaceResult{}     //map of race results, bib number is key
	ft := map[int]time.Time{}      // map of times with bib number being key

	event, err := inputEvents.GetRaceEvent(context.TODO(), 0)

	for event != nil && err == nil {
		switch event.GetType() {
		case events.StartEventType:
			start = append(start, event.(events.StartEvent))

			// iterate over result to get times for all finish events that came before the start event
			for _, result := range rr {
				result.Time = ft[result.Bib].Sub(start[len(start)-1].GetStartTime())
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
			//when the result does not exist the result is empty
			//if the ranking of the new event source is higher than the old create a new result
			if ranking[fe.GetSource()] <= ranking[result.FinishSource] || ranking[result.FinishSource] == 0 {

				result.Bib = fe.GetBib()
				result.Athlete = athletes[fe.GetBib()]
				result.FinishSource = fe.GetSource()
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
			}
		case events.PlaceEventType:
			pe := event.(events.PlaceEvent)

			result := rr[pe.GetBib()]

			if ranking[pe.GetSource()] <= ranking[result.PlaceSource] || ranking[result.PlaceSource] == 0 {
				result.Bib = pe.GetBib()
				result.Athlete = athletes[pe.GetBib()]
				result.Place = pe.GetPlace()
				result.PlaceSource = pe.GetSource()
				rr[pe.GetBib()] = result

				if rr[pe.GetBib()].IsComplete() {
					results.SendResult(context.TODO(), rr[pe.GetBib()])
				}
			}
		}

		event, err = inputEvents.GetRaceEvent(context.TODO(), 0)
	}
	return nil
}
