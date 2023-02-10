package server

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func tracing() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// get request id from header
		nextRequestID := func() string {
			return fmt.Sprintf("%d", time.Now().UnixNano())
		}
		requestID := c.Request.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = nextRequestID()
		}
		c.Writer.Header().Set("X-Request-Id", requestID)
		c.Next()
	})
}

