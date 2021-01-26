package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// CustomLogger logs the requests
func CustomLogger(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		// Process request
		c.Next()
		// end time
		endTime := time.Since(startTime)

		// Process request info
		msg := logrus.Fields{
			"latency":    endTime,
			"reqMethod":  c.Request.Method,
			"reqURL":     c.Request.RequestURI,
			"statusCode": c.Writer.Status(),
			"clientIP":   c.ClientIP(),
		}

		// format log based on status Code
		if c.Writer.Status() == 404 {
			log.WithFields(msg).Warn()
		} else {
			log.WithFields(msg).Info()
		}
	}
}
