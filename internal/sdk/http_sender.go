package sdk

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type HTTPSender struct {
	client httpClient
}

func NewHTTPSender(client httpClient) *HTTPSender {
	return &HTTPSender{client: client}
}

func (s *HTTPSender) Send(request *http.Request) (content []byte, err error) {
	response, err := s.client.Do(request)
	if err != nil {
		return nil, err
	}

	content, err = ioutil.ReadAll(response.Body)
	if err != nil {
		response.Body.Close()
		return nil, err
	}

	switch response.StatusCode {
	case 400:
		return nil, BadRequest
	case 401:
		return nil, Unauthorized
	case 402:
		return nil, PaymentRequired
	case 413:
		return nil, TooLarge
	case 429:
		return nil, TooManyRequests
	case 200:
		return content, response.Body.Close()
	default:
		response.Body.Close()
		return nil, fmt.Errorf("Non-200 status: %s\n%s", response.Status, string(content))
	}
}

var (
	Unauthorized    = errors.New("401 Unauthorized")
	PaymentRequired = errors.New("402 Payment Required")
	BadRequest      = errors.New("400 Bad Request")
	TooLarge        = errors.New("413 Request entity too large")
	TooManyRequests = errors.New("429 Too many requests")
)
