package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type emptyResponse struct {
}

func NewVerifyTimingHandler(logger *slog.Logger) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		logger.Info("handling get timing events")
		c.IndentedJSON(http.StatusNoContent, emptyResponse{})
	}
	return gin.HandlerFunc(fn)
}
