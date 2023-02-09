package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"telegrambot/sanzhar/config"
	"telegrambot/sanzhar/httpclient"
)

func main() {

	logger := NewLog
	log.Println(logger)
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.MyToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		switch update.Message.Location {
		case update.Message.Location:
			if update.Message.Location != nil {
				list, err := httpclient.GetWeatherForecast(update.Message.Location)
				if err != nil {
					log.Fatal("unable to parse", err)
				}
				msg.Text = Markdown(list)
				msg.ParseMode = "HTML"
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			}
		}
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil))

}

func Markdown(list *httpclient.WeatherResponse) string {
	message := "<b>%s</b>: <b>%.2fdegC</b>\n" + "Feels like <b>%.2fdegC</b>. %s\n"

	reply := fmt.Sprintf(message, list.Name, list.Main.Temp, list.Main.Temp, list.Weather[0].Description)

	return reply
}
