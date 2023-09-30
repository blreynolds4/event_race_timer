package places

import (
	"blreynolds4/event-race-timer/events"
	"blreynolds4/event-race-timer/eventstream"
	"context"
	"sort"
)

// Write a placer that takes event source and event target
type PlaceGenerator interface {
	GeneratePlaces(map[string]int) error
}

type defaultPlaceGenerator struct {
	eventSource events.EventSource
	eventTarget events.EventTarget
}

func NewPlaceGenerator(src events.EventSource, target events.EventTarget) PlaceGenerator {
	return &defaultPlaceGenerator{
		src,
		target,
	}
}

func (dpg *defaultPlaceGenerator) GeneratePlaces(sourceRanks map[string]int) error {
	// cache of finishes with bibs
	// Sorted by finish time, soonest to latest
	finishCache := make(map[int]events.FinishEvent)

	// read from the source any finish events with bibs (default_placer consumer group)
	event, err := dpg.eventSource.GetRaceEvent(context.TODO(), 0)
	if err != nil {
		return err
	}

	for event != nil {
		if event.GetType() == events.FinishEventType {
			// cast to FinishEvent, add to the cache
			finish := event.(events.FinishEvent)
			if finish.GetBib() != events.NoBib {
				// only cache finishes with bibs for placement
				// where the new finish is from a better source
				if dpg.isBetterFinish(finish, finishCache, sourceRanks) {
					finishCache[finish.GetBib()] = finish
					// create a slice of bibs in finish order
					finishedBibs := make([]int, 0, len(finishCache))
					for k := range finishCache {
						finishedBibs = append(finishedBibs, k)
					}
					sort.SliceStable(finishedBibs, func(i, j int) bool {
						return finishCache[finishedBibs[i]].GetFinishTime().Before(finishCache[finishedBibs[j]].GetFinishTime())
					})

					// loop through the slice and send events for the current bib
					// and everything after it
					send := false
					for i := 0; i < len(finishedBibs); i++ {
						current := finishCache[finishedBibs[i]]
						// don't send events till we get to the bib
						// from last finish
						if finish.GetBib() == current.GetBib() {
							send = true
						}

						if send {
							// create place event
							placeEvent := eventstream.NewPlaceEvent("default-placer", current.GetBib(), i+1)
							dpg.eventTarget.SendRaceEvent(context.TODO(), placeEvent)
						}
					}
				}
			}
		}
		event, err = dpg.eventSource.GetRaceEvent(context.TODO(), 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dpg *defaultPlaceGenerator) isBetterFinish(finish events.FinishEvent, finishCache map[int]events.FinishEvent, sourceRanks map[string]int) bool {
	previous, exists := finishCache[finish.GetBib()]
	if !exists {
		return true
	}

	// return true if new source is better then previous
	if sourceRanks[finish.GetSource()] < sourceRanks[previous.GetSource()] {
		return true
	}

	// source isn't better
	return false
}
