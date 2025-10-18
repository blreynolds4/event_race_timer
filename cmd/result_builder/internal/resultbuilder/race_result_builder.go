package resultbuilder

import (
	"blreynolds4/event-race-timer/internal/meets"
	"blreynolds4/event-race-timer/internal/raceevents"
	"context"
	"log/slog"
)

type RaceResultBuilder interface {
	BuildRaceResults(inputEvents raceevents.EventStream, athletes meets.AthleteLookup, ranking map[string]int, resultWriter meets.RaceResultWriter) error
}

func NewRaceResultBuilder(l *slog.Logger) RaceResultBuilder {
	return &raceResultBuilder{
		logger: l.With("app", "result-builder"),
	}
}

type raceResultBuilder struct {
	logger    *slog.Logger
	startTime *raceevents.StartEvent
}

func (rb *raceResultBuilder) hasStartTime() bool {
	return rb.startTime != nil
}

func (rb *raceResultBuilder) getStartTime() raceevents.StartEvent {
	if !rb.hasStartTime() {
		panic("no start time available")
	}
	return *rb.startTime
}

func (rb *raceResultBuilder) BuildRaceResults(inputEvents raceevents.EventStream,
	athletes meets.AthleteLookup,
	ranking map[string]int,
	resultWriter meets.RaceResultWriter) error {

	resultCache := make(map[int]*meets.RaceResult) //map of race results, bib number is key
	pendingFinishEvents := make(map[int]raceevents.FinishEvent)

	var event raceevents.Event
	gotEvent, err := inputEvents.GetRaceEvent(context.TODO(), 0, &event)
	if err != nil {
		return err
	}
	rb.logger.Info("GotEvent ", "event", gotEvent)

	for gotEvent {
		switch event.Data.(type) {
		case raceevents.StartEvent:
			se := event.Data.(raceevents.StartEvent)
			if !rb.hasStartTime() {
				rb.startTime = &se
			}

			// we have a start time now
			// update all pending finish
			for bib, pendingFinish := range pendingFinishEvents {
				resultCache[bib].FinishSource = pendingFinish.Source
				startTime := rb.getStartTime()
				resultCache[bib].Time = pendingFinish.FinishTime.Sub(startTime.StartTime)

				resultWriter.SaveResult(resultCache[bib])
				delete(pendingFinishEvents, bib)
			}

		case raceevents.FinishEvent:
			fe := event.Data.(raceevents.FinishEvent)

			// only handle bibs for athletes that exist
			if _, bibFound := athletes[fe.Bib]; bibFound {
				result := resultCache[fe.Bib]
				if result == nil {
					// the result doesn't exist in the cache
					result = new(meets.RaceResult)
					result.Bib = fe.Bib
					result.Athlete = athletes[fe.Bib]
					rb.logger.Info("New result created for bib", "bib", fe.Bib, "athlete", result.Athlete.DaID)
					resultCache[fe.Bib] = result
				}

				//if the ranking of the new event source is higher than the old create a new result
				if ranking[fe.Source] <= ranking[result.FinishSource] || ranking[result.FinishSource] == 0 {
					result.FinishSource = fe.Source
					if rb.hasStartTime() {
						startTime := rb.getStartTime()
						result.Time = fe.FinishTime.Sub(startTime.StartTime)
						rb.logger.Info("Result updated for bib", "bib", fe.Bib, "athlete", result.Athlete.LastName, "time", result.Time)
						resultWriter.SaveResult(resultCache[fe.Bib])
					} else {
						// save the whole finish event so we have the time and the source
						// information needed to build a result when a start time is available
						pendingFinishEvents[fe.Bib] = fe
					}
				}
			}
		case raceevents.PlaceEvent:
			pe := event.Data.(raceevents.PlaceEvent)
			if _, bibFound := athletes[pe.Bib]; bibFound {
				// see if a result exists for this place
				// get the result for the bib
				bibResult := resultCache[pe.Bib]
				if bibResult == nil {
					// this is a new result
					bibResult = new(meets.RaceResult)
					bibResult.Bib = pe.Bib
					bibResult.Athlete = athletes[pe.Bib]
					resultCache[pe.Bib] = bibResult
				}

				if ranking[pe.Source] <= ranking[bibResult.PlaceSource] || ranking[bibResult.PlaceSource] == 0 {
					bibResult.Place = pe.Place
					bibResult.PlaceSource = pe.Source
					resultWriter.SaveResult(resultCache[pe.Bib])
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
