package httpclient

import (
	"net/http"
	"time"
)

func NewHTTPCLient() *http.Client {
	return &http.Client{
		Timeout: time.Minute * 10,
	}
}
