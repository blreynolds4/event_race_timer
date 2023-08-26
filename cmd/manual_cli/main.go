package main

import (
	"blreynolds4/event-race-timer/command"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	redis "github.com/go-redis/redis/v7"
)

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

	// create the command map
	loopCommands := make(map[string]command.Command)
	finishCommand := command.NewFinishCommand(rdb, claRacename)
	loopCommands["quit"] = command.NewQuitCommand()
	loopCommands["q"] = command.NewQuitCommand()
	loopCommands["exit"] = command.NewQuitCommand()
	loopCommands["stop"] = command.NewQuitCommand()
	loopCommands["ping"] = command.NewPingCommand(rdb)
	loopCommands["start"] = command.NewStartCommand(rdb, claRacename)
	loopCommands["s"] = command.NewStartCommand(rdb, claRacename)
	loopCommands["place"] = command.NewPlaceCommand(rdb, claRacename)
	loopCommands["p"] = command.NewPlaceCommand(rdb, claRacename)
	loopCommands["list"] = command.NewListFinishCommand(rdb, claRacename)
	loopCommands["finish"] = finishCommand
	loopCommands["f"] = finishCommand

	scanner := bufio.NewScanner(os.Stdin)
	done := false
	var err error
	for !done {
		// read a line of input into an array of strings
		fmt.Printf("race-cmd:%s>", claRacename)
		if scanner.Scan() {
			fmt.Println()
			line := scanner.Text()
			// look up the first string as the command and pass the rest to the command if one is found.
			cmdArgs := strings.Split(line, " ")
			if len(cmdArgs) > 0 {
				cmd := cmdArgs[0]
				cmdFunc, found := loopCommands[cmd]
				if found {
					done, err = cmdFunc.Run(cmdArgs[1:])
					if err != nil {
						fmt.Println("error during", cmd, ":", err.Error())
					}
				} else {
					// default to a finish event and send cmd as the bib
					fmt.Println("defaulting to finish for", cmdArgs)
					finishCommand.Run(cmdArgs)
				}
			}
		}
	}
}
