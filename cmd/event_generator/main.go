package main

// Program to read a race result file and generate events
// Phase 1:  Generate 1 event for each row, fully populated
// Phase 2:  Generate 2-3 events for each row within a few ms, some missing info, some with all
import (
	"blreynolds4/event-race-timer/events"
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v7"
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
	eventTarget := events.NewRedisStreamEventTarget(rdb, claRacename)

	// send a start event
	startTime := time.Now().UTC()
	err = eventTarget.SendStart(events.StartEvent{
		Source:    "generator",
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

				// generate a finish event for each place with the correct timestamp from the finish
				eventTarget.SendFinish(events.FinishEvent{
					Source:     "generator-1",
					Bib:        split[0],
					FinishTime: startTime.Add(runDuration),
				})

				// get a random number 1 - 3 to decide on additional finish events for the athlete
				random := rand.Intn(3)
				if random >= 1 {
					// add a manual event a little slower than first event with no bib
					// set a bib about half the time
					bib := split[0]
					if rand.Intn(2) > 0 {
						bib = ""
					}
					eventTarget.SendFinish(events.FinishEvent{
						Source:     "generator-manual",
						Bib:        bib,
						FinishTime: startTime.Add(runDuration).Add(time.Millisecond * 500),
					})
				}

				if random >= 2 {
					// add a third reader event a little faster
					eventTarget.SendFinish(events.FinishEvent{
						Source:     "generator-2",
						Bib:        split[0],
						FinishTime: startTime.Add(runDuration).Add(time.Millisecond * -500),
					})
				}
			}
		} else {
			done = true
		}
	}
}
