package results

import (
	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/raceevents"
	"context"
	"time"
)

type ResultBuilder interface {
	BuildResults(inputEvents *raceevents.EventStream, athletes competitors.CompetitorLookup, results *ResultStream, ranking map[string]int) error
}

func NewResultBuilder() ResultBuilder {
	return &resultBuilder{}
}

type resultBuilder struct {
}

func (os *resultBuilder) BuildResults(inputEvents *raceevents.EventStream, athletes competitors.CompetitorLookup, results *ResultStream, ranking map[string]int) error {

	start := []raceevents.StartEvent{} //array to store all of the start events
	rr := map[int]RaceResult{}         //map of race results, bib number is key
	ft := map[int]time.Time{}          // map of times with bib number being key

	var event raceevents.Event
	gotEvent, err := inputEvents.GetRaceEvent(context.TODO(), 0, &event)
	if err != nil {
		return err
	}

	for gotEvent {
		switch event.Data.(type) {
		case raceevents.StartEvent:
			se := event.Data.(raceevents.StartEvent)
			start = append(start, se)

			// iterate over result to get times for all finish events that came before the start event
			for _, result := range rr {
				result.Time = ft[result.Bib].Sub(start[len(start)-1].StartTime)
				rr[result.Bib] = result

				if rr[result.Bib].IsComplete() {
					results.SendResult(context.TODO(), rr[result.Bib])
				}
			}
		case raceevents.FinishEvent:
			// add code here to always keep the best rated finish source
			// we need a way to define event source ranking
			// will require chaning the interface (function signature)
			fe := event.Data.(raceevents.FinishEvent)

			result := rr[fe.Bib]
			//when the result does not exist the result is empty
			//if the ranking of the new event source is higher than the old create a new result
			if ranking[fe.Source] <= ranking[result.FinishSource] || ranking[result.FinishSource] == 0 {

				result.Bib = fe.Bib
				result.Athlete = athletes[fe.Bib]
				result.FinishSource = fe.Source
				if len(start) > 0 {
					latest_start := len(start) - 1 // use go slice for last item
					result.Time = fe.FinishTime.Sub(start[latest_start].StartTime)
				} else {
					ft[fe.Bib] = fe.FinishTime
				}
				rr[fe.Bib] = result

				if rr[fe.Bib].IsComplete() {
					results.SendResult(context.TODO(), rr[fe.Bib])
				}
			}
		case raceevents.PlaceEvent:
			pe := event.Data.(raceevents.PlaceEvent)

			result := rr[pe.Bib]

			if ranking[pe.Source] <= ranking[result.PlaceSource] || ranking[result.PlaceSource] == 0 {
				result.Bib = pe.Bib
				result.Athlete = athletes[pe.Bib]
				result.Place = pe.Place
				result.PlaceSource = pe.Source
				rr[pe.Bib] = result

				if rr[pe.Bib].IsComplete() {
					results.SendResult(context.TODO(), rr[pe.Bib])
				}
			}
		}

		gotEvent, err = inputEvents.GetRaceEvent(context.TODO(), 0, &event)
	}
	return nil
}
