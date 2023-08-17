package sdk

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"
)

// RetryClient sends failed requests multiple times depending on the parameters passed to NewRetryClient.
type RetryClient struct {
	inner      HTTPClient
	maxRetries int
	sleeper    func(time.Duration)
	rateLimit  int
}

func NewRetryClient(inner HTTPClient, maxRetries int, sleeper func(time.Duration)) HTTPClient {
	if maxRetries == 0 {
		return inner
	}
	return &RetryClient{
		inner:      inner,
		maxRetries: maxRetries,
		sleeper:    sleeper,
		rateLimit:  -1,
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
		if response.Header != nil {
			if i, err := strconv.Atoi(response.Header.Get("Retry-After")); err == nil {
				r.rateLimit = i
				*attempt = 0
				return true
			}
		}
		if *attempt == 0 {
			r.rateLimit = 1
		} else {
			r.rateLimit += 1
		}
		*attempt = 0
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
	if attempt > r.maxRetries {
		return false
	}
	backOffCap := 0
	if r.rateLimit != -1 {
		backOffCap = r.rateLimit
	} else {
		backOffCap = max(0, min(maxBackOffDuration, attempt))
	}
	backOff := time.Second * time.Duration(backOffCap)
	r.sleeper(backOff)
	return true
}

// TODO: delete in favor of built-in function after upgrading to Go 1.21
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// TODO: delete in favor of built-in function after upgrading to Go 1.21
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

const (
	maxBackOffDuration = 10
)
