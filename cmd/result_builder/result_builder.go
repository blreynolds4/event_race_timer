package main

import (
	"blreynolds4/event-race-timer/cmd/result_builder/internal/resultbuilder"
	"blreynolds4/event-race-timer/internal/config"
	"blreynolds4/event-race-timer/internal/meets"
	"blreynolds4/event-race-timer/internal/raceevents"
	"blreynolds4/event-race-timer/internal/redis_stream"
	"flag"
	"log/slog"
	"os"
	"strings"

	_ "github.com/lib/pq" // PostgreSQL driver

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
	var claPostgresConnect string
	var claRacename string
	var claConfigPath string
	var claDebug bool
	var claPlaceFile string

	flag.StringVar(&claDbAddress, "redisAddress", "localhost:6379", "The host and port ie localhost:6379")
	flag.IntVar(&claDbNumber, "redisDbNumber", 0, "The database to use, defaults to 0")
	flag.StringVar(&claPostgresConnect, "pgConnect", "postgres://eventtimer:eventtimer@localhost:5432/eventtimer?sslmode=disable", "PostgreSQL connection string")
	flag.StringVar(&claRacename, "raceName", "", "The name of the race being timed")
	flag.StringVar(&claConfigPath, "config", "", "The path to the config file (json)")
	flag.BoolVar(&claDebug, "debug", false, "Flag to debug")
	flag.StringVar(&claPlaceFile, "places", "", "The path to the place file")

	// parse command line
	flag.Parse()

	if strings.TrimSpace(claRacename) == "" {
		logger.Error("raceName is required")
		os.Exit(1)
	}

	// connect to redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     claDbAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	defer rdb.Close()

	raceReader, err := meets.NewRaceReader(claPostgresConnect)
	if err != nil {
		logger.Error("ERROR creating race reader", "error", err)
		os.Exit(1)
	}
	defer raceReader.Close()

	race, err := raceReader.GetRaceByName(claRacename)
	if err != nil {
		logger.Error("ERROR loading race", "race", claRacename, "error", err)
		os.Exit(1)
	}

	resultsWriter, err := meets.NewRaceResultWriter(race, claPostgresConnect)
	if err != nil {
		logger.Error("ERROR creating results writer", "error", err)
		os.Exit(1)
	}

	athletes := make(meets.AthleteLookup)
	err = meets.LoadAthleteLookup(claPostgresConnect, claRacename, athletes)
	if err != nil {
		logger.Error("error loading athletes", "error", err)
		os.Exit(1)
	}

	logger.Info("Loaded athletes", "count", len(athletes), "athletes", athletes)

	var raceConfig config.RaceConfig
	err = config.LoadConfigData(claConfigPath, &raceConfig)
	if err != nil {
		logger.Error("ERROR loading config", "filename", claConfigPath, "error", err)
		os.Exit(1)
	}

	rawStream := redis_stream.NewRedisStream(rdb, claRacename)
	eventStream := raceevents.NewEventStream(rawStream)

	resultBuilder := resultbuilder.NewRaceResultBuilder(logger)

	err = resultBuilder.BuildRaceResults(eventStream, athletes, raceConfig.SourceRanks, resultsWriter)
	if err != nil {
		logger.Error("ERROR generating results", "error", err)
	}
}
