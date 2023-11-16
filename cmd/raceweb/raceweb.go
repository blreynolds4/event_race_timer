package main

import (
	"flag"
	"log/slog"
	"os"

	"blreynolds4/event-race-timer/cmd/raceweb/internal/raceweb"
	"blreynolds4/event-race-timer/internal/competitors"
	"blreynolds4/event-race-timer/internal/config"
	"blreynolds4/event-race-timer/internal/raceevents"
	"blreynolds4/event-race-timer/internal/redis_stream"

	redis "github.com/redis/go-redis/v9"
)

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func main() {
	// create a logger
	logger := newLogger()

	var claSourceConfig string
	var claDbAddress string
	var claDbNumber int
	var claRacename string
	var claCompetitorsPath string

	flag.StringVar(&claSourceConfig, "config", "", "The config file for sources")
	flag.StringVar(&claDbAddress, "dbAddress", "localhost:6379", "The host and port ie localhost:6379")
	flag.IntVar(&claDbNumber, "dbNumber", 0, "The database to use, defaults to 0")
	flag.StringVar(&claRacename, "raceName", "race", "The name of the race being timed (no spaces)")
	flag.StringVar(&claCompetitorsPath, "competitors", "", "The path to the competitor lookup file (json)")

	flag.Parse()

	// load config data
	var sources config.SourceConfig
	err := config.LoadAnyConfigData[config.SourceConfig](claSourceConfig, &sources)
	if err != nil {
		logger.Error("error loading config", "fileName", claSourceConfig, "error", err)
		os.Exit(1)
	}

	athletes := make(competitors.CompetitorLookup)
	err = competitors.LoadCompetitorLookup(claCompetitorsPath, athletes)
	if err != nil {
		logger.Error("ERROR loading competitors", "fileName", claCompetitorsPath, "error", err)
		os.Exit(1)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     claDbAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	defer rdb.Close()

	rawStream := redis_stream.NewRedisStream(rdb, claRacename)
	eventStream := raceevents.NewEventStream(rawStream)

	app := raceweb.NewApplication(sources, athletes, eventStream, logger)

	app.Run(":8080")
}
