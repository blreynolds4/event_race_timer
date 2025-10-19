package cli

import (
	"blreynolds4/event-race-timer/cmd/cli/internal/command"
	"blreynolds4/event-race-timer/cmd/cli/internal/repl"
	"blreynolds4/event-race-timer/internal/meets"
	"blreynolds4/event-race-timer/internal/raceevents"
	"blreynolds4/event-race-timer/internal/redis_stream"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

const sourceName = "manual"

type CliApp struct {
	replCommands map[string]command.Command
}

func NewCliApp() CliApp {
	return CliApp{
		replCommands: make(map[string]command.Command),
	}
}

func (ca CliApp) Run(claDbAddress string, claDbNumber int, claRacename string, pgConnect string) {
	// connect to redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     claDbAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	defer rdb.Close()

	rawStream := redis_stream.NewRedisStream(rdb, claRacename)
	eventStream := raceevents.NewEventStream(rawStream)

	// create the command map
	ca.createCommandMap(rdb, claRacename, eventStream, pgConnect)

	inputRepl := repl.NewReadEvalPrintLoop(fmt.Sprintf("race-cli:%s", claRacename), os.Stdin, ca.commandRunner)
	inputRepl.Run()
}

func (ca CliApp) createCommandMap(rdb *redis.Client, raceName string, eventStream raceevents.EventStream, pgConnect string) {
	athleteReader, err := meets.NewAthleteReader(pgConnect)
	if err != nil {
		fmt.Println("error creating athlete reader:", err.Error())
		os.Exit(1)
	}

	raceReader, err := meets.NewRaceReader(pgConnect)
	if err != nil {
		fmt.Println("error creating race reader:", err.Error())
		os.Exit(1)
	}

	raceWriter, err := meets.NewRaceWriter(pgConnect)
	if err != nil {
		fmt.Println("error creating race writer:", err.Error())
		os.Exit(1)
	}

	ca.replCommands["quit"] = command.NewQuitCommand()
	ca.replCommands["q"] = ca.replCommands["quit"]
	ca.replCommands["exit"] = ca.replCommands["quit"]
	ca.replCommands["stop"] = ca.replCommands["quit"]

	ca.replCommands["ping"] = command.NewPingCommand(rdb)

	ca.replCommands["start"] = command.NewStartCommand(sourceName, eventStream)
	ca.replCommands["s"] = ca.replCommands["start"]

	ca.replCommands["place"] = command.NewPlaceCommand(sourceName, eventStream)
	ca.replCommands["p"] = ca.replCommands["place"]

	ca.replCommands["placeRange"] = command.NewPlaceRangeCommand(sourceName, eventStream)
	ca.replCommands["pr"] = ca.replCommands["placeRange"]

	ca.replCommands["list"] = command.NewListFinishCommand(eventStream)

	ca.replCommands["bib"] = command.NewAddBibCommand(eventStream)

	ca.replCommands["removeAthleteFromRace"] = command.NewDeleteAthleteFromRaceCommand(athleteReader, raceReader, raceWriter)
	ca.replCommands["rar"] = ca.replCommands["removeAthleteFromRace"]

	ca.replCommands["addAthleteToRace"] = command.NewAddAthleteToRaceCommand(athleteReader, raceReader, raceWriter)
	ca.replCommands["aar"] = ca.replCommands["addAthleteToRace"]

	ca.replCommands["finish"] = command.NewFinishCommand(sourceName, eventStream)
	ca.replCommands["f"] = ca.replCommands["finish"]
}

func (ca CliApp) commandRunner(args []string) bool {
	if len(args) > 0 {
		cmd := args[0]
		cmdFunc, found := ca.replCommands[cmd]
		if found {
			done, err := cmdFunc.Run(args[1:])
			if err != nil {
				fmt.Println("error during", cmd, ":", err.Error())
			}
			return done
		} else {
			// default to a finish event and send cmd as the bib
			fmt.Println("unknown command", cmd, "skipping")
		}
	}

	return false
}
