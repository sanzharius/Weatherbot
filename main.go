package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {

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

	cfg, err := Init()
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

		/*msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)*/
		q := url.Values{}

		switch update.Message.Text {
		case "":
			if update.Message.Location != nil {
				weather := BuildURL(*update.Message.Location)
				body := HTTPGet(weather)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, body)
				msg.Text = Markdown(body)
				q.Add("text", msg.Text)
				q.Add("parse_mode", "HTML")
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			}

			port := os.Getenv("PORT")
			if port == "" {
				port = "8080"
			}
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

		}
	}
}

func Markdown(_ string) string {
	var list List
	var reply strings.Builder
	_, err := fmt.Fprintf(&reply, "<b>%s</b>: <b>%.2fdegC</b>\n", list.Name, list.Main.Temp)
	if err != nil {
		log.Fatal(err)
	}
	wet, err := fmt.Fprintf(&reply, "Feels like <b>%.2fdegC</b>. %s\n", list.Main.Temp, list.Weather[0].Description)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(wet)

	return reply.String()
}
