package places

import (
	"blreynolds4/event-race-timer/events"
	"container/list"
	"fmt"
)

// Write a placer that takes event source and event target
type PlaceGenerator interface {
	GeneratePlaces() error
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

func (dpg *defaultPlaceGenerator) GeneratePlaces() error {
	// cache of finishes with bibs
	// Sorted by finish time, soonest to latest
	finishCache := list.New()

	// read from the source any finish events with bibs (default_placer consumer group)
	event, err := dpg.eventSource.GetRaceEvent()
	fmt.Println("GOT EVENT", event)
	if err != nil {
		return err
	}

	for event != nil {
		fmt.Println("Processing event ", event.GetType())
		if event.GetType() == events.FinishEventType {
			// cast to FinishEvent, add to the cache
			finish := event.(events.FinishEvent)
			if finish.GetBib() != events.NoBib {
				fmt.Println("Has bib")
				// only cache finishes with bibs for placement
				// sort them by finish time earliest to latest
				var inserted *list.Element
				currentPlace := 1
				if finishCache.Len() == 0 {
					fmt.Println("Empty list, pushing to front")
					inserted = finishCache.PushFront(finish)
				} else {
					fmt.Println("Non Empty list, inserting...")
					for e := finishCache.Front(); e != nil; e = e.Next() {
						current := e.Value.(events.FinishEvent)
						fmt.Println("Finish time ", finish.GetFinishTime(), "current list time", current.GetFinishTime())
						if finish.GetFinishTime().Before(current.GetFinishTime()) {
							//insert in front of current and stop
							inserted = finishCache.InsertBefore(finish, e)
							fmt.Println("Inserted", inserted)
							break
						} else {
							fmt.Println("event is after")
							currentPlace = currentPlace + 1
						}
					}

					if inserted == nil {
						fmt.Println("Add to the end")
						inserted = finishCache.PushBack(finish)
					}
				}

				// send Place events for the new event and everything after it in the cache order
				fmt.Println("INSERTED", inserted.Value)
				for e := inserted; e != nil; e = e.Next() {
					// create place event
					current := e.Value.(events.FinishEvent)
					placeEvent := events.NewPlaceEvent("default_placer", current.GetBib(), currentPlace)
					currentPlace = currentPlace + 1
					fmt.Println("SENDING PLACE", placeEvent)
					dpg.eventTarget.SendRaceEvent(placeEvent)
				}
			}
		}
		event, err = dpg.eventSource.GetRaceEvent()
		if err != nil {
			return err
		}
	}

	return nil
}
