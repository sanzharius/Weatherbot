package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type config struct {
	MyToken        string `env:"MYTOKEN"`
	Port           string `env:"PORT"`
	WeatherApiHost string `env:"WEATHERAPIHOST"`
	AppId          string `env:"APPID"`
}

type List struct {
	Main    Main
	Weather Weather
}

type Weather struct {
	Id          int    `json:"id,omitempty"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func Init() (*config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	MyToken := os.Getenv("MYTOKEN")
	fmt.Println(MyToken)

	Port := os.Getenv("PORT")
	fmt.Println(Port)

	WeatherApiHost := os.Getenv("WEATHERAPIHOST")
	fmt.Println(WeatherApiHost)

	AppId := os.Getenv("APPID")
	fmt.Println(AppId)

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Printf("%+v\n", cfg)
	return &cfg, nil

}
