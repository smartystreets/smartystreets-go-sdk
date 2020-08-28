package sdk

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/smartystreets/clock"
)

// RetryClient sends failed requests multiple times depending on the parameters passed to NewRetryClient.
type RetryClient struct {
	inner      HTTPClient
	maxRetries int
	sleeper    *clock.Sleeper
	rand       *rand.Rand
}

func NewRetryClient(inner HTTPClient, maxRetries int, rand *rand.Rand) HTTPClient {
	if maxRetries == 0 {
		return inner
	}
	return &RetryClient{
		inner:      inner,
		maxRetries: maxRetries,
		rand:       rand,
	}
}

func (r *RetryClient) Do(request *http.Request) (*http.Response, error) {
	if request.Method == "POST" {
		return r.doBufferedPost(request)
	}
	return r.doGet(request)
}

func (r *RetryClient) doGet(request *http.Request) (response *http.Response, err error) {
	for attempt := 0; r.backOff(attempt); attempt++ {
		if response, err = r.inner.Do(request); err == nil && response.StatusCode == http.StatusOK {
			break
		}
	}
	return response, err
}

func (r *RetryClient) doBufferedPost(request *http.Request) (response *http.Response, err error) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	for attempt := 0; r.backOff(attempt); attempt++ {
		request.Body = ioutil.NopCloser(bytes.NewReader(body))
		if response, err = r.inner.Do(request); err == nil && response.StatusCode == http.StatusOK {
			break
		}
	}
	return response, err
}

func (r *RetryClient) backOff(attempt int) bool {
	if attempt == 0 {
		return true
	}
	if attempt > r.maxRetries {
		return false
	}
	backOffCap := min(maxBackOffDuration, 2<<attempt)
	backOff := time.Second * time.Duration(r.rand.Intn(backOffCap))
	r.sleeper.Sleep(backOff)
	return true
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

const maxBackOffDurationInSeconds = 10
const maxBackOffDuration = 10
