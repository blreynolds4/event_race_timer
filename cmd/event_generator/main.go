package main

// Program to read a race result file and generate events
// Phase 1:  Generate 1 event for each row, fully populated
// Phase 2:  Generate 2-3 events for each row within a few ms, some missing info, some with all
import (
	"blreynolds4/event-race-timer/events"
	"blreynolds4/event-race-timer/eventstream"
	"blreynolds4/event-race-timer/redis_stream"
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
	rawWrite := redis_stream.NewRedisStreamWriter(rdb, claRacename)
	eventTarget := eventstream.NewRaceEventTarget(rawWrite, eventstream.RaceEventToStreamMessage)

	// send a start event
	startTime := time.Now().UTC()
	err = eventTarget.SendRaceEvent(context.TODO(), eventstream.NewStartEvent("generator", startTime))
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

			split := strings.Split(line, " ")
			if len(split) > 0 {
				finishCount++

				// use place for bib
				timeString := split[len(split)-2]
				// durations are in minutes:seconds.tenths
				// convert to go duration format
				durationString := strings.Replace(timeString, ":", "m", 1) + "s"
				runDuration, err := time.ParseDuration(durationString)
				if err != nil {
					fmt.Printf("error getting duration from %s: %s -> %s %s", split[0], timeString, durationString, err.Error())
					os.Exit(-1)
				}

				// generate a finish event for each timestamp from the finish
				bib, err := strconv.Atoi(split[0])
				if err != nil {
					fmt.Printf("error getting bib from %s: %s", split[0], err.Error())
					os.Exit(-1)
				}
				eventTarget.SendRaceEvent(context.TODO(), eventstream.NewFinishEvent("reader-1", startTime.Add(runDuration), bib))

				// get a random number 1 - 3 to decide on additional finish events for the athlete
				random := rand.Intn(3)
				if random >= 1 {
					// add a manual event a little slower than first event with no bib
					// set a bib about half the time
					if rand.Intn(2) > 0 {
						bib = events.NoBib
					}
					eventTarget.SendRaceEvent(context.TODO(), eventstream.NewFinishEvent("generator-manual", startTime.Add(runDuration).Add(time.Millisecond*500), bib))
				}

				if random >= 2 {
					// add a third reader event a little faster
					eventTarget.SendRaceEvent(context.TODO(), eventstream.NewFinishEvent("generator-manual", startTime.Add(runDuration).Add(time.Millisecond*-500), bib))
				}
			}
		} else {
			done = true
		}
	}

	// send a place event for each finish using the place as the bib
	for i := 1; i <= finishCount; i++ {
		eventTarget.SendRaceEvent(context.TODO(), eventstream.NewPlaceEvent("chute-manual", i, i))
	}
}
