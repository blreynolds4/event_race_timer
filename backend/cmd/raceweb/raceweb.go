package main

import (
	"flag"
	"log/slog"
	"os"
	"strings"

	"blreynolds4/event-race-timer/cmd/raceweb/internal/raceweb"
	"blreynolds4/event-race-timer/internal/config"
	"blreynolds4/event-race-timer/internal/meets"
	"blreynolds4/event-race-timer/internal/raceevents"
	"blreynolds4/event-race-timer/internal/redis_stream"

	_ "github.com/lib/pq" // PostgreSQL driver
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
	var claPostgresConnect string

	flag.StringVar(&claSourceConfig, "config", "", "The config file for sources")
	flag.StringVar(&claDbAddress, "redisAddress", "localhost:6379", "The host and port ie localhost:6379")
	flag.IntVar(&claDbNumber, "redisDbNumber", 0, "The database to use, defaults to 0")
	flag.StringVar(&claPostgresConnect, "pgConnect", "postgres://eventtimer:eventtimer@localhost:5432/eventtimer?sslmode=disable", "PostgreSQL connection string")
	flag.StringVar(&claRacename, "raceName", "", "The name of the race being timed")

	flag.Parse()

	// load config data for host/ip to source name like mat and chute reader
	var sources config.SourceConfig
	err := config.LoadAnyConfigData[config.SourceConfig](claSourceConfig, &sources)
	if err != nil {
		logger.Error("error loading config", "fileName", claSourceConfig, "error", err)
		os.Exit(1)
	}

	if strings.TrimSpace(claRacename) == "" {
		logger.Error("raceName is required")
		os.Exit(1)
	}

	athletes := make(meets.AthleteLookup)
	err = meets.LoadAthleteLookup(claPostgresConnect, claRacename, athletes)
	if err != nil {
		logger.Error("error loading athletes", "error", err)
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

	meetReader, err := meets.NewMeetReader(claPostgresConnect)
	if err != nil {
		logger.Error("error creating meet reader", "error", err)
		panic(err)
	}

	app := raceweb.NewApplication(sources, athletes, meetReader, eventStream, logger)

	app.Run(":8080")
}
