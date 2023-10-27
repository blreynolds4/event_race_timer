package places

import (
	"blreynolds4/event-race-timer/internal/competitors"
	"blreynolds4/event-race-timer/internal/raceevents"
	"context"
	"fmt"
	"sort"
)

const placerSourceName = "default-placer"

// Write a placer that takes event source and event target
type PlaceGenerator interface {
	GeneratePlaces(competitors.CompetitorLookup, map[string]int) error
}

type defaultPlaceGenerator struct {
	stream raceevents.EventStream
}

func NewPlaceGenerator(es raceevents.EventStream) PlaceGenerator {
	return &defaultPlaceGenerator{
		stream: es,
	}
}

// preserve order of arrival of the bibs

func (dpg *defaultPlaceGenerator) GeneratePlaces(athletes competitors.CompetitorLookup, sourceRanks map[string]int) error {
	// cache of finishes with bibs
	finishCache := make(map[int]raceevents.FinishEvent)
	// start sorting with bibs in arrival order so the sort can use arrival to break ties
	finishedBibs := make([]int, 0)

	// read from the source any finish events with bibs (default_placer consumer group)
	var event raceevents.Event
	gotEvent, err := dpg.stream.GetRaceEvent(context.TODO(), 0, &event)
	if err != nil {
		return err
	}

	for gotEvent {
		switch event.Data.(type) {
		case raceevents.FinishEvent:
			finish := event.Data.(raceevents.FinishEvent)
			_, bibFound := athletes[finish.Bib]
			if bibFound {
				previous, existed := finishCache[finish.Bib]
				if !existed {
					finishedBibs = append(finishedBibs, finish.Bib)
					finishCache[finish.Bib] = finish
				}

				// only cache finishes with bibs of known athletes for placement
				// where the new finish is from a better source
				if sourceRanks[finish.Source] < sourceRanks[previous.Source] || !existed {
					finishCache[finish.Bib] = finish
					// create a slice of bibs in finish order
					sorted := make([]int, len(finishedBibs))
					copy(sorted, finishedBibs)
					sort.SliceStable(sorted, func(i, j int) bool {
						return finishCache[sorted[i]].FinishTime.Before(finishCache[sorted[j]].FinishTime)
					})

					// loop through the slice and send events for the current bib
					// and everything after it
					send := false
					for i := 0; i < len(sorted); i++ {
						current := finishCache[sorted[i]]
						// don't send events till we get to the bib
						// from last finish
						if finish.Bib == current.Bib {
							send = true
						}

						if send {
							// send place event
							dpg.stream.SendPlaceEvent(context.TODO(), raceevents.PlaceEvent{
								Source: placerSourceName,
								Place:  i + 1,
								Bib:    current.Bib,
							})
							fmt.Printf("Place sent for bib %d %d\n", current.Bib, i+1)
						}
					}
				}
			}
		}

		gotEvent, err = dpg.stream.GetRaceEvent(context.TODO(), 0, &event)
		if err != nil {
			return err
		}
	}

	return nil
}
