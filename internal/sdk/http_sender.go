package sdk

import (
	"io/ioutil"
	"net/http"

	"bitbucket.org/smartystreets/smartystreets-go-sdk"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type HTTPSender struct {
	client HTTPClient
}

func NewHTTPSender(client HTTPClient) *HTTPSender {
	return &HTTPSender{client: client}
}

func (s *HTTPSender) Send(request *http.Request) ([]byte, error) {
	if response, err := s.client.Do(request); err != nil {
		return nil, err
	} else if content, err := readResponseBody(response); err != nil {
		return content, err
	} else {
		return interpret(response, content)
	}
}

func readResponseBody(response *http.Response) ([]byte, error) {
	if content, err := ioutil.ReadAll(response.Body); err != nil {
		response.Body.Close()
		return nil, err
	} else {
		err = response.Body.Close()
		return content, err
	}
}

func interpret(response *http.Response, content []byte) ([]byte, error) {
	switch response.StatusCode {
	case http.StatusOK:
		return content, nil
	case http.StatusBadRequest:
		return nil, smarty_sdk.StatusBadRequest
	case http.StatusUnauthorized:
		return nil, smarty_sdk.StatusUnauthorized
	case http.StatusPaymentRequired:
		return nil, smarty_sdk.StatusPaymentRequired
	case http.StatusRequestEntityTooLarge:
		return nil, smarty_sdk.StatusRequestEntityTooLarge
	case http.StatusTooManyRequests:
		return nil, smarty_sdk.StatusTooManyRequests
	default:
		return nil, smarty_sdk.StatusUncataloguedError(response.Status, content)
	}
}
