package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"io"
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
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		switch update.Message.Text {
		case "":
			if update.Message.Text != nil {
				weather, _ := MakeRequest(update.Message.Text)
				fmt.Printf("message: %s\n", weather)
			}

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

func MakeRequest(location Location) (*List, error) {

	cfg, err := Init()
	if err != nil {
		log.Fatal(err)
	}

	r := url.Values{}
	r.Add("appid", cfg.AppId)
	r.Add("lat", fmt.Sprint(location.Lat))
	r.Add("lon", fmt.Sprint(location.Lon))

	resp, err := http.Get(cfg.WeatherApiHost)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))

	var list List
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal to struct: %w", err)

	}
	return &list, nil

}

func Markdown(list List) string {
	var reply strings.Builder
	fmt.Fprintf(&reply, "<b>%s</b>: <b>%.2fdegC<b>\n", list.Weather.Main, list.Main.Temp)
	fmt.Fprintf(&reply, "Feels like <b>%.2fdegC<b>. %s\n", list.Main.Temp, strings.ToTitle(list.Weather.Description))

	return reply.String()
}
