package handler

import (
	"blreynolds4/event-race-timer/internal/meets"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewMeetListHandler(meetReader meets.MeetReader, logger *slog.Logger) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		logger.Info("handling meet list")
		meets, err := meetReader.GetMeets()
		if err != nil {
			logger.Error("error getting meets", "error", err)
			c.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, meets)
	}
	return gin.HandlerFunc(fn)
}
