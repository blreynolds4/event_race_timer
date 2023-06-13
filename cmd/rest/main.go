package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// methods to handle rest requests
func createTimingEvent(c *gin.Context) {
	fmt.Println("Handling post timing events")
	c.String(http.StatusOK, "implement create timing event")
}

func getLiveTimingEventsEvents(c *gin.Context) {
	fmt.Println("Handling get live timing events")
	c.String(http.StatusOK, "implement get live timing event")
}

func getReplayTimingEventsEvents(c *gin.Context) {
	fmt.Println("Handling get replay timing events")
	c.String(http.StatusOK, "implement get replay timing event")
}

func main() {
	router := gin.Default()

	// Setup route group for the API
	api := router.Group("/api")
	api.GET("/liveTimingEvents", getLiveTimingEventsEvents)
	api.GET("/replayTimingEvents", getReplayTimingEventsEvents)
	api.POST("/timingEvents", createTimingEvent)

	router.Run("localhost:8080")
}
