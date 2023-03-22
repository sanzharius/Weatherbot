package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
	"telegrambot/sanzhar/config"
)

func InitLog(cfg *config.Config) {

	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	logLevel, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		logLevel = log.InfoLevel
	}

	log.SetLevel(logLevel)
}
