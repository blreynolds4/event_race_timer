package handler

import (
	"blreynolds4/event-race-timer/internal/config"
	"blreynolds4/event-race-timer/internal/meets"
	"blreynolds4/event-race-timer/internal/raceevents"
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func NewTimingHandler(sourceLookup config.SourceConfig, athletes meets.AthleteLookup, eventStream raceevents.EventStream, logger *slog.Logger) gin.HandlerFunc {
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
			// send the finish event
			// using the event time from the JSON payload as the finish time
			finishTime := time.UnixMilli(int64(data.EventTime)).UTC()
			eventStream.SendFinishEvent(context.TODO(), raceevents.FinishEvent{
				Source:     sourceLookup.SourceMap[data.Host],
				Bib:        bib,
				FinishTime: finishTime,
			})

			// respond with the data and a 201 created
			c.IndentedJSON(http.StatusCreated, data)
			logger.Info("sent race finish event",
				"bib", bib,
				"finish time", finishTime,
				"antenna", data.Antenna,
				"source", sourceLookup.SourceMap[data.Host],
				"host", data.Host)
		} else {
			logger.Warn("skipping unknown bib", "bib", bib, "antenna", data.Antenna, "host", data.Host)
		}
	}

	return gin.HandlerFunc(fn)
}
