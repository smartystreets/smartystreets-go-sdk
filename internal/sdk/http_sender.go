package sdk

import (
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
		// BUG: if we get a response (and it has a body), it should always be closed
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		return nil, fmt.Errorf("Non-200 status: %s\n%s", response.Status, string(content))
	}

	return content, response.Body.Close()

}
