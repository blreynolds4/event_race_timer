package main

import (
	"blreynolds4/event-race-timer/cmd/scorer/internal/overall"
	"blreynolds4/event-race-timer/cmd/scorer/internal/xc"
	"blreynolds4/event-race-timer/internal/meets"
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func main() {
	// create a logger
	logger := newLogger()

	var claRacename string
	var claPostgresConnect string
	var claOverall bool
	var claXCTeam bool
	var claDebug bool

	flag.StringVar(&claRacename, "raceName", "race", "The name of the race being timed (no spaces)")
	flag.StringVar(&claPostgresConnect, "pgConnect", "postgres://eventtimer:eventtimer@localhost:5432/eventtimer?sslmode=disable", "PostgreSQL connection string")
	flag.BoolVar(&claOverall, "overall", false, "Use this flag to turn on overall scoring")
	flag.BoolVar(&claXCTeam, "xc", false, "Use this flag to turn on XC team scoring")
	flag.BoolVar(&claDebug, "debug", false, "Use this flag to debug")

	// parse command line
	flag.Parse()

	if strings.TrimSpace(claRacename) == "" {
		logger.Error("raceName is required")
		os.Exit(1)
	}

	raceReader, err := meets.NewRaceReader(claPostgresConnect)
	if err != nil {
		logger.Error("ERROR creating race reader", "error", err)
		os.Exit(1)
	}

	race, err := raceReader.GetRaceByName(claRacename)
	if err != nil {
		logger.Error("ERROR getting race by name", "error", err)
		os.Exit(1)
	}

	raceResultsReader, err := meets.NewRaceResultReader(race, claPostgresConnect)
	if err != nil {
		logger.Error("ERROR creating race results reader", "error", err)
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	for {
		if claXCTeam {
			xcScorer := xc.NewXCTeamScorer(race, logger)
			err := xcScorer.ScoreResults(raceResultsReader)
			if err != nil {
				logger.Error("ERROR scoring xc results", "error", err)
			}
		}

		if claOverall {
			resultScorer := overall.NewOverallRaceResults(race, logger)
			err := resultScorer.ScoreResults(context.TODO(), raceResultsReader)
			if err != nil {
				logger.Error("ERROR scoring overall race results", "error", err)
			}
		}

		t := time.NewTicker(time.Second * 2)
		select {
		case <-c:
			logger.Info("Scorer Exiting")
			os.Exit(0)
		case <-t.C:
		}
	}
}
