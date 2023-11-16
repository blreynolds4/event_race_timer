package main

import (
	"blreynolds4/event-race-timer/cmd/result_builder/internal/resultbuilder"
	"blreynolds4/event-race-timer/internal/competitors"
	"blreynolds4/event-race-timer/internal/config"
	"blreynolds4/event-race-timer/internal/raceevents"
	"blreynolds4/event-race-timer/internal/redis_stream"
	"blreynolds4/event-race-timer/internal/results"
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
	var claCompetitorsPath string
	var claConfigPath string
	var claDebug bool
	var claPlaceFile string

	flag.StringVar(&claDbAddress, "dbAddress", "localhost:6379", "The host and port ie localhost:6379")
	flag.IntVar(&claDbNumber, "dbNumber", 0, "The database to use, defaults to 0")
	flag.StringVar(&claRacename, "raceName", "race", "The name of the race being timed (no spaces)")
	flag.StringVar(&claCompetitorsPath, "competitors", "", "The path to the competitor lookup file (json)")
	flag.StringVar(&claConfigPath, "config", "", "The path to the config file (json)")
	flag.BoolVar(&claDebug, "debug", false, "Flag to debug")
	flag.StringVar(&claPlaceFile, "places", "", "The path to the place file")

	// parse command line
	flag.Parse()

	// connect to redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     claDbAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	defer rdb.Close()

	athletes := make(competitors.CompetitorLookup)
	err := competitors.LoadCompetitorLookup(claCompetitorsPath, athletes)
	if err != nil {
		logger.Error("ERROR loading competitors", "filename", claCompetitorsPath, "error", err)
		os.Exit(1)
	}

	var raceConfig config.RaceConfig
	err = config.LoadConfigData(claConfigPath, &raceConfig)
	if err != nil {
		logger.Error("ERROR loading config", "filename", claConfigPath, "error", err)
		os.Exit(1)
	}

	rawStream := redis_stream.NewRedisStream(rdb, claRacename)
	eventStream := raceevents.NewEventStream(rawStream)

	resultStreamName := claRacename + "_results"
	if claDebug {
		resultStreamName = resultStreamName + "_debug"
	}
	rawResultStream := redis_stream.NewRedisStream(rdb, resultStreamName)
	resultStream := results.NewResultStream(rawResultStream)

	resultBuilder := resultbuilder.NewResultBuilder(logger)
	if claDebug {
		logger.Info("DEBUGGING RESULTS")
		resultBuilder = resultbuilder.NewStartFinishResultBuilder(claPlaceFile, logger)
	}

	err = resultBuilder.BuildResults(eventStream, athletes, resultStream, raceConfig.SourceRanks)
	if err != nil {
		logger.Error("ERROR generating results", "error", err)
	}
}
