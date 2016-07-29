package sdk

import (
	"net/http"
	"time"

	"github.com/smartystreets/clock"
)

type RetryClient struct {
	inner      httpClient
	maxRetries int
	sleeper    *clock.Sleeper
}

func NewRetryClient(client httpClient, maxRetries int) *RetryClient {
	return &RetryClient{
		inner:      client,
		maxRetries: maxRetries,
	}
}

func (r *RetryClient) Do(request *http.Request) (response *http.Response, err error) {
	for attempt := 0; attempt <= r.maxRetries; attempt++ {
		if attempt > 0 {
			r.sleeper.Sleep(time.Second * time.Duration(attempt))
		}

		response, err = r.inner.Do(request)

		if err == nil && response.StatusCode == http.StatusOK {
			break
		}
	}
	return response, err
}
