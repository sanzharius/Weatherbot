package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"os"
	"telegrambot/sanzhar/apperrors"
)

type Config struct {
	MyToken        string `env:"TELEGRAM_BOT_TOKEN"`
	Port           string `env:"PORT"`
	WeatherApiHost string `env:"WEATHERAPIHOST"`
	AppId          string `env:"APPID"`
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, apperrors.ConfigReadErr.AppendMessage(err)
	}

	Port := os.Getenv("PORT")

	weatherApiHost := os.Getenv("WEATHERAPIHOST")

	fmt.Printf("Port: %s; WeatherApiHost: %s", Port, weatherApiHost)

	Cfg := Config{}
	if err := env.Parse(&Cfg); err != nil {
		return nil, apperrors.ConfigReadErr.AppendMessage(err)
	}

	fmt.Printf("%+v\n", Cfg)
	return &Cfg, nil

}
