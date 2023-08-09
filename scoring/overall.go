package scoring

import (
	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/events"
	"fmt"
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
	//start := []events.StartEvent{}
	//finishes := map[int]events.FinishEvent{} //key is bib number
	places := map[int]events.PlaceEvent{} //key is place
	result := []ScoredResult{}

	for event, err := inputEvents.GetRaceEvent(); event != nil && err == nil; {
		fmt.Println("read event", event)
		// switch event.GetType() {
		// case events.StartEventType:
		// 	start = append(start, event.(events.StartEvent))
		// case events.FinishEventType:
		// 	finishes[event.(events.FinishEvent).GetBib()] = event.(events.FinishEvent)
		// case events.PlaceEventType:
		// 	places[event.(events.PlaceEvent).GetPlace()] = event.(events.PlaceEvent)
		// }

	}

	//find first place
	//fill in info for a scored result and append to result
	//remove from arrays
	//repeat with next postion

	for place, event := range places {
		result = append(result, ScoredResult{athletes[event.GetBib()], place, time.Duration(1)})
	}

	return result, nil
}
