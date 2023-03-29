package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"telegrambot/sanzhar/bot"
	"telegrambot/sanzhar/config"
	"telegrambot/sanzhar/httpclient"
	"telegrambot/sanzhar/logger"
)

func main() {
	cfg, err := config.NewConfig(".env")
	if err != nil {
		log.Fatal(err)
	}

	logger.InitLog(cfg)

	httpWeatherClient := httpclient.NewHTTPCLient()
	weatherClient := httpclient.NewWeatherClient(cfg, httpWeatherClient)

	httpTgClient := httpclient.NewHTTPCLient()
	tgClient, err := tgbotapi.NewBotAPIWithClient(cfg.TelegramBotTok, "https://api.telegram.org/bot%s/%s", httpTgClient)
	if err != nil {
		log.Fatal(err)
	}
	tgClient.Debug = true

	tgBot, err := bot.NewBot(cfg, tgClient, weatherClient)
	if err != nil {
		log.Fatal(err)
	}

	tgBot.ReplyingOnMessages()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil))
}
