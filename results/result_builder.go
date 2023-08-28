package results

import (
	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/events"
	"fmt"
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
	start := []events.StartEvent{}

	rr := map[int]RaceResult{}
	ft := map[int]time.Time{}

	event, err := inputEvents.GetRaceEvent()

	for event != nil && err == nil {
		switch event.GetType() {
		case events.StartEventType:
			fmt.Println("start event")
			start = append(start, event.(events.StartEvent))

			for _, result := range rr {
				latest_start := len(start) - 1
				result.Time = ft[result.Bib].Sub(start[latest_start].GetStartTime())
				rr[result.Bib] = result

				if rr[result.Bib].IsComplete() {
					fmt.Println("send result")
					results.SendResult(rr[result.Bib])
				}
			}
		case events.FinishEventType:
			fe := event.(events.FinishEvent)

			result := rr[fe.GetBib()]
			result.Bib = fe.GetBib()
			result.Athlete = athletes[fe.GetBib()]
			fmt.Println("finishEvent2")
			if len(start) > 0 {
				fmt.Println("finishEvent2.1")
				latest_start := len(start) - 1
				result.Time = fe.GetFinishTime().Sub(start[latest_start].GetStartTime())
			} else {
				fmt.Println("finishEvent2.2")
				ft[fe.GetBib()] = fe.GetFinishTime()
			}

			fmt.Println("finishEvent3")
			rr[fe.GetBib()] = result

			if rr[fe.GetBib()].IsComplete() {
				fmt.Println("sent")
				results.SendResult(rr[fe.GetBib()])
			}
		case events.PlaceEventType:
			pe := event.(events.PlaceEvent)

			result := rr[pe.GetBib()]
			result.Bib = pe.GetBib()
			result.Athlete = athletes[pe.GetBib()]
			result.Place = pe.GetPlace()
			rr[pe.GetBib()] = result

			if rr[pe.GetBib()].IsComplete() {
				fmt.Println("sent")
				results.SendResult(rr[pe.GetBib()])
			}
		}

		event, err = inputEvents.GetRaceEvent()
	}
	return nil
}