package places

import (
	"blreynolds4/event-race-timer/raceevents"
	"context"
	"sort"
)

// Write a placer that takes event source and event target
type PlaceGenerator interface {
	GeneratePlaces(map[string]int) error
}

type defaultPlaceGenerator struct {
	stream *raceevents.EventStream
}

func NewPlaceGenerator(es *raceevents.EventStream) PlaceGenerator {
	return &defaultPlaceGenerator{
		stream: es,
	}
}

func (dpg *defaultPlaceGenerator) GeneratePlaces(sourceRanks map[string]int) error {
	// cache of finishes with bibs
	// Sorted by finish time, soonest to latest
	finishCache := make(map[int]raceevents.FinishEvent)

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
			if finish.Bib != raceevents.NoBib {
				// only cache finishes with bibs for placement
				// where the new finish is from a better source
				if dpg.isBetterFinish(finish, finishCache, sourceRanks) {
					finishCache[finish.Bib] = finish
					// create a slice of bibs in finish order
					finishedBibs := make([]int, 0, len(finishCache))
					for k := range finishCache {
						finishedBibs = append(finishedBibs, k)
					}
					sort.SliceStable(finishedBibs, func(i, j int) bool {
						return finishCache[finishedBibs[i]].FinishTime.Before(finishCache[finishedBibs[j]].FinishTime)
					})

					// loop through the slice and send events for the current bib
					// and everything after it
					send := false
					for i := 0; i < len(finishedBibs); i++ {
						current := finishCache[finishedBibs[i]]
						// don't send events till we get to the bib
						// from last finish
						if finish.Bib == current.Bib {
							send = true
						}

						if send {
							// send place event
							dpg.stream.SendPlaceEvent(context.TODO(), raceevents.PlaceEvent{
								Source: "default-placer",
								Place:  i + 1,
								Bib:    current.Bib,
							})
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

func (dpg *defaultPlaceGenerator) isBetterFinish(finish raceevents.FinishEvent, finishCache map[int]raceevents.FinishEvent, sourceRanks map[string]int) bool {
	previous, exists := finishCache[finish.Bib]
	if !exists {
		return true
	}

	// return true if new source is better then previous
	if sourceRanks[finish.Source] < sourceRanks[previous.Source] {
		return true
	}

	// source isn't better
	return false
}
