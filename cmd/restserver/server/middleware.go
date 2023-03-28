package server

import (
	"moodle/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(tenantName, token string) bool {
	if tenantName != config.USER {
		return false
	}

	if token != config.MOODLE_TOKEN {
		return false
	}
	return true
}

func MoodleMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// ignore for health and /api/v1/swagger/*any
		if c.Request.URL.Path == "/health" || strings.HasPrefix(c.Request.URL.Path, "/api/v1/swagger") {
			c.Next()
			return
		}

		// handle cors
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Tenant-Name, X-Auth-Token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// set response header
		c.Header("Content-Type", "application/json")

		// handle preflight request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		tenentName := c.Request.Header.Get("X-Tenant-Name")
		token := c.Request.Header.Get("X-Auth-Token")
		if !Auth(tenentName, token) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	})
}


