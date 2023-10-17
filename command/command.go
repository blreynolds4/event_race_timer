package command

import (
	"blreynolds4/event-race-timer/raceevents"
	"context"
	"fmt"
	"strconv"
	"time"

	redis "github.com/redis/go-redis/v9"
)

// unique name for a client
// var clientName string

type Command interface {
	Run(args []string) (bool, error)
}

// func init() {
// 	var err error
// 	clientName, err = os.Hostname()
// 	if err != nil {
// 		// use a random number
// 		clientName = fmt.Sprintf("race-cli-%d", rand.Intn(100))
// 	}
// }

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

func NewFinishCommand(sourceName string, eventTarget *raceevents.EventStream) Command {
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

			return false, eventTarget.SendFinishEvent(context.TODO(), raceevents.FinishEvent{
				Source:     sourceName,
				FinishTime: time.Now().UTC(),
				Bib:        bib,
			})
		},
	}
}

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
