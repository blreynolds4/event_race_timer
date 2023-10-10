package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type emptyResponse struct {
}

type OpenSignupsTimingEvent struct {
	ElapsedTime int    `json:"elapsedTime"`
	CaptureMode string `json:"captureMode"`
	StartTime   int    `json:"startTime"`
	Antenna     int    `json:"antenna"`
	Bib         int    `json:"bib"`
}

// methods to handle rest requests
func verityTimingEvent(c *gin.Context) {
	fmt.Println("Handling get timing events")
	c.IndentedJSON(http.StatusOK, emptyResponse{})
	// c.String(http.StatusNoContent, "ok")
}

func createTimingEvent(c *gin.Context) {
	postData, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println("Handling post timing events", string(postData))
	var data OpenSignupsTimingEvent
	if err := c.BindJSON(&data); err != nil {
		return
	}
	fmt.Println("Got bib elapsed start antenna", data.Bib, data.ElapsedTime, data.StartTime, data.Antenna)
	c.IndentedJSON(http.StatusCreated, data)
}

func main() {
	router := gin.Default()

	// Setup route group for the API
	api := router.Group("/api")
	api.GET("/timingEvents", verityTimingEvent)
	api.POST("/timingEvents/finishes", createTimingEvent)

	router.Run(":8080")
}
