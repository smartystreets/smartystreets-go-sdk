package sdk

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"bitbucket.org/smartystreets/smartystreets-go-sdk"
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
		return nil, smarty_sdk.StatusBadRequest
	case 401:
		return nil, smarty_sdk.StatusUnauthorized
	case 402:
		return nil, smarty_sdk.StatusPaymentRequired
	case 413:
		return nil, smarty_sdk.StatusRequestEntityTooLarge
	case 429:
		return nil, smarty_sdk.StatusTooManyRequests
	case 200:
		return content, response.Body.Close()
	default:
		response.Body.Close()
		return nil, fmt.Errorf("Non-200 status: %s\n%s", response.Status, string(content))
	}
}
