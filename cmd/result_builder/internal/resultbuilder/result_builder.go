package resultbuilder

import (
	"blreynolds4/event-race-timer/internal/competitors"
	"blreynolds4/event-race-timer/internal/raceevents"
	"blreynolds4/event-race-timer/internal/results"
	"context"
	"log/slog"
	"time"
)

type ResultBuilder interface {
	BuildResults(inputEvents raceevents.EventStream, athletes competitors.CompetitorLookup, results results.ResultStream, ranking map[string]int) error
}

func NewResultBuilder(l *slog.Logger) ResultBuilder {
	return &resultBuilder{
		logger: l.With("app", "result-builder"),
	}
}

type resultBuilder struct {
	logger *slog.Logger
}

func (rb *resultBuilder) BuildResults(inputEvents raceevents.EventStream,
	athletes competitors.CompetitorLookup,
	resultOutput results.ResultStream,
	ranking map[string]int) error {

	start := []raceevents.StartEvent{}      //array to store all of the start events
	rr := make(map[int]*results.RaceResult) //map of race results, bib number is key
	ft := make(map[int]time.Time)           // map of times with bib number being key
	placeIndex := make(map[int]*results.RaceResult)

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
					rb.sendResult(context.TODO(), rr[result.Bib], resultOutput)
				}
			}
		case raceevents.FinishEvent:
			fe := event.Data.(raceevents.FinishEvent)

			// only handle bibs for athletes that exist
			if _, bibFound := athletes[fe.Bib]; bibFound {
				result := rr[fe.Bib]
				if result == nil {
					// the result doesn't exist in the cache
					result = new(results.RaceResult)
					result.Bib = fe.Bib
					result.Athlete = athletes[fe.Bib]
					rr[fe.Bib] = result
				}

				//if the ranking of the new event source is higher than the old create a new result
				if ranking[fe.Source] <= ranking[result.FinishSource] || ranking[result.FinishSource] == 0 {

					result.FinishSource = fe.Source
					if len(start) > 0 {
						latest_start := len(start) - 1
						result.Time = fe.FinishTime.Sub(start[latest_start].StartTime)
					} else {
						// no start event yet, just save the finish
						ft[fe.Bib] = fe.FinishTime
					}
					rr[fe.Bib] = result

					if rr[fe.Bib].IsComplete() {
						rb.sendResult(context.TODO(), rr[result.Bib], resultOutput)
					}
				}
			}
		case raceevents.PlaceEvent:
			pe := event.Data.(raceevents.PlaceEvent)
			if _, bibFound := athletes[pe.Bib]; bibFound {
				// see if a result exists for this place
				// get the result for the bib
				bibResult := rr[pe.Bib]
				if bibResult == nil {
					// this is a new result
					bibResult = new(results.RaceResult)
					bibResult.Bib = pe.Bib
					bibResult.Athlete = athletes[pe.Bib]
					bibResult.PlaceSource = pe.Source
					rr[pe.Bib] = bibResult
				}

				previousPlace := bibResult.Place
				if ranking[pe.Source] <= ranking[bibResult.PlaceSource] || ranking[bibResult.PlaceSource] == 0 {
					// send updated results for the new place and everything after
					switch {
					case previousPlace == 0:
						// add and send result for new place
						addPlaceResult(bibResult, pe, placeIndex)
						rb.sendResult(context.TODO(), bibResult, resultOutput)
					case previousPlace < pe.Place:
						// require promotions pe.Place < previousPlace
						demotePlaceResult(bibResult, pe, placeIndex)
						// Loop and send from previous place to pe place inclusive
						for i := previousPlace; i <= pe.Place; i++ {
							if place, exists := placeIndex[i]; exists {
								// if place >= pe.Place {
								rb.sendResult(context.TODO(), place, resultOutput)
							}
						}
					default:
						// add and update
						addPlaceResult(bibResult, pe, placeIndex)
						// send updated results from new place to old place inclusive
						for i := pe.Place; i <= previousPlace; i++ {
							if place, exists := placeIndex[i]; exists {
								rb.sendResult(context.TODO(), place, resultOutput)
							}
						}
					}
				}
			} else {
				rb.logger.Info("skipping unknown bib", "bib", pe.Bib)
			}
		}

		gotEvent, err = inputEvents.GetRaceEvent(context.TODO(), 0, &event)
		if err != nil {
			return err
		}
	}
	return nil
}

func addPlaceResult(rr *results.RaceResult, pe raceevents.PlaceEvent, places map[int]*results.RaceResult) {
	// if the there is a result in the new place already, make room and save the result
	delete(places, rr.Place)
	_, found := places[pe.Place]
	if found {
		// move every result back a place (without updating the place source)
		// place currently in result is deleted so start one less
		// assumes new place is better (less than existing)
		for i := rr.Place - 1; i >= pe.Place; i-- {
			// there could be place gaps
			if _, exists := places[i]; exists {
				places[i+1] = places[i]
				places[i+1].Place = i + 1
			}
		}
	}

	// put the new place in
	rr.Place = pe.Place
	rr.PlaceSource = pe.Source
	places[pe.Place] = rr
}

func demotePlaceResult(rr *results.RaceResult, pe raceevents.PlaceEvent, places map[int]*results.RaceResult) {
	// if the there is a result in the new place already, make room and save the result
	delete(places, rr.Place)
	_, found := places[pe.Place]
	if found {
		// move every result up a place (without updating the place source)
		// place currently in result is deleted so start one more
		// assumes new place is worse (more than existing)
		for i := rr.Place + 1; i <= pe.Place; i++ {
			// there could be place gaps
			if _, exists := places[i]; exists {
				places[i-1] = places[i]
				places[i-1].Place = i - 1
			}
		}
	}

	// put the new place in
	rr.Place = pe.Place
	rr.PlaceSource = pe.Source
	places[pe.Place] = rr
}

func (rb *resultBuilder) sendResult(ctx context.Context, rr *results.RaceResult, s results.ResultStream) {
	copy := *rr

	s.SendResult(ctx, copy)
	rb.logger.Info("result sent", "bib", rr.Bib, "place", rr.Place, "elapsedTime", rr.Time.String())
}
