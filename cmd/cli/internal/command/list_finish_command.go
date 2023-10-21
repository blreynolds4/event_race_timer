package command

import (
	"blreynolds4/event-race-timer/raceevents"
	"context"
	"fmt"
	"time"
)

func NewListFinishCommand(eventSource *raceevents.EventStream) Command {
	return &noStateCommand{
		CmdFunc: func(args []string) (bool, error) {
			var err error
			var startEvent raceevents.StartEvent
			finishes := make([]raceevents.Event, 0, 100)
			hasStart := false
			// read all the events and print them out
			var current raceevents.Event
			readEvent, err := eventSource.GetRaceEvent(context.TODO(), time.Second, &current)
			if err != nil {
				return false, err
			}

			for readEvent {
				switch current.Data.(type) {
				case raceevents.StartEvent:
					startEvent = current.Data.(raceevents.StartEvent)
					hasStart = true
				case raceevents.FinishEvent:
					finishes = append(finishes, current)
				default:
				}

				readEvent, err = eventSource.GetRaceEvent(context.TODO(), time.Second, &current)
				if err != nil {
					return false, err
				}
			}

			// print the finish events in order with a duration base on the start event
			// can't print finishes with out a start event
			if hasStart {
				fmt.Printf("%20s %20s %6s\n", "Event ID", "Time", "Bib")
				for _, e := range finishes {
					fe := e.Data.(raceevents.FinishEvent)
					fmt.Printf("%20s %20s %6d\n", e.ID, fe.FinishTime.Sub(startEvent.StartTime), fe.Bib)
				}
			}

			return false, err
		},
	}
}
