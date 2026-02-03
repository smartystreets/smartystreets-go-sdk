package sdk

import (
	"bytes"
	"context"
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
	sleeper    func(context.Context, time.Duration)
	lock       *sync.Mutex
	rand       *rand.Rand
}

func NewRetryClient(inner HTTPClient, maxRetries int, rand *rand.Rand, sleeper func(context.Context, time.Duration)) HTTPClient {
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
	ctx := request.Context()
	for attempt := 0; r.backOff(ctx, attempt); attempt++ {
		if err = ctx.Err(); err != nil {
			return nil, err
		}
		if response, err = r.inner.Do(request); err == nil && response.StatusCode == http.StatusOK {
			if r.readBody(response) {
				break
			}
		}
		if !r.handleHttpStatusCode(ctx, response, &attempt) {
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

	ctx := request.Context()
	for attempt := 0; r.backOff(ctx, attempt); attempt++ {
		if err = ctx.Err(); err != nil {
			return nil, err
		}
		request.Body = io.NopCloser(bytes.NewReader(body))
		if response, err = r.inner.Do(request); err == nil && response.StatusCode == http.StatusOK {
			if r.readBody(response) {
				break
			}
		}
		if !r.handleHttpStatusCode(ctx, response, &attempt) {
			break
		}
	}
	return response, err
}

func (r *RetryClient) handleHttpStatusCode(ctx context.Context, response *http.Response, attempt *int) bool {
	if response == nil {
		return true
	}
	if response.StatusCode >= http.StatusBadRequest && response.StatusCode <= http.StatusUnprocessableEntity {
		return false
	}
	if response.StatusCode == http.StatusNotModified {
		return false
	}
	if response.StatusCode == http.StatusTooManyRequests {
		r.sleeper(ctx, time.Second*time.Duration(r.random(backOffRateLimit)))
		if ctx.Err() != nil {
			return false
		}
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

func (r *RetryClient) backOff(ctx context.Context, attempt int) bool {
	if attempt == 0 {
		return true
	}
	if attempt > r.maxRetries {
		return false
	}
	backOffCap := max(0, min(maxBackOffDuration, attempt))
	backOff := time.Second * time.Duration(r.random(backOffCap))
	r.sleeper(ctx, backOff)
	return true
}

// ContextSleep is a context-aware sleep function suitable for production use.
// It returns immediately if the context is cancelled.
func ContextSleep(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}

func (r *RetryClient) random(cap int) int {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.rand.Intn(cap)
}

const (
	backOffRateLimit   = 5
	maxBackOffDuration = 10
)
