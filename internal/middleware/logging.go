package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process the request
		c.Next()

		// Log information after the request is processed
		end := time.Now()
		latency := end.Sub(start)

		// Log the request details
		log.Infof("%s %s %s %s %d %v\n",
			end.Format("2006-01-02 15:04:05"),
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			latency,
		)
	}
}
