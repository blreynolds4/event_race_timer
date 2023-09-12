package command

import (
	"blreynolds4/event-race-timer/events"
	"blreynolds4/event-race-timer/eventstream"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	redis "github.com/redis/go-redis/v9"
)

// unique name for a client
var clientName string

type Command interface {
	Run(args []string) (bool, error)
}

func init() {
	var err error
	clientName, err = os.Hostname()
	if err != nil {
		// use a random number
		clientName = fmt.Sprintf("race-cli-%d", rand.Intn(100))
	}
}

type noStateCommand struct {
	CmdFunc func(args []string) (bool, error)
}

// supported commands
//
//	ping to make sure server is ok
//	quit/stop/exit the progeam
//	start event with or without seed time to start at
//	finish event with or without bib
func (nsc *noStateCommand) Run(args []string) (bool, error) {
	return nsc.CmdFunc(args)
}

func NewQuitCommand() Command {
	return &noStateCommand{
		CmdFunc: func(args []string) (bool, error) {
			fmt.Println("quitting...")
			return true, nil
		},
	}
}

func NewPingCommand(rdb *redis.Client) Command {
	return &noStateCommand{
		CmdFunc: func(args []string) (bool, error) {
			cmdResult := rdb.Ping(context.TODO())
			fmt.Println(cmdResult.String())
			return false, cmdResult.Err()
		},
	}
}

func NewStartCommand(eventTarget events.EventTarget) Command {
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

			return false, eventTarget.SendRaceEvent(context.TODO(), eventstream.NewStartEvent(clientName, startTime.Add(-seedDuration)))
		},
	}

}

func NewFinishCommand(eventTarget events.EventTarget) Command {
	return &noStateCommand{
		CmdFunc: func(args []string) (bool, error) {
			var err error
			bib := events.NoBib
			if len(args) > 0 && len(args[0]) > 0 {
				bib, err = strconv.Atoi(args[0])
				if err != nil {
					fmt.Println("Finish Error (no event sent)", err, len(args), args)
					return false, err
				}
			}

			return false, eventTarget.SendRaceEvent(context.TODO(), eventstream.NewFinishEvent(clientName, time.Now().UTC(), bib))
		},
	}
}

func NewPlaceCommand(eventTarget events.EventTarget) Command {
	return &noStateCommand{
		CmdFunc: func(args []string) (bool, error) {
			var err error
			bib, place := events.NoBib, 0
			if len(args) > 1 {
				bib, err = strconv.Atoi(args[0])
				if err != nil {
					return false, err
				}

				place, err = strconv.Atoi(args[1])
				if err != nil {
					return false, err
				}

				return false, eventTarget.SendRaceEvent(context.TODO(), eventstream.NewPlaceEvent(clientName, bib, place))
			}

			return false, fmt.Errorf("missing bib or place argument")
		},
	}
}

func NewListFinishCommand(eventSource events.EventSource) Command {
	return &noStateCommand{
		CmdFunc: func(args []string) (bool, error) {
			var err error
			var startEvent events.StartEvent
			finishes := make([]events.FinishEvent, 0, 100)
			// read all the events and print them out
			var current events.RaceEvent
			current, err = eventSource.GetRaceEvent(context.TODO(), time.Second)
			if err != nil {
				return false, err
			}

			for current != nil {
				switch current.GetType() {
				case events.StartEventType:
					startEvent = current.(events.StartEvent)
				case events.FinishEventType:
					finishes = append(finishes, current.(events.FinishEvent))
				default:
				}

				current, err = eventSource.GetRaceEvent(context.TODO(), time.Second)
				if err != nil {
					return false, err
				}
			}

			// print the finish events in order with a duration base on the start event
			// can't print finishes with out a start event
			if startEvent != nil {
				fmt.Printf("%20s %20s %6s\n", "Event ID", "Time", "Bib")
				for _, fe := range finishes {
					fmt.Printf("%20s %20s %6d\n", fe.GetID(), fe.GetFinishTime().Sub(startEvent.GetStartTime()), fe.GetBib())
				}
			}

			return false, err
		},
	}
}

func NewAddBibCommand(eventSource events.EventSource, eventTarget events.EventTarget) Command {
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

			eventRange, err := eventSource.GetRaceEventRange(context.TODO(), args[0], args[0])
			if err != nil {
				return false, err
			}
			if len(eventRange) != 1 {
				return false, fmt.Errorf("event id did not return 1 event")
			}

			if len(eventRange) == 1 {
				finishEvent, ok := eventRange[0].(events.FinishEvent)
				if !ok || finishEvent.GetType() != events.FinishEventType {
					return false, fmt.Errorf("expected event id to be for finish event, skipping")
				}

				// create updated event with new bib
				updated := eventstream.NewFinishEvent(finishEvent.GetSource(), finishEvent.GetFinishTime(), bib)
				eventTarget.SendRaceEvent(context.TODO(), updated)
			}

			return false, nil
		},
	}
}
