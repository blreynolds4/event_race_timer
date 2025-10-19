package main

import (
	"blreynolds4/event-race-timer/cmd/race_archiver/internal/racearchive/archiver"
	"blreynolds4/event-race-timer/cmd/race_archiver/internal/racearchive/restorer"
	"blreynolds4/event-race-timer/internal/raceevents"
	"blreynolds4/event-race-timer/internal/redis_stream"
	"flag"
	"log"
	"os"

	redis "github.com/redis/go-redis/v9"
)

const (
	saveAction    = "save"
	restoreAction = "restore"

	archiveExtension  = ".archive.json"
	restoredExtension = "_restored"
)

func main() {
	// connect to redis
	// cli for db address, username, password, db, stream name?
	// stream is specific to a race
	var claDbAddress string
	var claDbNumber int
	var claRacename string
	var claAction string

	flag.StringVar(&claDbAddress, "dbAddress", "localhost:6379", "The host and port ie localhost:6379")
	flag.IntVar(&claDbNumber, "dbNumber", 0, "The database to use, defaults to 0")
	flag.StringVar(&claRacename, "raceName", "race", "The name of the race being timed (no spaces)")
	flag.StringVar(&claAction, "action", saveAction, "The action to take:  save or restore")

	// parse command line
	flag.Parse()

	// connect to redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     claDbAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()

	switch {

	case claAction == saveAction:
		// create the race event stream to archive
		rawStream := redis_stream.NewRedisStream(rdb, claRacename)
		eventStream := raceevents.NewEventStream(rawStream)

		// create the file for the events
		f, err := os.Create(claRacename + archiveExtension)
		if err != nil {
			log.Fatal("failed to open output file")
		}
		defer f.Close()

		// archive to the file
		archiver := archiver.NewJsonFileArchiver(f)
		err = archiver.Archive(eventStream)
		if err != nil {
			log.Fatalf("error archiving %s: %s\n", claRacename, err)
		}
	case claAction == restoreAction:
		rawStream := redis_stream.NewRedisStream(rdb, claRacename+restoredExtension)
		eventStream := raceevents.NewEventStream(rawStream)

		f, err := os.Open(claRacename + archiveExtension)
		if err != nil {
			log.Fatal("failed to open archive file")
		}
		defer f.Close()

		restorer := restorer.NewRestorer()
		err = restorer.Restore(f, eventStream)
		if err != nil {
			log.Fatalf("error restoring %s: %s\n", claRacename, err)
		}
	default:
		log.Fatal("Unknown action", claAction)
	}
}
