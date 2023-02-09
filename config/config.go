package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"os"
	"telegrambot/sanzhar/apperrors"
)

type Config struct {
	MyToken        string `env:"MYTOKEN"`
	Port           string `env:"PORT"`
	WeatherApiHost string `env:"WEATHERAPIHOST"`
	AppId          string `env:"APPID"`
}

const errMsg = "parse failed"

func Init() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, apperrors.WrapNil(errMsg, err)
	}

	MyToken := os.Getenv("MYTOKEN")
	fmt.Println(MyToken)

	Port := os.Getenv("PORT")
	fmt.Println(Port)

	WeatherApiHost := os.Getenv("WEATHERAPIHOST")
	fmt.Println(WeatherApiHost)

	AppId := os.Getenv("APPID")
	fmt.Println(AppId)

	Cfg := Config{}
	if err := env.Parse(&Cfg); err != nil {
		return nil, apperrors.WrapNil(errMsg, err)
	}

	fmt.Printf("%+v\n", Cfg)
	return &Cfg, nil

}
