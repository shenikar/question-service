package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// NewLogger инициализирует и возвращает новый экземпляр Logrus.
func NewLogger() *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(logrus.InfoLevel) // По умолчанию уровень INFO

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		if level, err := logrus.ParseLevel(logLevel); err == nil {
			log.SetLevel(level)
		} else {
			log.Warnf("Invalid LOG_LEVEL '%s', defaulting to INFO", logLevel)
		}
	}

	return log
}
