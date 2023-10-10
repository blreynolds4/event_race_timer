package main

import (
	"blreynolds4/event-race-timer/overall"
	"blreynolds4/event-race-timer/redis_stream"
	"blreynolds4/event-race-timer/results"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	redis "github.com/redis/go-redis/v9"
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

	resultStreamName := claRacename + "_results"

	fmt.Println("building result reader reader for", resultStreamName)
	rawResultStream := redis_stream.NewRedisEventStream(rdb, claRacename+"_results")
	resultStream := results.NewResultStream(rawResultStream)

	resultScorer := overall.NewOverallResults()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	bgCtx := context.Background()
	cancelCtx, cancel := context.WithCancel(bgCtx)

	go func() {
		<-c
		// cancel the scorer
		cancel()
		fmt.Println("Cancelled Scorer")
	}()

	fmt.Println("building overall results")
	err := resultScorer.ScoreResults(cancelCtx, resultStream)
	if err != nil {
		fmt.Println("ERROR scoring results:", err)
	}

	fmt.Println("Scoring Exiting")
}
