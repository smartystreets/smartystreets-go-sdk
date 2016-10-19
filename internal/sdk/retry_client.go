package sdk

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/smartystreets/clock"
)

// RetryClient sends failed requests multiple times depending on the parameters passed to NewRetryClient.
type RetryClient struct {
	inner      HTTPClient
	maxRetries int
	sleeper    *clock.Sleeper
}

func NewRetryClient(inner HTTPClient, maxRetries int) HTTPClient {
	if maxRetries == 0 {
		return inner
	}
	return &RetryClient{
		inner:      inner,
		maxRetries: maxRetries,
	}
}

func (r *RetryClient) Do(request *http.Request) (response *http.Response, err error) {
	var body []byte
	if request.Method == "POST" {
		body, err = ioutil.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
	}

	for attempt := 0; r.backOff(attempt); attempt++ {
		if len(body) > 0 {
			request.Body = ioutil.NopCloser(bytes.NewReader(body))
		}
		if response, err = r.inner.Do(request); err == nil && response.StatusCode == http.StatusOK {
			break
		}
	}
	return response, err
}

func (r *RetryClient) backOff(attempt int) bool {
	if attempt > r.maxRetries {
		return false
	}
	backOff := time.Second * time.Duration(attempt)
	r.sleeper.Sleep(minDuration(backOff, maxBackOffDuration))
	return true
}

func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

const maxBackOffDuration = time.Second * 10
