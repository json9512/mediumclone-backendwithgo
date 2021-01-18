package logger

import (
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// InitLogger ...
// Returns a formatted logger
func InitLogger() *logrus.Logger {
	log := logrus.StandardLogger()

	log.SetFormatter(&formatter.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"category"},
	})

	return log
}

// MiddleWare ...
// logs the requests
func MiddleWare(log *logrus.Logger) gin.HandlerFunc {
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.TextFormatter{})
	return func(c *gin.Context) {
		startTime := time.Now()
		// Process request
		c.Next()
		// end time
		endTime := time.Now()

		// Process request info
		latency := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqURL := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		// format log based on status Code

		if statusCode == 404 {
			log.Warnf(
				"| %3d | %13v | %15s | %s | %s",
				statusCode,
				latency,
				clientIP,
				reqMethod,
				reqURL,
			)
		} else {
			log.Infof(
				"| %3d | %13v | %15s | %s | %s",
				statusCode,
				latency,
				clientIP,
				reqMethod,
				reqURL,
			)
		}
	}
}
