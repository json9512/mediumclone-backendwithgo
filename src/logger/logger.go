package logger

import (
	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

// InitLogger ...
// Returns a formatted logger
func InitLogger() *logrus.Logger {
	log := logrus.New()

	log.SetFormatter(&formatter.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"component", "category"},
	})
	return log
}
