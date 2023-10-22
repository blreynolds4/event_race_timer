package command

import (
	"blreynolds4/event-race-timer/internal/raceevents"
	"context"
	"fmt"
	"strconv"
)

func NewPlaceCommand(sourceName string, eventTarget *raceevents.EventStream) Command {
	return &noStateCommand{
		CmdFunc: func(args []string) (bool, error) {
			var err error
			bib, place := raceevents.NoBib, 0
			if len(args) > 1 {
				bib, err = strconv.Atoi(args[0])
				if err != nil {
					return false, err
				}

				place, err = strconv.Atoi(args[1])
				if err != nil {
					return false, err
				}

				return false, eventTarget.SendPlaceEvent(context.TODO(), raceevents.PlaceEvent{
					Source: sourceName,
					Bib:    bib,
					Place:  place,
				})
			}

			return false, fmt.Errorf("missing bib or place argument")
		},
	}
}
