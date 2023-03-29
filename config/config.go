package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"os"
	"telegrambot/sanzhar/apperrors"
)

type Config struct {
	TelegramHost                string `env:"TELEGRAM_HOST"`
	TelegramBotTok              string `env:"TELEGRAM_BOT_TOKEN"`
	Port                        string `env:"PORT"`
	WeatherApiHost              string `env:"WEATHERAPIHOST"`
	AppId                       string `env:"APPID"`
	LogLevel                    string `env:"LOGLEVEL"`
	TelegramMessageTimeoutInSec int    `env:"TELEGRAM_MESSAGE_TIMEOUT_IN_SEC"`
}

func NewConfig(path string) (*Config, error) {
	err := godotenv.Load(path)
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
