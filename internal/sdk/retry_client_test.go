package sdk

import (
	"context"
	"errors"
	"io"
	"net/http"
	"slices"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/smarty/assertions"
	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestRetryClientFixture(t *testing.T) {
	gunit.Run(new(RetryClientFixture), t)
}

type RetryClientFixture struct {
	*gunit.Fixture
	inner    *FakeMultiHTTPClient
	response *http.Response
	err      error

	naps []time.Duration
}

func (f *RetryClientFixture) TestRequestBodyCannotBeBuffered_ErrorReturnedImmediately() {
	f.response, f.err = f.sendErrorProneRequest()
	f.assertReadErrorReturnedAndRequestNotSent()
}
func (f *RetryClientFixture) sendErrorProneRequest() (*http.Response, error) {
	f.inner = &FakeMultiHTTPClient{}
	client := NewRetryClient(f.inner, 10, f.sleep).(*RetryClient)
	request, _ := http.NewRequest("POST", "/", &ErrorProneReadCloser{readError: errors.New("GOPHERS!")})
	return client.Do(request)
}
func (f *RetryClientFixture) sleep(_ context.Context, duration time.Duration) {
	f.naps = append(f.naps, duration)
}
func (f *RetryClientFixture) assertReadErrorReturnedAndRequestNotSent() {
	f.So(f.response, should.BeNil)
	f.So(f.err, should.Resemble, errors.New("GOPHERS!"))
	f.So(f.inner.call, should.Equal, 0)
}

func (f *RetryClientFixture) TestGetRequestRetryUntilSuccess() {
	f.simulateNetworkOutageUntilSuccess()
	f.response, f.err = f.sendGetWithRetry(4)
	f.assertRequestAttempted5TimesWithBackOff_EachTimeWithSameBody()
}

/**************************************************************************/

func (f *RetryClientFixture) TestRetryOnClientErrorUntilSuccess() {
	f.simulateNetworkOutageUntilSuccess()
	f.response, f.err = f.sendPostWithRetry(4)
	f.assertRequestAttempted5TimesWithBackOff_EachTimeWithSameBody()
}
func (f *RetryClientFixture) simulateNetworkOutageUntilSuccess() {
	clientError := errors.New("Simulating Network Outage")
	f.inner = NewErringHTTPClient(clientError, clientError, clientError, clientError, nil)
}
func (f *RetryClientFixture) assertRequestAttempted5TimesWithBackOff_EachTimeWithSameBody() {
	f.assertRequestWasSuccessful()
	f.assertBackOffStrategyWasObserved()
	f.So(f.inner.bodies, should.Resemble, []string{"request", "request", "request", "request", "request"})
}
func (f *RetryClientFixture) assertRequestWasSuccessful() {
	f.So(f.err, should.BeNil)
	if f.So(f.response, should.NotBeNil) {
		f.So(f.response.StatusCode, should.Equal, 200)
	}
}
func (f *RetryClientFixture) assertBackOffStrategyWasObserved() {
	f.So(f.inner.call, should.Equal, 5)
	f.So(len(f.naps), should.Equal, 4) // 4 backoff sleeps for 5 attempts (first attempt has no backoff)
	for i, nap := range f.naps {
		cap := time.Second * time.Duration(min(i+1, maxBackOffDuration))
		f.So(nap, should.BeGreaterThanOrEqualTo, time.Second)
		f.So(nap, should.BeLessThanOrEqualTo, cap)
	}
}

/**************************************************************************/

func (f *RetryClientFixture) TestRetryOnBadResponseUntilSuccess() {
	f.inner = NewFailingHTTPClient(500, 501, 502, 522, 200)

	f.response, f.err = f.sendPostWithRetry(4)

	f.assertRequestWasSuccessful()
	f.assertBackOffStrategyWasObserved()
}

func (f *RetryClientFixture) TestPost404ErrorDoesNotRetry() {
	f.inner = NewFailingHTTPClient(404, 429)

	f.response, f.err = f.sendPostWithRetry(1)

	if f.So(f.response, should.NotBeNil) {
		f.So(f.response.StatusCode, should.Equal, 404)
	}
	f.So(f.err, should.BeNil)
}

func (f *RetryClientFixture) TestGet404ErrorDoesNotRetry() {
	f.inner = NewFailingHTTPClient(404, 429)

	f.response, f.err = f.sendGetWithRetry(1)

	if f.So(f.response, should.NotBeNil) {
		f.So(f.response.StatusCode, should.Equal, 404)
	}
	f.So(f.err, should.BeNil)
}

/**************************************************************************/

func (f *RetryClientFixture) TestFailureReturnedIfRetryExceeded() {
	f.inner = NewFailingHTTPClient(500, 500, 500, 500, 500)

	f.response, f.err = f.sendPostWithRetry(4)

	f.assertInternalServerError()
	f.assertBackOffStrategyWasObserved()
}
func (f *RetryClientFixture) assertInternalServerError() {
	if f.So(f.response, should.NotBeNil) {
		f.So(f.response.StatusCode, should.Equal, 500)
	}
	f.So(f.err, should.BeNil)
}

/**************************************************************************/

func (f *RetryClientFixture) TestNoRetryRequestedReturnsInnerClientInstead() {
	inner := &FakeHTTPClient{}
	client := NewRetryClient(inner, 0, f.sleep)
	f.So(client, should.Equal, inner)
}

/**************************************************************************/

func (f *RetryClientFixture) TestBackOffNeverToExceedHardCodedMaximum() {
	retries := 2000
	f.inner = NewFailingHTTPClient(make([]int, retries)...)

	_, f.err = f.sendPostWithRetry(retries - 1)

	f.So(f.err, should.BeNil)
	f.So(f.inner.call, should.Equal, retries)
	// The maximum must actually be reachable, not just an exclusive upper bound.
	f.So(slices.Contains(f.naps, maxBackOffDuration*time.Second), should.BeTrue)
	for i, nap := range f.naps {
		cap := time.Second * time.Duration(min(i+1, maxBackOffDuration))
		f.So(nap, should.BeGreaterThanOrEqualTo, time.Second)
		f.So(nap, should.BeLessThanOrEqualTo, cap)
	}
}

func (f *RetryClientFixture) TestBackOffRateLimitedGet() {
	retries := 4
	x := http.StatusTooManyRequests
	f.inner = NewFailingHTTPClient(x, x, x, x, x) // more 429s than retries

	response, _ := f.sendGetWithRetry(retries)

	// 429s are bounded by maxRetries
	f.So(f.inner.call, should.Equal, retries+1)
	f.So(response.StatusCode, should.Equal, http.StatusTooManyRequests)
	// Each retry triggers exactly one sleep (rate limit sleep via backOff)
	f.So(len(f.naps), should.Equal, retries)
}

func (f *RetryClientFixture) TestBackOffRateLimitedPost() {
	retries := 4
	x := http.StatusTooManyRequests
	f.inner = NewFailingHTTPClient(x, x, x, x, x) // more 429s than retries

	response, _ := f.sendPostWithRetry(retries)

	// 429s are bounded by maxRetries
	f.So(f.inner.call, should.Equal, retries+1)
	f.So(response.StatusCode, should.Equal, http.StatusTooManyRequests)
	// Each retry triggers exactly one sleep (rate limit sleep via backOff)
	f.So(len(f.naps), should.Equal, retries)
}

func (f *RetryClientFixture) TestRetryAfterHeaderUsedFor429Post() {
	f.inner = NewFailingHTTPClient(429, 429, 429, http.StatusOK)
	retryAfterSeconds := 7
	f.inner.headerKey = "Retry-After"
	f.inner.rateLimitTime = retryAfterSeconds
	f.inner.responses[3].Body = io.NopCloser(strings.NewReader("Success"))

	_, f.err = f.sendPostWithRetry(3) // 3 retries allows for 4 attempts

	f.So(f.err, should.BeNil)
	f.So(f.inner.call, should.Equal, 4)
	// At least one nap should be the exact Retry-After value (from 429 handling)
	hasRetryAfterNap := slices.Contains(f.naps, time.Second*time.Duration(retryAfterSeconds))
	f.So(hasRetryAfterNap, should.BeTrue)
}

func (f *RetryClientFixture) TestFallsBackToDefaultSleepWithoutRetryAfterHeaderPost() {
	f.inner = NewFailingHTTPClient(429, 429, 429, http.StatusOK)
	f.inner.headerKey = "X-Invalid-Header" // Not Retry-After
	f.inner.rateLimitTime = 100            // Would be obvious if used
	f.inner.responses[3].Body = io.NopCloser(strings.NewReader("Success"))

	_, f.err = f.sendPostWithRetry(3) // 3 retries allows for 4 attempts

	f.So(f.err, should.BeNil)
	f.So(f.inner.call, should.Equal, 4)
	// Without Retry-After header, each 429 uses the fixed default sleep
	for _, nap := range f.naps {
		f.So(nap, should.Equal, defaultRateLimitSleep)
	}
}

func (f *RetryClientFixture) TestRetryAfterHeaderUsedFor429Get() {
	f.inner = NewFailingHTTPClient(429, 429, 429, http.StatusOK)
	retryAfterSeconds := 7
	f.inner.headerKey = "Retry-After"
	f.inner.rateLimitTime = retryAfterSeconds
	f.inner.responses[3].Body = io.NopCloser(strings.NewReader("Success"))

	_, f.err = f.sendGetWithRetry(3) // 3 retries allows for 4 attempts

	f.So(f.err, should.BeNil)
	f.So(f.inner.call, should.Equal, 4)
	hasRetryAfterNap := slices.Contains(f.naps, time.Second*time.Duration(retryAfterSeconds))
	f.So(hasRetryAfterNap, should.BeTrue)
}

func (f *RetryClientFixture) TestFallsBackToDefaultSleepWithoutRetryAfterHeaderGet() {
	f.inner = NewFailingHTTPClient(429, 429, 429, http.StatusOK)
	f.inner.headerKey = "X-Invalid-Header" // Not Retry-After
	f.inner.rateLimitTime = 100            // Would be obvious if used
	f.inner.responses[3].Body = io.NopCloser(strings.NewReader("Success"))

	_, f.err = f.sendGetWithRetry(3) // 3 retries allows for 4 attempts

	f.So(f.err, should.BeNil)
	f.So(f.inner.call, should.Equal, 4)
	for _, nap := range f.naps {
		f.So(nap, should.Equal, defaultRateLimitSleep)
	}
}

/**************************************************************************/

func (f *RetryClientFixture) sendGetWithRetry(retries int) (*http.Response, error) {
	if len(f.inner.responses) <= retries {
		f.T().Fatalf("The number of retries is greater than or equal to the number of status codes provided. Please ensure that the number of retries is less than the number of status codes provided.")
	}

	client := NewRetryClient(f.inner, retries, f.sleep).(*RetryClient)
	request, _ := http.NewRequest("GET", "/?body=request", nil)
	return client.Do(request)
}
func (f *RetryClientFixture) sendPostWithRetry(retries int) (*http.Response, error) {
	if len(f.inner.responses) <= retries {
		f.T().Fatalf("The number of retries is greater than or equal to the number of status codes provided. Please ensure that the number of retries is less than the number of status codes provided.")
	}

	client := NewRetryClient(f.inner, retries, f.sleep).(*RetryClient)
	request, _ := http.NewRequest("POST", "/", strings.NewReader("request"))
	return client.Do(request)
}

/**************************************************************************/

func (f *RetryClientFixture) TestContextAlreadyCancelledReturnsImmediately() {
	f.inner = NewFailingHTTPClient(500, 500, 500, 500, 500)
	ctx, cancel := context.WithCancel(f.T().Context())
	cancel() // Cancel immediately

	client := NewRetryClient(f.inner, 10, f.sleep).(*RetryClient)
	request, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
	response, err := client.Do(request)

	f.So(response, should.BeNil)
	f.So(err, should.Equal, context.Canceled)
	f.So(f.inner.call, should.Equal, 0) // No requests should be made
}

func (f *RetryClientFixture) TestContextAlreadyCancelledReturnsImmediatelyForPost() {
	f.inner = NewFailingHTTPClient(500, 500, 500, 500, 500)
	ctx, cancel := context.WithCancel(f.T().Context())
	cancel() // Cancel immediately

	client := NewRetryClient(f.inner, 10, f.sleep).(*RetryClient)
	request, _ := http.NewRequestWithContext(ctx, "POST", "/", strings.NewReader("body"))
	response, err := client.Do(request)

	f.So(response, should.BeNil)
	f.So(err, should.Equal, context.Canceled)
	f.So(f.inner.call, should.Equal, 0) // No requests should be made
}

func (f *RetryClientFixture) TestContextCancelledDuringBackoffStopsRetryingGet() {
	f.inner = NewFailingHTTPClient(500, 500, 500)
	ctx, cancel := context.WithCancel(f.T().Context())

	sleepCount := 0
	cancellingSleeper := func(_ context.Context, d time.Duration) {
		sleepCount++
		if sleepCount == 2 {
			cancel()
		}
	}

	client := NewRetryClient(f.inner, 10, cancellingSleeper).(*RetryClient)
	request, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
	response, err := client.Do(request)

	f.So(response, should.BeNil)
	f.So(err, should.Equal, context.Canceled)
	f.So(f.inner.call, should.Equal, 2)
}

func (f *RetryClientFixture) TestContextCancelledDuringBackoffStopsRetryingPost() {
	f.inner = NewFailingHTTPClient(500, 500, 500)
	ctx, cancel := context.WithCancel(f.T().Context())

	sleepCount := 0
	cancellingSleeper := func(_ context.Context, d time.Duration) {
		sleepCount++
		if sleepCount == 2 {
			cancel()
		}
	}

	client := NewRetryClient(f.inner, 10, cancellingSleeper).(*RetryClient)
	request, _ := http.NewRequestWithContext(ctx, "POST", "/", strings.NewReader("body"))
	response, err := client.Do(request)

	f.So(response, should.BeNil)
	f.So(err, should.Equal, context.Canceled)
	f.So(f.inner.call, should.Equal, 2)
}

func (f *RetryClientFixture) TestContextCancelledDuringRequestStopsRetryingGet() {
	ctx, cancel := context.WithCancel(f.T().Context())

	cancellingClient := &ContextCancellingHTTPClient{
		cancelOnCall: 2,
		cancel:       cancel,
		inner:        NewFailingHTTPClient(500, 500, 500, 500, 500),
	}

	client := NewRetryClient(cancellingClient, 10, f.sleep).(*RetryClient)
	request, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
	response, err := client.Do(request)

	f.So(response, should.BeNil)
	f.So(err, should.Equal, context.Canceled)
	f.So(cancellingClient.inner.call, should.Equal, 2) // Should not retry after cancellation
}

func (f *RetryClientFixture) TestContextCancelledDuringRequestStopsRetryingPost() {
	ctx, cancel := context.WithCancel(f.T().Context())

	cancellingClient := &ContextCancellingHTTPClient{
		cancelOnCall: 2,
		cancel:       cancel,
		inner:        NewFailingHTTPClient(500, 500, 500, 500, 500),
	}

	client := NewRetryClient(cancellingClient, 10, f.sleep).(*RetryClient)
	request, _ := http.NewRequestWithContext(ctx, "POST", "/", strings.NewReader("body"))
	response, err := client.Do(request)

	f.So(response, should.BeNil)
	f.So(err, should.Equal, context.Canceled)
	f.So(cancellingClient.inner.call, should.Equal, 2) // Should not retry after cancellation
}

func (f *RetryClientFixture) TestContextCancelledDuringRateLimitBackoffStopsRetrying() {
	f.inner = NewFailingHTTPClient(429, 429, 429)
	ctx, cancel := context.WithCancel(f.T().Context())

	sleepCount := 0
	cancellingSleeper := func(_ context.Context, d time.Duration) {
		sleepCount++
		if sleepCount == 1 {
			cancel()
		}
	}

	// 429 handling: request made, then single sleep in backOff (count=1, cancels), then ctx.Err() exits.
	client := NewRetryClient(f.inner, 10, cancellingSleeper).(*RetryClient)
	request, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
	response, err := client.Do(request)

	f.So(response, should.BeNil)
	f.So(err, should.Equal, context.Canceled)
	f.So(f.inner.call, should.Equal, 1)
}

func (f *RetryClientFixture) TestRequestContextIsPassedToSleeper() {
	type contextKey string
	ctx := context.WithValue(f.T().Context(), contextKey("key"), "value")
	f.inner = NewFailingHTTPClient(500, 429, http.StatusOK)
	f.inner.responses[2].Body = io.NopCloser(strings.NewReader("Success"))

	var sleeperContexts []context.Context
	recordingSleeper := func(ctx context.Context, _ time.Duration) {
		sleeperContexts = append(sleeperContexts, ctx)
	}

	client := NewRetryClient(f.inner, 10, recordingSleeper).(*RetryClient)
	request, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
	_, err := client.Do(request)

	f.So(err, should.BeNil)
	f.So(f.inner.call, should.Equal, 3)
	// Both the randomized and the rate-limited backoff paths must hand the
	// request's own context to the sleeper so cancellation can cut a sleep short.
	if f.So(len(sleeperContexts), should.Equal, 2) {
		for _, sleeperContext := range sleeperContexts {
			f.So(sleeperContext.Value(contextKey("key")), should.Equal, "value")
		}
	}
}

/**************************************************************************/

func TestContextSleep(t *testing.T) {
	// Each case runs in a synctest bubble: the clock is virtual, so sleeps
	// complete instantly and elapsed durations are exact rather than approximate.
	bubble := func(name string, test func(t *testing.T)) {
		t.Run(name, func(t *testing.T) { synctest.Test(t, test) })
	}

	bubble("SleepsForFullDurationWhenContextNotCancelled", func(t *testing.T) {
		start := time.Now()

		ContextSleep(t.Context(), 50*time.Millisecond)

		assertions.New(t).So(time.Since(start), should.Equal, 50*time.Millisecond)
	})

	bubble("ReturnsImmediatelyWhenContextAlreadyCancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		start := time.Now()

		ContextSleep(ctx, time.Second)

		assertions.New(t).So(time.Since(start), should.BeZeroValue)
	})

	bubble("ReturnsEarlyWhenContextCancelledDuringSleep", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		start := time.Now()

		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		ContextSleep(ctx, time.Second)

		assertions.New(t).So(time.Since(start), should.Equal, 50*time.Millisecond)
	})
}
