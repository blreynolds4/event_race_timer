package command

import (
	"blreynolds4/event-race-timer/cmd/cli/internal/repl"
	"blreynolds4/event-race-timer/internal/raceevents"
	"context"
	"fmt"
	"os"
	"strconv"
)

type placeRangeCommand struct {
	source      string
	eventTarget *raceevents.EventStream
}

func NewPlaceRangeCommand(sourceName string, eventTarget *raceevents.EventStream) Command {
	return &placeRangeCommand{
		source:      sourceName,
		eventTarget: eventTarget,
	}
}

func (pr *placeRangeCommand) String() string {
	return "PlaceRangeCommand"
}

func (pr *placeRangeCommand) Run(args []string) (bool, error) {
	// create a repl with place range command and quit command
	var err error
	nextPlace := 1
	if len(args) > 0 {
		nextPlace, err = strconv.Atoi(args[0])
		if err != nil {
			return false, err
		}
	}

	placeCmd := NewPlaceRangePlaceCommand(nextPlace, pr.source, pr.eventTarget)

	cmdRun := func(args []string) bool {
		if len(args) > 0 {
			if args[0] == "q" || args[0] == "quit" {
				return true
			}

			done, err := placeCmd.Run(args)
			if err != nil {
				fmt.Println("error in place command", err)
			}
			return done
		}
		return false
	}

	fmt.Println("next place is", nextPlace)
	repl := repl.NewReadEvalPrintLoop("places", os.Stdin, cmdRun)
	repl.Run()

	return false, nil
}

func NewPlaceRangePlaceCommand(nextPlace int, source string, eventTarget *raceevents.EventStream) Command {
	return &placeRangePlaceCommand{
		nextPlace:   nextPlace - 1, // gets incremented before use
		source:      source,
		eventTarget: eventTarget,
	}
}

type placeRangePlaceCommand struct {
	nextPlace   int
	source      string
	eventTarget *raceevents.EventStream
}

func (prc *placeRangePlaceCommand) Run(args []string) (bool, error) {
	var err error
	bib := raceevents.NoBib
	if len(args) > 0 {
		bib, err = strconv.Atoi(args[0])
		if err != nil {
			return false, err
		}

		prc.nextPlace++
		err := prc.eventTarget.SendPlaceEvent(context.TODO(), raceevents.PlaceEvent{
			Source: prc.source,
			Bib:    bib,
			Place:  prc.nextPlace,
		})
		if err != nil {
			return false, err
		}

		fmt.Println("sent place", prc.nextPlace, "for bib", bib)

		return false, nil
	}

	return false, fmt.Errorf("missing place argument")
}
