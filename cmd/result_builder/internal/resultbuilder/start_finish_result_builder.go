package resultbuilder

import (
	"blreynolds4/event-race-timer/internal/competitors"
	"blreynolds4/event-race-timer/internal/raceevents"
	"blreynolds4/event-race-timer/internal/results"
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"
)

func NewStartFinishResultBuilder(places string) ResultBuilder {
	// read in place file
	finishes := make(map[int]int)
	loadPlaces(places, finishes)
	return &startFinishResultBuilder{places: finishes}
}

func loadPlaces(fname string, places map[int]int) error {
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	place := 1
	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		bib, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return err
		}

		places[bib] = place
		place++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

type startFinishResultBuilder struct {
	// map of bib to place
	places map[int]int
}

func (rb *startFinishResultBuilder) BuildResults(inputEvents *raceevents.EventStream,
	athletes competitors.CompetitorLookup,
	outputResults *results.ResultStream,
	ranking map[string]int) error {

	start := make([]raceevents.StartEvent, 0) //array to store all of the start events
	rr := make(map[int]*results.RaceResult)   //map of race results, bib number is key
	ft := make(map[int]time.Time)             // map of times with bib number being key

	fmt.Println("START FINISH RESULT BUILDER IS FOR DEBUGGING")
	startCount := 0
	finishCount := 0
	resultSent := 0

	var event raceevents.Event
	gotEvent, err := inputEvents.GetRaceEvent(context.TODO(), 0, &event)
	if err != nil {
		return err
	}

	for gotEvent {
		switch event.Data.(type) {
		case raceevents.StartEvent:
			startCount++
			se := event.Data.(raceevents.StartEvent)
			start = append(start, se)
			fmt.Println("Got Start Event for", se.StartTime, "from", se.Source)

			// iterate over result to get times for all finish events that came before the start event
			for _, result := range rr {
				result.Time = ft[result.Bib].Sub(start[len(start)-1].StartTime)
				rr[result.Bib] = result
				resultSent++
				sendResult(context.TODO(), rr[result.Bib], outputResults)
			}
		case raceevents.FinishEvent:
			finishCount++
			fe := event.Data.(raceevents.FinishEvent)

			// only handle bibs for athletes that exist
			if athlete, bibFound := athletes[fe.Bib]; bibFound {
				fmt.Println("Got finish", fe.Bib, athlete.Name, athlete.Team)
				result := rr[fe.Bib]
				if result == nil {
					// the result doesn't exist in the cache
					result = new(results.RaceResult)
					result.Bib = fe.Bib
					result.Athlete = athlete
					rr[fe.Bib] = result
				}

				if place, found := rb.places[fe.Bib]; found {
					result.Place = place
					result.PlaceSource = "manual"
					// clear the place as used
					delete(rb.places, fe.Bib)
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
						resultSent++
						sendResult(context.TODO(), rr[result.Bib], outputResults)
					} else {
						fmt.Println("NOT COMPLETE", fe.Bib, result)
					}
				}
			} else {
				fmt.Println("BIB NOT FOUND", fe.Bib)
			}
		case raceevents.PlaceEvent:
		}

		gotEvent, err = inputEvents.GetRaceEvent(context.TODO(), 5, &event)
		if err != nil {
			fmt.Println("Bibs with results", len(rr), "athlete count", len(athletes))
			printMissingAthletes(rr, athletes)

			// send missed places
			for bib, place := range rb.places {
				fmt.Println("No finish for", bib, "chute place", place)
				result := new(results.RaceResult)
				result.Bib = bib
				result.Place = place
				result.PlaceSource = "manual"
				if athlete, bibFound := athletes[bib]; bibFound {
					result.Athlete = athlete
					sendResult(context.TODO(), result, outputResults)
					resultSent++
				} else {
					fmt.Println("chute bib", bib, "not in athletes")
				}
			}

			fmt.Println("Start count", startCount)
			fmt.Println("Finish count", finishCount)
			fmt.Println("Result Sent", resultSent)
			return err
		}
	}

	return nil
}

func printMissingAthletes(outputResults map[int]*results.RaceResult, athletes competitors.CompetitorLookup) {
	for bib, athlete := range athletes {
		if _, found := outputResults[bib]; !found {
			fmt.Println("No Result: ", bib, athlete.Name, athlete.Team)
		}
	}
}
