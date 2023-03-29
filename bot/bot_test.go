package bot

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strings"
	"telegrambot/sanzhar/config"
	"telegrambot/sanzhar/httpclient"
	"testing"
)

type testTransport func(r *http.Request) (*http.Response, error)

func (t testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t(req)
}

func fakeBotWithWeatherClient(weatherClient *httpclient.WeatherClient) *Bot {
	apiToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	testConfig := &config.Config{
		TelegramBotTok: apiToken,
	}

	response, err := generateBotOkJsonApiResponse()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	testTgClientHttp := fakeHTTPBotClient(200, response)
	tgClient, err := tgbotapi.NewBotAPIWithClient(testConfig.TelegramBotTok, "https://api.telegram.org/bot%s/%s", testTgClientHttp)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := NewBot(testConfig, tgClient, weatherClient)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return bot
}

func fakeHTTPBotClient(statusCode int, jsonResponse string) *http.Client {
	return &http.Client{
		Transport: testTransport(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: statusCode,
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(jsonResponse)),
			}, nil
		}),
	}
}

func generateBotOkJsonApiResponse() (string, error) {
	testUser := tgbotapi.User{
		ID:        123,
		FirstName: "John",
		LastName:  "Doe",
	}

	user, err := json.Marshal(&testUser)
	if err != nil {
		return "", err
	}

	testApiResponse := tgbotapi.APIResponse{
		Ok:          true,
		Result:      user,
		ErrorCode:   0,
		Description: "",
		Parameters:  nil,
	}

	response, err := json.Marshal(&testApiResponse)
	if err != nil {
		return "", err
	}

	return string(response), nil
}

func TestNewBot(t *testing.T) {
	apiToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	testConfig := &config.Config{
		TelegramBotTok: apiToken,
	}

	response, err := generateBotOkJsonApiResponse()
	if err != nil {
		t.Error(err)
		return
	}

	testTgClientHttp := fakeHTTPBotClient(200, response)
	tgClient, err := tgbotapi.NewBotAPIWithClient(testConfig.TelegramBotTok, "https://api.telegram.org/bot%s/%s", testTgClientHttp)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := NewBot(testConfig, tgClient, nil); err != nil {
		t.Log(err)
	}
}

func TestReplyingOnMessages(t *testing.T) {
	tbClient := &Bot{}
	testU := tgbotapi.NewUpdate(0)
	testU.Timeout = tbClient.cfg.TelegramMessageTimeoutInSec
	testUpdates := tbClient.tgClient.GetUpdatesChan(testU)
	for testUpdate := range testUpdates {
		testMsg, err := tbClient.GetMessageByUpdate(&testUpdate)
		if err != nil {
			t.Log(err)
			continue
		}

		_, err = tbClient.tgClient.Send(testMsg)
		if err != nil {
			t.Log(err)
		}
	}
}

func TestGetMessageByUpdate(t *testing.T) {

	ttPass := []struct {
		name                 string
		messageChatId        int64
		givenWeatherResponse *httpclient.GetWeatherResponse
		givenMessage         *tgbotapi.Update
		botReply             *tgbotapi.MessageConfig
	}{
		{
			"existing location command",
			300,
			&httpclient.GetWeatherResponse{
				Name: "Warsaw",
				Main: &httpclient.MainForecast{
					Temp:      10,
					FeelsLike: 15,
				},
				Weather: []*httpclient.Weather{
					{Description: "It's Warsaw, baby :)"},
				},
			},
			&tgbotapi.Update{
				UpdateID: 0,
				Message: &tgbotapi.Message{
					Chat: &tgbotapi.Chat{
						ID: 300,
					},
					Location: &tgbotapi.Location{
						Longitude: 21.017532,
						Latitude:  52.237049,
					},
				},
			},
			&tgbotapi.MessageConfig{
				Text:      "",
				ParseMode: "HTML",
			},
		},
	}

	testConfig, err := config.NewConfig("../.env")
	if err != nil {
		log.Fatal(err)
	}

	for _, tc := range ttPass {

		responseJSON, err := json.Marshal(tc.givenWeatherResponse)
		if err != nil {
			t.Fatal(err)
		}

		httpClient := fakeHTTPBotClient(200, string(responseJSON))
		weatherClient := httpclient.NewWeatherClient(testConfig, httpClient)
		bot := fakeBotWithWeatherClient(weatherClient)
		bot.tgClient.Debug = true
		msg, err := bot.GetMessageByUpdate(tc.givenMessage)
		if err != nil {
			t.Error(err)
		}

		expMessage := MapGetWeatherResponseToHTML(tc.givenWeatherResponse)
		if msg != nil && msg.Text != expMessage {
			t.Errorf("bot reply should be %s, but got %s", tc.botReply.Text, msg.Text)
		}
	}
}

func TestMapGetWeatherResponseToHTML(t *testing.T) {
	var testList *httpclient.GetWeatherResponse

	MapGetWeatherResponseToHTML(testList)
	want := "<b>%s</b>: <b>%.2fdegC</b>\n" + "Feels like <b>%.2fdegC</b>. %s\n"

	got := fmt.Sprintf(want, testList.Name, testList.Main.Temp, testList.Main.Temp, testList.Weather[0].Description)

	if got != want {
		t.Errorf("Unexpected result returned. Got %v, want %v", got, want)
	}

}
