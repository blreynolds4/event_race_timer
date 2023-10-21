package main

import (
	"blreynolds4/event-race-timer/cmd/cli/internal/command"
	"blreynolds4/event-race-timer/cmd/cli/internal/repl"
	"blreynolds4/event-race-timer/raceevents"
	"blreynolds4/event-race-timer/redis_stream"
	"flag"
	"fmt"
	"os"

	redis "github.com/redis/go-redis/v9"
)

const sourceName = "manual"

func main() {
	// connect to redis
	// cli for db address, username, password, db, stream name?
	// stream is specific to a race
	var claDbAddress string
	var claDbNumber int
	var claRacename string

	flag.StringVar(&claDbAddress, "dbAddress", "localhost:6379", "The host and port ie localhost:6379")
	flag.IntVar(&claDbNumber, "dbNumber", 0, "The database to use, defaults to 0")
	flag.StringVar(&claRacename, "raceName", "race", "The name of the race being timed (no spaces)")

	// parse command line
	flag.Parse()

	// connect to redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     claDbAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	defer rdb.Close()

	rawStream := redis_stream.NewRedisEventStream(rdb, claRacename)
	eventStream := raceevents.NewEventStream(rawStream)

	// create the command map
	replCommands := make(map[string]command.Command)

	replCommands["quit"] = command.NewQuitCommand()
	replCommands["q"] = replCommands["quit"]
	replCommands["exit"] = replCommands["quit"]
	replCommands["stop"] = replCommands["quit"]

	replCommands["ping"] = command.NewPingCommand(rdb)

	replCommands["start"] = command.NewStartCommand(sourceName, eventStream)
	replCommands["s"] = replCommands["start"]

	replCommands["place"] = command.NewPlaceCommand(sourceName, eventStream)
	replCommands["p"] = replCommands["place"]

	replCommands["placeRange"] = command.NewPlaceRangeCommand(sourceName, eventStream)
	replCommands["pr"] = replCommands["placeRange"]

	replCommands["list"] = command.NewListFinishCommand(eventStream)

	replCommands["bib"] = command.NewAddBibCommand(eventStream)

	replCommands["finish"] = command.NewFinishCommand(sourceName, eventStream)
	replCommands["f"] = replCommands["finish"]

	cmdRun := func(args []string) bool {
		if len(args) > 0 {
			cmd := args[0]
			cmdFunc, found := replCommands[cmd]
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

	inputRepl := repl.NewReadEvalPrintLoop(fmt.Sprintf("race-cli:%s", claRacename), os.Stdin, cmdRun)
	inputRepl.Run()
}
