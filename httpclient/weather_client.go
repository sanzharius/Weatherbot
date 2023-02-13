package httpclient

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"telegrambot/sanzhar/apperrors"
	"telegrambot/sanzhar/config"
)

type GetWeatherResponse struct {
	Name    string        `json:"name"`
	Main    *MainForecast `json:"main"`
	Wind    *Wind         `json:"wind"`
	Weather []*Weather    `json:"weather"`
}

type Weather struct {
	Id          int    `json:"id,omitempty"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type MainForecast struct {
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

func AppendQueryParamsToGetWeather(loc *tgbotapi.Location) (Parsed string) {
	cfg, err := config.NewConfig()
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

func GetWeatherForecast(loc *tgbotapi.Location) (*GetWeatherResponse, error) {
	weatherURL := AppendQueryParamsToGetWeather(loc)
	resp, err := http.Get(weatherURL)
	if err != nil {
		return nil, apperrors.MessageUnmarshallingError.AppendMessage(err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, apperrors.DataNotFoundErr.AppendMessage(err)
		default:
			errMsg := fmt.Sprintf("Got unknown err, while calling API to get weather forecast. HTTP code: %v", resp.StatusCode)
			return nil, apperrors.APICallingErr.AppendMessage(errMsg)
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, apperrors.MessageUnmarshallingError.AppendMessage(err)
	}
	var list GetWeatherResponse
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, apperrors.MessageUnmarshallingError.AppendMessage(err)
	}
	log.Printf("%+v\n", list)
	return &list, nil
}
