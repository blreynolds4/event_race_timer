package main

import (
	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/events"
	"blreynolds4/event-race-timer/redis_stream"
	"blreynolds4/event-race-timer/results"
	"flag"
	"fmt"
	"os"

	redis "github.com/redis/go-redis/v9"
)

func main() {
	// connect to redis
	// cli for db address, username, password, db, stream name?
	// stream is specific to a race
	var claDbAddress string
	var claDbNumber int
	var claRacename string
	var claCompetitorsPath string

	flag.StringVar(&claDbAddress, "dbAddress", "localhost:6379", "The host and port ie localhost:6379")
	flag.IntVar(&claDbNumber, "dbNumber", 0, "The database to use, defaults to 0")
	flag.StringVar(&claRacename, "raceName", "race", "The name of the race being timed (no spaces)")
	flag.StringVar(&claCompetitorsPath, "competitors", "", "The path to the competitor lookup file (json)")

	// parse command line
	flag.Parse()

	// connect to redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     claDbAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	defer rdb.Close()

	athletes, err := competitors.LoadCompetitorLookup(claCompetitorsPath)
	if err != nil {
		fmt.Printf("ERROR loading competitors from '%s': %v\n", claCompetitorsPath, err)
		os.Exit(1)
	}

	rawRead := redis_stream.NewRedisStreamReader(rdb, claRacename)
	raceEventSrc := events.NewRaceEventSource(rawRead)

	rawWrite := redis_stream.NewRedisStreamWriter(rdb, claRacename+"_results")
	resultEventTarget := results.NewResultTarget(rawWrite)

	resultBuilder := results.NewResultBuilder()
	err = resultBuilder.BuildResults(raceEventSrc, athletes, resultEventTarget)
	if err != nil {
		fmt.Println("ERROR generating results:", err)
	}
}
