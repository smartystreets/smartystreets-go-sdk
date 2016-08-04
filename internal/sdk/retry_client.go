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
	for attempt := 0; r.backOff(attempt); attempt++ {
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