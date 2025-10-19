package command

import (
	"blreynolds4/event-race-timer/internal/raceevents"
	"context"
	"fmt"
	"strconv"
	"time"
)

func NewFinishCommand(sourceName string, eventTarget raceevents.EventStream) Command {
	return &noStateCommand{
		CmdFunc: func(args []string) (bool, error) {
			var err error
			bib := raceevents.NoBib
			if len(args) > 0 && len(args[0]) > 0 {
				bib, err = strconv.Atoi(args[0])
				if err != nil {
					fmt.Println("Finish Error (no event sent)", err, len(args), args)
					return false, err
				}
			}

			finish := time.Now().UTC()
			if len(args) > 1 && len(args[1]) > 0 {
				// parse a finish time in RFC3339Nano format
				finish, err = time.Parse(time.RFC3339Nano, args[1])
				if err != nil {
					return false, err
				}
			}

			return false, eventTarget.SendFinishEvent(context.TODO(), raceevents.FinishEvent{
				Source:     sourceName,
				FinishTime: finish,
				Bib:        bib,
			})
		},
	}
}
