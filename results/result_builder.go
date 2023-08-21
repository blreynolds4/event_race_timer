package results

import (
	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/events"
	"time"
)

type ScoredResult struct {
	Athlete competitors.Competitor
	Place   int
	Time    time.Duration
}

type Scorer interface {
	Score(inputEvents events.EventSource, athletes competitors.CompetitorLookup) ([]ScoredResult, error)
}

func NewOverallScoring() Scorer {
	return &overallScorer{}
}

type overallScorer struct {
	eventSource events.EventSource
	athletes    competitors.CompetitorLookup
}

func (os *overallScorer) Score(inputEvents events.EventSource, athletes competitors.CompetitorLookup) ([]ScoredResult, error) {
	start := []events.StartEvent{}
	finishes := map[int]events.FinishEvent{} //key is bib number
	places := map[int]events.PlaceEvent{}    //key is place
	result := []ScoredResult{}

	event, err := inputEvents.GetRaceEvent()

	for event != nil && err == nil {
		switch event.GetType() {
		case events.StartEventType:
			start = append(start, event.(events.StartEvent))
		case events.FinishEventType:
			finishes[event.(events.FinishEvent).GetBib()] = event.(events.FinishEvent)
		case events.PlaceEventType:
			places[event.(events.PlaceEvent).GetPlace()] = event.(events.PlaceEvent)
		}

		event, err = inputEvents.GetRaceEvent()
	}

	for place, event := range places {
		result = append(result, ScoredResult{athletes[event.GetBib()], place, finishes[event.GetBib()].GetFinishTime().Sub(start[0].GetStartTime())})
	}

	return result, nil
}
