package main

import (
	"blreynolds4/event-race-timer/cmd/cli/internal/cli"
	"flag"
	"log/slog"
	"os"
	"strings"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func main() {
	// connect to redis
	// cli for db address, username, password, db, stream name?
	// stream is specific to a race
	var claDbAddress string
	var claDbNumber int
	var claRacename string
	var claPostgresConnect string

	flag.StringVar(&claDbAddress, "dbAddress", "localhost:6379", "The host and port ie localhost:6379")
	flag.IntVar(&claDbNumber, "dbNumber", 0, "The database to use, defaults to 0")
	flag.StringVar(&claRacename, "raceName", "", "The name of the race being timed (no spaces)")
	flag.StringVar(&claPostgresConnect, "pgConnect", "postgres://eventtimer:eventtimer@localhost:5432/eventtimer?sslmode=disable", "PostgreSQL connection string")

	// parse command line
	flag.Parse()

	logger := newLogger()

	if strings.TrimSpace(claRacename) == "" {
		logger.Error("raceName is required")
		os.Exit(1)
	}

	app := cli.NewCliApp()

	app.Run(claDbAddress, claDbNumber, claRacename, claPostgresConnect)
}
