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
	"blreynolds4/event-race-timer/internal/competitors"
	"blreynolds4/event-race-timer/internal/raceevents"
	"blreynolds4/event-race-timer/internal/redis_stream"
	"bufio"
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	redis "github.com/redis/go-redis/v9"
)

func main() {
	// connect to redis
	// cli for db address, username, password, db, stream name?
	// stream is specific to a race
	var claDbAddress string
	var claDbNumber int
	var claRacename string
	var claSourceFile string

	flag.StringVar(&claDbAddress, "dbAddress", "localhost:6379", "The host and port ie localhost:6379")
	flag.IntVar(&claDbNumber, "dbNumber", 0, "The database to use, defaults to 0")
	flag.StringVar(&claRacename, "raceName", "race", "The name of the race being timed (no spaces)")
	flag.StringVar(&claSourceFile, "sourceFile", "", "The name of the source file for the race events.")

	// parse command line
	flag.Parse()

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
	rawStream := redis_stream.NewRedisEventStream(rdb, claRacename)
	eventStream := raceevents.NewEventStream(rawStream)

	// create and save competitor data
	athletes := make(competitors.CompetitorLookup)

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

				c := competitors.Competitor{
					Name:  split[3] + " " + split[2],
					Team:  split[5],
					Grade: int(grade),
				}

				athletes[bib] = &c

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

		// Save the competitor data for these events
		athletes.Store(claSourceFile + "_athletes.json")
	}
}
