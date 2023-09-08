package main

import (
	"blreynolds4/event-race-timer/eventstream"
	"blreynolds4/event-race-timer/places"
	"blreynolds4/event-race-timer/redis_stream"
	"flag"
	"fmt"

	redis "github.com/redis/go-redis/v9"
)

func main() {
	// connect to redis
	// cli for db address, username, password, db, stream name?
	// stream is specific to a race
	var claDbAddress string
	var claDbNumber int
	var claRacename string

	flag.StringVar(&claDbAddress, "dbAddress", "localhost:6379", "The host and port ie localhost:6379")
	flag.IntVar(&claDbNumber, "dbNumber", 0, "The database to use, defaults to 0")
	flag.StringVar(&claRacename, "raceName", "race", "The name of the race being timed (no spaces)")

	// parse command line
	flag.Parse()

	// connect to redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     claDbAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	defer rdb.Close()

	rawRead := redis_stream.NewRedisStreamReader(rdb, claRacename)
	src := eventstream.NewRaceEventSource(rawRead, eventstream.StreamMessageToRaceEvent)

	rawWrite := redis_stream.NewRedisStreamWriter(rdb, claRacename)
	target := eventstream.NewRaceEventTarget(rawWrite, eventstream.RaceEventToStreamMessage)

	placer := places.NewPlaceGenerator(src, target)

	err := placer.GeneratePlaces()
	if err != nil {
		fmt.Println("ERROR generating places:", err)
	}
}
