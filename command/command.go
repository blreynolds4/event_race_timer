package command

import (
	"blreynolds4/event-race-timer/events"
	"fmt"
	"math/rand"
	"os"
	"time"

	redis "github.com/go-redis/redis/v7"
)

//unique name for a client
var clientName string

type CommandFunction func(args []string) (bool, error)

func init() {
	var err error
	clientName, err = os.Hostname()
	if err != nil {
		// use a random number
		clientName = fmt.Sprintf("race-cli-%d", rand.Intn(100))
	}
}

// supported commands
//		ping to make sure server is ok
//		quit/stop/exit the progeam
//		start event with or without seed time to start at
//		finish event with or without bib

func QuitCommand(args []string) (bool, error) {
	fmt.Println("quitting...")
	return true, nil
}

func NewPingCommand(rdb *redis.Client) CommandFunction {
	return func(args []string) (bool, error) {
		cmdResult := rdb.Ping()
		fmt.Println(cmdResult.String())
		return false, cmdResult.Err()
	}
}

func NewStartCommand(rdb *redis.Client, streamName string) CommandFunction {
	return func(args []string) (bool, error) {
		eventTarget := events.NewRedisStreamEventTarget(rdb, streamName)
		startTime := time.Now().UTC()
		seedTime := "0s"
		if len(args) > 0 {
			seedTime = args[0]
		}
		seedDuration, err := time.ParseDuration(seedTime)
		if err != nil {
			return false, err
		}

		return false, eventTarget.SendStart(events.StartEvent{
			Source:    clientName,
			StartTime: startTime.Add(-seedDuration),
		})
	}
}

func NewFinishCommand(rdb *redis.Client, streamName string) CommandFunction {
	return func(args []string) (bool, error) {
		eventTarget := events.NewRedisStreamEventTarget(rdb, streamName)

		bib := ""
		if len(args) > 0 {
			bib = args[0]
		}

		return false, eventTarget.SendFinish(events.FinishEvent{
			Source:     clientName,
			Bib:        bib,
			FinishTime: time.Now().UTC(),
		})
	}
}
