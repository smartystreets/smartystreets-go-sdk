package sdk

import (
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
	for attempt := 0; attempt <= r.maxRetries; attempt++ {
		r.sleeper.Sleep(time.Second * time.Duration(attempt)) // FUTURE: upper threshold for back-off

		response, err = r.inner.Do(request)

		if err == nil && response.StatusCode == http.StatusOK {
			break
		}
	}
	return response, err
}
