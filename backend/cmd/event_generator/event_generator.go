package main

// Program to read a race result file and generate events
// Phase 1:  Generate 1 event for each row, fully populated
// Phase 2:  Generate 2-3 events for each row within a few ms, some missing info, some with all
// Phase 3: emulate event volume by generating on schedule.  Send a start, then finishes, default placer should do place events
// then allow manual place fixes
// run score in real time or after all finshers are in
// verify scoring can be run multiple times to get results
// (ie the stream doesn't care it's been read already)

// Modify this to generate: start, finish (and place events?) required to use existing results to drive a test of
// whole system
// place events may be needed to distinguish the order of finish if times are the same

import (
	"blreynolds4/event-race-timer/internal/meets"
	"blreynolds4/event-race-timer/internal/raceevents"
	"blreynolds4/event-race-timer/internal/redis_stream"
	"bufio"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver

	redis "github.com/redis/go-redis/v9"
)

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func main() {
	// create a logger
	logger := newLogger()

	var claDbAddress string
	var claDbNumber int
	var claRacename string
	var claPostgresConnect string
	var claSourceFile string

	flag.StringVar(&claDbAddress, "redisAddress", "localhost:6379", "The host and port ie localhost:6379")
	flag.IntVar(&claDbNumber, "redisDbNumber", 0, "The database to use, defaults to 0")
	flag.StringVar(&claPostgresConnect, "pgConnect", "postgres://eventtimer:eventtimer@localhost:5432/eventtimer?sslmode=disable", "PostgreSQL connection string")
	flag.StringVar(&claRacename, "raceName", "", "The name of the race being timed")
	flag.StringVar(&claSourceFile, "sourceFile", "", "The name of the source file for the race events.")

	// parse command line
	flag.Parse()

	if strings.TrimSpace(claRacename) == "" {
		logger.Error("raceName is required")
		os.Exit(1)
	}

	// create a meet
	meetWriter, err := meets.NewMeetWriter(claPostgresConnect)
	if err != nil {
		logger.Error("ERROR creating meet writer", "error", err)
		os.Exit(1)
	}
	defer meetWriter.Close()

	meet, err := meetWriter.SaveMeet(&meets.Meet{
		Name: "Generated Meet " + time.Now().Format("2006-01-02-15-04-05"),
	})
	if err != nil {
		logger.Error("ERROR saving meet", "error", err)
		os.Exit(1)
	}

	raceWriter, err := meets.NewRaceWriter(claPostgresConnect)
	if err != nil {
		logger.Error("ERROR creating race writer", "error", err)
		os.Exit(1)
	}
	defer raceWriter.Close()

	race, err := raceWriter.SaveRace(&meets.Race{Name: claRacename}, meet)
	if err != nil {
		logger.Error("ERROR saving race", "error", err)
		os.Exit(1)
	}

	athleteWriter, err := meets.NewAthleteWriter(claPostgresConnect)
	if err != nil {
		logger.Error("ERROR creating athlete writer", "error", err)
		os.Exit(1)
	}
	defer athleteWriter.Close()

	// connect to redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     claDbAddress,
		Password: "",          // no password set
		DB:       claDbNumber, // use default DB
	})

	defer rdb.Close()

	// read the source file
	eventFile, err := os.Open(claSourceFile)
	if err != nil {
		fmt.Printf("error opening file %s: %s", claSourceFile, err.Error())
		os.Exit(-1)
	}

	// create event target
	rawStream := redis_stream.NewRedisStream(rdb, claRacename)
	eventStream := raceevents.NewEventStream(rawStream)

	// create and save competitor data
	athletes := make(meets.AthleteLookup)

	// send a start event
	startTime := time.Now().UTC()
	err = eventStream.SendStartEvent(context.TODO(), raceevents.StartEvent{
		Source:    "manual",
		StartTime: startTime,
	})
	if err != nil {
		fmt.Printf("error sending start event: %s", err.Error())
		os.Exit(-1)
	}

	// scan the lines of the race file and generate finish events to
	// match the race times based on starTime
	scanner := bufio.NewScanner(eventFile)
	done := false
	finishCount := 0
	for !done {
		// read a line of event input
		if scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "#") {
				// skip comment line
				continue
			}

			split := strings.Split(line, "|")
			if len(split) > 0 {
				// add competitor to lookup
				bib, err := strconv.Atoi(split[1])
				if err != nil {
					fmt.Printf("error getting bib from %s: %s", split[1], err.Error())
					os.Exit(-1)
				}
				grade, _ := strconv.Atoi(split[4])
				if err != nil {
					fmt.Printf("error getting grade from %s: %s", split[4], err.Error())
					os.Exit(-1)
				}

				c := new(meets.Athlete)
				c.DaID = fmt.Sprintf("gen-%d", bib)
				c.FirstName = split[3]
				c.LastName = split[2]
				c.Team = split[5]
				c.Grade = grade
				athletes[bib] = c

				athlete, err := athleteWriter.SaveAthlete(c)
				if err != nil {
					fmt.Printf("error saving athlete %s: %s", c.DaID, err.Error())
					os.Exit(-1)
				}

				err = raceWriter.AddAthlete(race, athlete, bib)
				if err != nil {
					fmt.Printf("error adding athlete %d to race %s: %s", bib, race.Name, err.Error())
					os.Exit(-1)
				}

				athletes[bib] = athlete

				finishCount++

				// col           0     1     2     3      4      5       6     7
				// splits are: place, bib, last, first, grade, school, time, score
				// use place for bib
				timeString := split[6]
				// durations are in minutes:seconds.tenths
				// convert to go duration format
				durationString := strings.Replace(timeString, ":", "m", 1) + "s"
				runDuration, err := time.ParseDuration(durationString)
				if err != nil {
					fmt.Printf("error getting duration from %s: %s -> %s %s", split[0], timeString, durationString, err.Error())
					os.Exit(-1)
				}

				// generate a finish event for each timestamp from the finish
				eventStream.SendFinishEvent(context.TODO(), raceevents.FinishEvent{
					Source:     "reader-1",
					FinishTime: startTime.Add(runDuration),
					Bib:        bib,
				})

				eventStream.SendPlaceEvent(context.TODO(), raceevents.PlaceEvent{
					Source: "manual",
					Bib:    bib,
					Place:  finishCount,
				})

				// get a random number 1 - 3 to decide on additional finish events for the athlete
				random := rand.Intn(3)
				if random >= 1 {
					// add a another event a little slower than first event with no bib
					// set a bib about half the time
					if rand.Intn(2) > 0 {
						bib = raceevents.NoBib
					}
					eventStream.SendFinishEvent(context.TODO(), raceevents.FinishEvent{
						Source:     "generator-slow",
						FinishTime: startTime.Add(runDuration).Add(time.Millisecond * 500),
						Bib:        bib,
					})
				}

				if random >= 2 {
					// add a third reader event a little faster
					eventStream.SendFinishEvent(context.TODO(), raceevents.FinishEvent{
						Source:     "generator-fast",
						FinishTime: startTime.Add(runDuration).Add(time.Millisecond * -500),
						Bib:        bib,
					})
				}
			}
		} else {
			done = true
		}
	}
}
