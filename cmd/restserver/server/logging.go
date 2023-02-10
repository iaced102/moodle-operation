package server

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func logging() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		requestID := c.Writer.Header().Get("X-Request-Id")
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()
		agent := c.Request.UserAgent()
		fmt.Printf("%v | %3d | %13v | %15s | %s  %s | %s", end.Format("2006/01/02 - 15:04:05"), statusCode, latency, clientIP, method, path, agent)
		if statusCode >= 400 {
			fmt.Printf(" | %s", requestID)
		}
		fmt.Println()
	})
}

