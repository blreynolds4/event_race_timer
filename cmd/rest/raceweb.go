package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/config"
	"blreynolds4/event-race-timer/raceevents"
	"blreynolds4/event-race-timer/redis_stream"

	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
)

type emptyResponse struct {
}

type OpenSignupsTimingEvent struct {
	ElapsedTime int    `json:"elapsedTime"`
	CaptureMode string `json:"captureMode"`
	StartTime   int    `json:"startTime"`
	Antenna     int    `json:"antenna"`
	Bib         string `json:"bib"`
	Host        string `json:"host"`
}

// methods to handle rest requests
func verifyTimingEvent(c *gin.Context) {
	fmt.Println("Handling get timing events")
	c.IndentedJSON(http.StatusOK, emptyResponse{})
}

func NewTimingHandler(sourceLookup config.SourceConfig, athletes competitors.CompetitorLookup, eventStream *raceevents.EventStream) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var data OpenSignupsTimingEvent
		if err := c.BindJSON(&data); err != nil {
			fmt.Println("bind json error:", err.Error())
			return
		}
		bib, err := strconv.Atoi(data.Bib)
		if err != nil {
			// just drop the bad bib
			fmt.Println("dropping bad bib", data.Bib)
			return
		}

		if _, bibFound := athletes[bib]; bibFound {
			startTime := time.UnixMilli(int64(data.StartTime))
			finishTime := startTime.Add(time.Duration(data.ElapsedTime * int(time.Millisecond)))

			// send the finish event
			eventStream.SendFinishEvent(context.TODO(), raceevents.FinishEvent{
				Source:     sourceLookup.SourceMap[data.Host],
				Bib:        bib,
				FinishTime: finishTime,
			})
			c.IndentedJSON(http.StatusCreated, data)
			fmt.Printf("Finish Sent bib %d %s\n", bib, time.Duration(data.ElapsedTime*int(time.Millisecond)).String())
		} else {
			fmt.Println("skipping unknown bib", bib)
		}
	}

	return gin.HandlerFunc(fn)
}

// event timer needs to use configration to pick a source
// Config needs to use host to map to source
// the source goes into the finish event

// long term source could also be used to identify event type like start or place
// but not right now

func main() {
	var claSourceConfig string
	var claDbAddress string
	var claDbNumber int
	var claRacename string
	var claCompetitorsPath string

	flag.StringVar(&claSourceConfig, "config", "", "The config file for sources")
	flag.StringVar(&claDbAddress, "dbAddress", "localhost:6379", "The host and port ie localhost:6379")
	flag.IntVar(&claDbNumber, "dbNumber", 0, "The database to use, defaults to 0")
	flag.StringVar(&claRacename, "raceName", "race", "The name of the race being timed (no spaces)")
	flag.StringVar(&claCompetitorsPath, "competitors", "", "The path to the competitor lookup file (json)")

	flag.Parse()

	// load config data
	var sources config.SourceConfig
	err := config.LoadAnyConfigData[config.SourceConfig](claSourceConfig, &sources)
	if err != nil {
		log.Fatalf("error loading %s config %v", claSourceConfig, &sources)
	}

	athletes := make(competitors.CompetitorLookup)
	err = competitors.LoadCompetitorLookup(claCompetitorsPath, athletes)
	if err != nil {
		fmt.Printf("ERROR loading competitors from '%s': %v\n", claCompetitorsPath, err)
		os.Exit(1)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     claDbAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	defer rdb.Close()

	rawStream := redis_stream.NewRedisEventStream(rdb, claRacename)
	eventStream := raceevents.NewEventStream(rawStream)

	router := gin.Default()

	// Setup route group for the API
	api := router.Group("/api")
	api.GET("/timingEvents", verifyTimingEvent)
	api.POST("/timingEvents/finishes", NewTimingHandler(sources, athletes, eventStream))

	// results paths
	router.StaticFile("/overall", "overall_results.html")

	router.Run(":8080")
}
