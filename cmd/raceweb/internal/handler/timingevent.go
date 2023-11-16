package handler

import (
	"blreynolds4/event-race-timer/internal/competitors"
	"blreynolds4/event-race-timer/internal/config"
	"blreynolds4/event-race-timer/internal/raceevents"
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type OpenSignupsTimingEvent struct {
	ElapsedTime int    `json:"elapsedTime"`
	CaptureMode string `json:"captureMode"`
	StartTime   int    `json:"startTime"`
	Antenna     int    `json:"antenna"`
	Bib         string `json:"bib"`
	Host        string `json:"host"`
}

func NewTimingHandler(sourceLookup config.SourceConfig, athletes competitors.CompetitorLookup, eventStream raceevents.EventStream, logger *slog.Logger) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var data OpenSignupsTimingEvent
		if err := c.BindJSON(&data); err != nil {
			logger.Error("bind json error", "error", err.Error())
			return
		}
		bib, err := strconv.Atoi(data.Bib)
		if err != nil {
			// just drop the bad bib
			logger.Error("dropping bad bib", "bib", data.Bib, "antenna", data.Antenna, "host", data.Host)
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
			logger.Info("sent finish event", "bib", bib, "elapsedTime", time.Duration(data.ElapsedTime*int(time.Millisecond)).String(), "antenna", data.Antenna, "source", sourceLookup.SourceMap[data.Host], "host", data.Host)
		} else {
			logger.Info("skipping unknown bib", "bib", bib, "antenna", data.Antenna, "host", data.Host)
		}
	}

	return gin.HandlerFunc(fn)
}
