package main

import (
	"blreynolds4/event-race-timer/cmd/placer/internal/places"
	"blreynolds4/event-race-timer/internal/competitors"
	"blreynolds4/event-race-timer/internal/config"
	"blreynolds4/event-race-timer/internal/raceevents"
	"blreynolds4/event-race-timer/internal/redis_stream"
	"flag"
	"log/slog"
	"os"

	redis "github.com/redis/go-redis/v9"
)

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func main() {
	// create a logger
	logger := newLogger()

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

	rawStream := redis_stream.NewRedisStream(rdb, claRacename)
	eventStream := raceevents.NewEventStream(rawStream)

	var raceConfig config.RaceConfig
	err := config.LoadConfigData(claConfigPath, &raceConfig)
	if err != nil {
		logger.Error("ERROR loading config", "fileName", claConfigPath, "error", err)
		os.Exit(1)
	}

	athletes := make(competitors.CompetitorLookup)
	err = competitors.LoadCompetitorLookup(claCompetitorsPath, athletes)
	if err != nil {
		logger.Error("ERROR loading competitors from", "fileName", claCompetitorsPath, "error", err)
		os.Exit(1)
	}

	placer := places.NewPlaceGenerator(eventStream, logger)

	err = placer.GeneratePlaces(athletes, raceConfig.SourceRanks)
	if err != nil {
		logger.Error("ERROR generating places", "error", err)
	}
}
