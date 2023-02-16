package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
	"telegrambot/sanzhar/apperrors"
	"telegrambot/sanzhar/config"
	"telegrambot/sanzhar/httpclient"
	"telegrambot/sanzhar/logger"
)

type Bot struct {
	cfg           *config.Config
	weatherClient *httpclient.WeatherClient
	log           *logger.Log
}

func NewBot(config *config.Config, weatherClient *httpclient.WeatherClient, log *logger.Log) *Bot {
	return &Bot{
		cfg:           config,
		weatherClient: weatherClient,
		log:           log,
	}
}

func (bot *Bot) ReplyingOnMessages() {
	b, err := tgbotapi.NewBotAPI(bot.cfg.TelegramBotTok)
	if err != nil {
		bot.log.Logger.Panic(err)
	}
	b.Debug = true
	bot.log.Logger.Printf("Authorized on account %s", b.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		switch update.Message.Location {
		case update.Message.Location:
			if update.Message.Location != nil {
				list, err := bot.weatherClient.GetWeatherForecast(update.Message.Location)
				if err != nil {
					bot.log.Logger.Fatal(apperrors.MessageUnmarshallingError.AppendMessage(err))
				}
				msg.Text = Markdown(list)
				msg.ParseMode = "HTML"
				if _, err := b.Send(msg); err != nil {
					bot.log.Logger.Panic(err)
				}
			}
		}
	}
	bot.log.Logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", bot.cfg.Port), nil))

}

func Markdown(list *httpclient.GetWeatherResponse) string {

	message := "<b>%s</b>: <b>%.2fdegC</b>\n" + "Feels like <b>%.2fdegC</b>. %s\n"

	reply := fmt.Sprintf(message, list.Name, list.Main.Temp, list.Main.Temp, list.Weather[0].Description)

	return reply
}
