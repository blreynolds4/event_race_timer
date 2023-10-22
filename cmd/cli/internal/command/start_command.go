package command

import (
	"blreynolds4/event-race-timer/internal/raceevents"
	"context"
	"time"
)

func NewStartCommand(sourceName string, eventTarget *raceevents.EventStream) Command {
	return &noStateCommand{
		CmdFunc: func(args []string) (bool, error) {
			startTime := time.Now().UTC()
			seedTime := "0s"
			if len(args) > 0 {
				seedTime = args[0]
			}
			seedDuration, err := time.ParseDuration(seedTime)
			if err != nil {
				return false, err
			}

			return false, eventTarget.SendStartEvent(context.TODO(), raceevents.StartEvent{
				Source:    sourceName,
				StartTime: startTime.Add(-seedDuration),
			})
		},
	}
}
