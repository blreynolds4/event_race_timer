package main

import (
	"blreynolds4/event-race-timer/events"
	"blreynolds4/event-race-timer/places"
	"flag"
	"fmt"

	redis "github.com/go-redis/redis/v7"
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

	src := events.NewRedisStreamEventSource(rdb, claRacename)
	target := events.NewRedisStreamEventTarget(rdb, claRacename)

	placer := places.NewPlaceGenerator(src, target)

	err := placer.GeneratePlaces()
	if err != nil {
		fmt.Println("ERROR generating places:", err)
	}
}
