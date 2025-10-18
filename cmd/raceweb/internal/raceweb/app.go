package raceweb

import (
	"blreynolds4/event-race-timer/cmd/raceweb/internal/handler"
	"blreynolds4/event-race-timer/internal/config"
	"blreynolds4/event-race-timer/internal/meets"
	"blreynolds4/event-race-timer/internal/raceevents"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type Application interface {
	Run(address string)
}

type application struct {
	router *gin.Engine
}

func NewApplication(sources config.SourceConfig, athletes meets.AthleteLookup, eventStream raceevents.EventStream, logger *slog.Logger) Application {
	router := gin.Default()

	// Setup route group for the API
	api := router.Group("/api")
	api.GET("/timingEvents", handler.NewVerifyTimingHandler(logger))
	api.POST("/timingEvents/finishes", handler.NewTimingHandler(sources, athletes, eventStream, logger))

	// results paths
	router.StaticFile("/overall", "overall_results.html")

	return &application{
		router: router,
	}
}

func (a *application) Run(address string) {
	a.router.Run(address)
}
