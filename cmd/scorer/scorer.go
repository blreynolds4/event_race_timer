package main

import (
	"blreynolds4/event-race-timer/cmd/scorer/internal/overall"
	"blreynolds4/event-race-timer/cmd/scorer/internal/xc"
	"blreynolds4/event-race-timer/internal/redis_stream"
	"blreynolds4/event-race-timer/internal/results"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	redis "github.com/redis/go-redis/v9"
)

func main() {
	// connect to redis
	// cli for db address, username, password, db, stream name?
	// stream is specific to a race
	var claDbAddress string
	var claDbNumber int
	var claRacename string
	var claOverall bool
	var claXCTeam bool
	var claDebug bool

	flag.StringVar(&claDbAddress, "dbAddress", "localhost:6379", "The host and port ie localhost:6379")
	flag.IntVar(&claDbNumber, "dbNumber", 0, "The database to use, defaults to 0")
	flag.StringVar(&claRacename, "raceName", "race", "The name of the race being timed (no spaces)")
	flag.BoolVar(&claOverall, "overall", false, "Use this flag to turn on overall scoring")
	flag.BoolVar(&claXCTeam, "xc", false, "Use this flag to turn on XC team scoring")
	flag.BoolVar(&claDebug, "debug", false, "Use this flag to debug")

	// parse command line
	flag.Parse()

	// connect to redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     claDbAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	defer rdb.Close()

	resultStreamName := claRacename + "_results"
	if claDebug {
		resultStreamName = resultStreamName + "_debug"
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	for {
		fmt.Println("building result reader reader for", resultStreamName)
		rawResultStream := redis_stream.NewRedisStream(rdb, resultStreamName)
		resultStream := results.NewResultStream(rawResultStream)

		if claXCTeam {
			xcScorer := xc.NewXCScorer()
			err := xcScorer.ScoreResults(context.TODO(), resultStream)
			if err != nil {
				fmt.Println("ERROR scoring xc results:", err)
			}
		}

		if claOverall {
			resultScorer := overall.NewOverallResults()
			err := resultScorer.ScoreResults(context.TODO(), resultStream)
			if err != nil {
				fmt.Println("ERROR scoring overall results:", err)
			}
		}

		t := time.NewTicker(time.Second * 2)
		select {
		case <-c:
			fmt.Println("Scorer Exiting")
			os.Exit(0)
		case <-t.C:
		}
	}
}
