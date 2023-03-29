package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"telegrambot/sanzhar/apperrors"
	"telegrambot/sanzhar/config"
	"telegrambot/sanzhar/httpclient"
)

type Bot struct {
	cfg           *config.Config
	weatherClient *httpclient.WeatherClient
	tgClient      *tgbotapi.BotAPI
}

func NewBot(config *config.Config, tgClient *tgbotapi.BotAPI, weatherClient *httpclient.WeatherClient) (*Bot, error) {
	log.Printf("Authorized on account %s", tgClient.Self.UserName)
	return &Bot{
		cfg:           config,
		weatherClient: weatherClient,
		tgClient:      tgClient,
	}, nil
}

func (bot *Bot) ReplyingOnMessages() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = bot.cfg.TelegramMessageTimeoutInSec

	updates := bot.tgClient.GetUpdatesChan(u)
	for update := range updates {
		msg, err := bot.GetMessageByUpdate(&update)
		if err != nil {
			log.Error(err)
			continue
		}

		_, err = bot.tgClient.Send(msg)
		if err != nil {
			log.Error(err)
		}
	}
}

func (bot *Bot) GetMessageByUpdate(update *tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
	if update.Message == nil {
		return nil, nil
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	if update.Message.Location != nil {
		getWeatherResponse, err := bot.weatherClient.GetWeatherForecast(update.Message.Location)
		if err != nil {
			return nil, apperrors.MessageUnmarshallingError.AppendMessage(err)
		}

		msg.Text = MapGetWeatherResponseToHTML(getWeatherResponse)
		msg.ParseMode = "HTML"
	}

	return &msg, nil
}

func MapGetWeatherResponseToHTML(list *httpclient.GetWeatherResponse) string {

	message := "<b>%s</b>: <b>%.2fdegC</b>\n" + "Feels like <b>%.2fdegC</b>. %s\n"

	reply := fmt.Sprintf(message, list.Name, list.Main.Temp, list.Main.Temp, list.Weather[0].Description)

	return reply
}
