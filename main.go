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

	/*bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotTok)

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
				weatherClient := httpclient.NewWeatherClient(cfg)
				list, err := weatherClient.GetWeatherForecast(update.Message.Location)
				if err != nil {
					log.Fatal(apperrors.MessageUnmarshallingError.AppendMessage(err))
				}
				msg.Text = Markdown(list)
				msg.ParseMode = "HTML"
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			}
		}
	}*/
	/*log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil))*/
}

/*func Markdown(list *httpclient.GetWeatherResponse) string {

	message := "<b>%s</b>: <b>%.2fdegC</b>\n" + "Feels like <b>%.2fdegC</b>. %s\n"

	reply := fmt.Sprintf(message, list.Name, list.Main.Temp, list.Main.Temp, list.Weather[0].Description)

	return reply
}*/
