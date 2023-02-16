package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type Log struct {
	Logger *log.Logger
}

func NewLog() *Log {
	log.WithFields(log.Fields{
		"out":  os.Stderr,
		"time": time.Now(),
	}).Info("A new message received")

	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	LogLevel, err := log.ParseLevel(os.Getenv("LOGLEVEL"))
	if err != nil {
		LogLevel = log.InfoLevel
	}

	log.SetLevel(LogLevel)
	return &Log{}

}
