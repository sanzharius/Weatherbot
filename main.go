package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"telegrambot/sanzhar/bot"
	"telegrambot/sanzhar/config"
	"telegrambot/sanzhar/httpclient"
	"telegrambot/sanzhar/logger"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger.InitLog(cfg)
	httpClient := httpclient.NewHTTPCLient()

	tgBot, err := bot.NewBot(cfg, httpClient)
	if err != nil {
		log.Fatal(err)
	}

	tgBot.ReplyingOnMessages()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil))
}
