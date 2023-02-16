package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"os"
	"telegrambot/sanzhar/apperrors"
)

type Config struct {
	TelegramBotTok string `env:"TELEGRAM_BOT_TOKEN"`
	Port           string `env:"PORT"`
	WeatherApiHost string `env:"WEATHERAPIHOST"`
	AppId          string `env:"APPID"`
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, apperrors.ConfigReadErr.AppendMessage(err)
	}

	port := os.Getenv("PORT")

	weatherApiHost := os.Getenv("WEATHERAPIHOST")

	fmt.Printf("Port: %s; WeatherApiHost: %s", port, weatherApiHost)

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, apperrors.ConfigReadErr.AppendMessage(err)
	}

	fmt.Printf("%+v\n", cfg)
	return &cfg, nil

}
