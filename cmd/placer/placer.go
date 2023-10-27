package main

import (
	"blreynolds4/event-race-timer/cmd/placer/internal/places"
	"blreynolds4/event-race-timer/internal/competitors"
	"blreynolds4/event-race-timer/internal/config"
	"blreynolds4/event-race-timer/internal/raceevents"
	"blreynolds4/event-race-timer/internal/redis_stream"
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
	var claConfigPath string
	var claCompetitorsPath string

	flag.StringVar(&claDbAddress, "dbAddress", "localhost:6379", "The host and port ie localhost:6379")
	flag.IntVar(&claDbNumber, "dbNumber", 0, "The database to use, defaults to 0")
	flag.StringVar(&claRacename, "raceName", "race", "The name of the race being timed (no spaces)")
	flag.StringVar(&claConfigPath, "config", "", "The path to the config file (json)")
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

	rawStream := redis_stream.NewRedisEventStream(rdb, claRacename)
	eventStream := raceevents.NewEventStream(rawStream)

	var raceConfig config.RaceConfig
	err := config.LoadConfigData(claConfigPath, &raceConfig)
	if err != nil {
		fmt.Printf("ERROR loading config from '%s': %v\n", claConfigPath, err)
		os.Exit(1)
	}

	athletes := make(competitors.CompetitorLookup)
	err = competitors.LoadCompetitorLookup(claCompetitorsPath, athletes)
	if err != nil {
		fmt.Printf("ERROR loading competitors from '%s': %v\n", claCompetitorsPath, err)
		os.Exit(1)
	}

	placer := places.NewPlaceGenerator(eventStream)

	err = placer.GeneratePlaces(athletes, raceConfig.SourceRanks)
	if err != nil {
		fmt.Println("ERROR generating places:", err)
	}
}
