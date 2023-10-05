package sdk

import (
	"bytes"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// RetryClient sends failed requests multiple times depending on the parameters passed to NewRetryClient.
type RetryClient struct {
	inner      HTTPClient
	maxRetries int
	sleeper    func(time.Duration)
	lock       *sync.Mutex
	rand       *rand.Rand
}

func NewRetryClient(inner HTTPClient, maxRetries int, rand *rand.Rand, sleeper func(time.Duration)) HTTPClient {
	if maxRetries == 0 {
		return inner
	}
	return &RetryClient{
		inner:      inner,
		maxRetries: maxRetries,
		lock:       new(sync.Mutex),
		rand:       rand,
		sleeper:    sleeper,
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
			if r.readBody(response) {
				break
			}
		}
		if !r.handleHttpStatusCode(response, &attempt) {
			break
		}
	}
	return response, err
}

func (r *RetryClient) doBufferedPost(request *http.Request) (response *http.Response, err error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	for attempt := 0; r.backOff(attempt); attempt++ {
		request.Body = io.NopCloser(bytes.NewReader(body))
		if response, err = r.inner.Do(request); err == nil && response.StatusCode == http.StatusOK {
			if r.readBody(response) {
				break
			}
		}
		if !r.handleHttpStatusCode(response, &attempt) {
			break
		}
	}
	return response, err
}

func (r *RetryClient) handleHttpStatusCode(response *http.Response, attempt *int) bool {
	if response == nil {
		return true
	}
	if response.StatusCode >= http.StatusBadRequest && response.StatusCode <= http.StatusUnprocessableEntity {
		return false
	}
	if response.StatusCode == http.StatusTooManyRequests {
		r.sleeper(time.Second * time.Duration(r.random(backOffRateLimit)))
		// Setting attempt to 1 will make 429s retry indefinitely; this is intended behavior.
		*attempt = 1
	}
	return true
}

func (r *RetryClient) readBody(response *http.Response) bool {
	if response.Body == nil {
		return false
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, response.Body); err == nil {
		_ = response.Body.Close()
		response.Body = io.NopCloser(&buf)
		return true
	}
	_ = response.Body.Close()
	return false
}

func (r *RetryClient) backOff(attempt int) bool {
	if attempt == 0 {
		return true
	}
	if attempt > r.maxRetries {
		return false
	}
	backOffCap := max(0, min(maxBackOffDuration, attempt))
	backOff := time.Second * time.Duration(r.random(backOffCap))
	r.sleeper(backOff)
	return true
}

func (r *RetryClient) random(cap int) int {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.rand.Intn(cap)
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

const (
	backOffRateLimit   = 5
	maxBackOffDuration = 10
)
