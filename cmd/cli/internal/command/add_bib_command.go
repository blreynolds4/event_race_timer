package command

import (
	"blreynolds4/event-race-timer/raceevents"
	"context"
	"fmt"
	"strconv"
)

func NewAddBibCommand(eventStream *raceevents.EventStream) Command {
	return &noStateCommand{
		CmdFunc: func(args []string) (bool, error) {
			//get the event with the event id and resend it with a bib attached
			if len(args) < 2 {
				return false, fmt.Errorf("add bib requires to arguments:  <finish event id> <bib>")
			}

			bib, err := strconv.Atoi(args[1])
			if err != nil {
				return false, err
			}

			msgBuffer := make([]raceevents.Event, 5)
			countRead, err := eventStream.GetRaceEventRange(context.TODO(), args[0], args[0], msgBuffer)
			if err != nil {
				return false, err
			}
			if countRead != 1 {
				return false, fmt.Errorf("event id did not return 1 event")
			}

			finishEvent, ok := msgBuffer[0].Data.(raceevents.FinishEvent)
			if !ok {
				return false, fmt.Errorf("expected event id to be for finish event, skipping")
			}

			// create updated event with new bib
			eventStream.SendFinishEvent(context.TODO(), raceevents.FinishEvent{
				Source:     finishEvent.Source,
				Bib:        bib,
				FinishTime: finishEvent.FinishTime,
			})

			return false, nil
		},
	}
}
