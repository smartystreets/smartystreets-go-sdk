package sdk

import (
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

func (this *HTTPSender) Send(request *http.Request) (content []byte, err error) {
	if response, err := this.client.Do(request); err != nil {
		return nil, err
	} else if content, err = ioutil.ReadAll(response.Body); err != nil {
		return nil, err
	} else {
		return content, response.Body.Close()
	}
}
