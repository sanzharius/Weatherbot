package main

import (
	tglog "github.com/sirupsen/logrus"
	"telegrambot/sanzhar/bot"
	"telegrambot/sanzhar/config"
	"telegrambot/sanzhar/httpclient"
	"telegrambot/sanzhar/logger"
)

func main() {

	log := logger.NewLog()
	cfg, err := config.NewConfig()
	if err != nil {
		tglog.Fatal(err)
	}

	httpClient := httpclient.NewHTTPCLient()
	weatherClient := httpclient.NewWeatherClient(cfg, httpClient)

	tgBot := bot.NewBot(cfg, weatherClient, log) /*bot.NewBot(cfg, weatherClient, botLogger)*/
	tgBot.ReplyingOnMessages()

}
