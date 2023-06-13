package command

import (
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
		startTime := time.Now().UTC()
		seedTime := "0s"
		if len(args) > 0 {
			seedTime = args[0]
		}
		seedDuration, err := time.ParseDuration(seedTime)
		if err != nil {
			return false, err
		}

		addArgs := redis.XAddArgs{
			Stream: streamName,
			Values: map[string]interface{}{
				"event_type": "start",
				"start_time": startTime.Add(-seedDuration).UnixMilli(),
				"source":     clientName,
			},
		}
		result := rdb.XAdd(&addArgs)
		if result.Err() != nil {
			return false, result.Err()
		}

		fmt.Println("ok -", result.Val())
		return false, nil
	}
}

func NewFinishCommand(rdb *redis.Client, streamName string) CommandFunction {
	return func(args []string) (bool, error) {
		bib := ""
		if len(args) > 0 {
			bib = args[0]
		}
		addArgs := redis.XAddArgs{
			Stream: streamName,
			Values: map[string]interface{}{
				"event_type":  "finish",
				"bib":         bib,
				"finish_time": time.Now().UTC().UnixMilli(),
				"source":      clientName,
			},
		}
		result := rdb.XAdd(&addArgs)
		if result.Err() != nil {
			return false, result.Err()
		}

		fmt.Println("ok -", result.Val())
		return false, nil
	}
}
