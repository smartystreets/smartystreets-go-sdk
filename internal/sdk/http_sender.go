package sdk

import (
	"io/ioutil"
	"net/http"

	"github.com/smartystreets/smartystreets-go-sdk"
)

// HTTPClient matches http.Client and allows us to define custom clients that wrap over http.Client.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// HTTPSender translates *http.Request into ([]byte, error) by
// - calling the provided HTTPClient,
// - reading the response body, and
// - interpreting the content (or error) received.
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
	// TODO: Since we already copy response.Body in retry_client.go -> readBody()
	//       It would behoove us to prevent a second copy in that case.

	if content, err := ioutil.ReadAll(response.Body); err != nil {
		_ = response.Body.Close()
		return nil, err
	} else {
		return content, response.Body.Close()
	}
}

func interpret(response *http.Response, content []byte) ([]byte, error) {
	if response.StatusCode == http.StatusOK {
		return content, nil
	}
	return nil, sdk.NewHTTPStatusError(response.StatusCode, content)
}
