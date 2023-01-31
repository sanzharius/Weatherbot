package main

import (
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env/v6"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type config struct {
	MyToken        string `env:"MYTOKEN"`
	Port           string `env:"PORT"`
	WeatherApiHost string `env:"WEATHERAPIHOST"`
	AppId          string `env:"APPID"`
}

type List struct {
	Name    string     `json:"name"`
	Main    Main       `json:"main"`
	Wind    Wind       `json:"wind"`
	Weather []*Weather `json:"weather"`
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

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
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

func BuildURL(loc *tgbotapi.Location) (Parsed string) {

	cfg, err := Init()
	if err != nil {
		log.Fatal(err)
	}
	URL, _ := url.Parse(cfg.WeatherApiHost)

	r := url.Values{}

	r.Add("appid", cfg.AppId)
	r.Add("lat", fmt.Sprint(loc.Latitude))
	r.Add("lon", fmt.Sprint(loc.Longitude))
	r.Add("units", "metric")

	URL.RawQuery = r.Encode()
	Parsed = URL.String()
	return Parsed

}

func HTTPGet(weatherURL string) *List {

	resp, err := http.Get(weatherURL)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if resp.StatusCode == http.StatusOK {
		log.Println("request succeeded: 200 OK")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var list List
	err = json.Unmarshal(body, &list)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", list)

	return &list

}
