package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
	"os"
	"strings"
	"telegrambot/sanzhar/config"
	"telegrambot/sanzhar/httpclient"
	"testing"
)

type tbFakeService func(r *http.Request) (*http.Response, error)

func (f tbFakeService) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestNewBot(t *testing.T) {
	apiToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if apiToken == "" {
		t.Skip("TELEGRAM_BOT_TOKEN not set, skipping test")
	}

	testHttp := httpclient.NewHTTPCLient()
	testConfig := &config.Config{
		TelegramBotTok: apiToken,
	}
	if _, err := NewBot(testConfig, testHttp); err != nil {
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

	client := &httpclient.WeatherClient{
		Http: &http.Client{
			Transport: tbFakeService(func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Header: http.Header{
						"content-type": []string{"application/json"},
					},
					Body: io.NopCloser(strings.NewReader("location: successfully transmitted")),
				}, nil
			}),
		}}

	tbClient := &Bot{
		weatherClient: client,
	}

	var update *tgbotapi.Update
	got, err := tbClient.GetMessageByUpdate(update)
	if err != nil {
		t.Fatal(err)
	}

	want := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	if got != &want {
		t.Errorf("Unexpected result returned. Got %v, want %v", got, want)
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
