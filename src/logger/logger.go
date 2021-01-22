package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// InitLogger returns a formatted logger
func InitLogger() *logrus.Logger {
	log := logrus.StandardLogger()
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}

// Middleware logs the requests
func Middleware(log *logrus.Logger) gin.HandlerFunc {
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
