package sdk

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"

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
	header   http.Header

	naps []time.Duration
}

func (f *RetryClientFixture) TestRequestBodyCannotBeBuffered_ErrorReturnedImmediately() {
	f.response, f.err = f.sendErrorProneRequest()
	f.assertReadErrorReturnedAndRequestNotSent()
}
func (f *RetryClientFixture) sendErrorProneRequest() (*http.Response, error) {
	f.inner = &FakeMultiHTTPClient{}
	client := NewRetryClient(f.inner, 10, rand.New(rand.NewSource(0)), f.sleep).(*RetryClient)
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
	client := NewRetryClient(inner, 0, rand.New(rand.NewSource(0)), f.sleep)
	f.So(client, should.Equal, inner)
}

/**************************************************************************/

func (f *RetryClientFixture) TestBackOffNeverToExceedHardCodedMaximum() {
	retries := 2000
	f.inner = NewFailingHTTPClient(make([]int, retries)...)

	_, f.err = f.sendPostWithRetry(retries - 1)

	f.So(f.err, should.BeNil)
	f.So(f.inner.call, should.Equal, retries)
	for i := 0; i < len(f.naps); i++ {
		cap := time.Second * time.Duration(min(i+1, maxBackOffDuration))
		f.So(f.naps[i], should.BeLessThanOrEqualTo, cap)
	}
}

func (f *RetryClientFixture) TestBackOffRateLimitedGet() {
	retries := 5
	x := http.StatusTooManyRequests
	f.inner = NewFailingHTTPClient(x, x, x, x, x, x, x, x, x, x, http.StatusOK) //x10 rate limits > retries
	f.inner.responses[10].Body = io.NopCloser(strings.NewReader("Alohomora"))

	_, f.err = f.sendGetWithRetry(retries - 1)

	// 429s retry indefinitely (attempt reset to 1), so all 11 requests should be made
	f.So(f.err, should.BeNil)
	f.So(f.inner.call, should.Equal, 11)
	// Each 429 triggers a rate limit sleep plus a backoff sleep
	f.So(len(f.naps), should.BeGreaterThan, 0)
}

func (f *RetryClientFixture) TestBackOffRateLimitedPost() {
	retries := 5
	x := http.StatusTooManyRequests
	f.inner = NewFailingHTTPClient(x, x, x, x, x, x, x, x, x, x, http.StatusOK) //x10 rate limits > retries
	f.inner.responses[10].Body = io.NopCloser(strings.NewReader("Alohomora"))

	_, f.err = f.sendPostWithRetry(retries - 1)

	// 429s retry indefinitely (attempt reset to 1), so all 11 requests should be made
	f.So(f.err, should.BeNil)
	f.So(f.inner.call, should.Equal, 11)
	// Each 429 triggers a rate limit sleep plus a backoff sleep
	f.So(len(f.naps), should.BeGreaterThan, 0)
}

func (f *RetryClientFixture) TestRetryAfterHeaderUsedFor429() {
	f.inner = NewFailingHTTPClient(429, 429, 429, http.StatusOK)
	retryAfterSeconds := 7
	f.inner.headerKey = "Retry-After"
	f.inner.rateLimitTime = retryAfterSeconds
	f.inner.responses[3].Body = io.NopCloser(strings.NewReader("Success"))

	_, f.err = f.sendPostWithRetry(3) // 3 retries allows for 4 attempts

	f.So(f.err, should.BeNil)
	f.So(f.inner.call, should.Equal, 4)
	// At least one nap should be the exact Retry-After value (from 429 handling)
	hasRetryAfterNap := false
	for _, nap := range f.naps {
		if nap == time.Second*time.Duration(retryAfterSeconds) {
			hasRetryAfterNap = true
			break
		}
	}
	f.So(hasRetryAfterNap, should.BeTrue)
}

func (f *RetryClientFixture) TestFallsBackToRandomBackoffWithoutRetryAfterHeader() {
	f.inner = NewFailingHTTPClient(429, 429, 429, http.StatusOK)
	f.inner.headerKey = "X-Invalid-Header" // Not Retry-After
	f.inner.rateLimitTime = 100            // Would be obvious if used
	f.inner.responses[3].Body = io.NopCloser(strings.NewReader("Success"))

	_, f.err = f.sendPostWithRetry(3) // 3 retries allows for 4 attempts

	f.So(f.err, should.BeNil)
	f.So(f.inner.call, should.Equal, 4)
	// Without Retry-After header, should use random backoff (0 to backOffRateLimit)
	for _, nap := range f.naps {
		f.So(nap, should.BeLessThanOrEqualTo, time.Second*time.Duration(max(backOffRateLimit, maxBackOffDuration)))
	}
}

/**************************************************************************/

func (f *RetryClientFixture) sendGetWithRetry(retries int) (*http.Response, error) {
	if len(f.inner.responses) <= retries {
		f.T().Fatalf("The number of retries is greater than or equal to the number of status codes provided. Please ensure that the number of retries is less than the number of status codes provided.")
	}

	client := NewRetryClient(f.inner, retries, rand.New(rand.NewSource(0)), f.sleep).(*RetryClient)
	request, _ := http.NewRequest("GET", "/?body=request", nil)
	return client.Do(request)
}
func (f *RetryClientFixture) sendPostWithRetry(retries int) (*http.Response, error) {
	if len(f.inner.responses) <= retries {
		f.T().Fatalf("The number of retries is greater than or equal to the number of status codes provided. Please ensure that the number of retries is less than the number of status codes provided.")
	}

	client := NewRetryClient(f.inner, retries, rand.New(rand.NewSource(0)), f.sleep).(*RetryClient)
	request, _ := http.NewRequest("POST", "/", strings.NewReader("request"))
	return client.Do(request)
}

/**************************************************************************/

func (f *RetryClientFixture) TestContextAlreadyCancelledReturnsImmediately() {
	f.inner = NewFailingHTTPClient(500, 500, 500, 500, 500)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	client := NewRetryClient(f.inner, 10, rand.New(rand.NewSource(0)), f.sleep).(*RetryClient)
	request, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
	response, err := client.Do(request)

	f.So(response, should.BeNil)
	f.So(err, should.Equal, context.Canceled)
	f.So(f.inner.call, should.Equal, 0) // No requests should be made
}

func (f *RetryClientFixture) TestContextAlreadyCancelledReturnsImmediatelyForPost() {
	f.inner = NewFailingHTTPClient(500, 500, 500, 500, 500)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	client := NewRetryClient(f.inner, 10, rand.New(rand.NewSource(0)), f.sleep).(*RetryClient)
	request, _ := http.NewRequestWithContext(ctx, "POST", "/", strings.NewReader("body"))
	response, err := client.Do(request)

	f.So(response, should.BeNil)
	f.So(err, should.Equal, context.Canceled)
	f.So(f.inner.call, should.Equal, 0) // No requests should be made
}

func (f *RetryClientFixture) TestContextCancelledDuringBackoffStopsRetryingGet() {
	f.inner = NewFailingHTTPClient(500, 500, 500)
	ctx, cancel := context.WithCancel(context.Background())

	sleepCount := 0
	cancellingSleeper := func(_ context.Context, d time.Duration) {
		sleepCount++
		if sleepCount == 2 {
			cancel()
		}
	}

	client := NewRetryClient(f.inner, 10, rand.New(rand.NewSource(0)), cancellingSleeper).(*RetryClient)
	request, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
	response, err := client.Do(request)

	f.So(response, should.BeNil)
	f.So(err, should.Equal, context.Canceled)
	f.So(f.inner.call, should.Equal, 2)
}

func (f *RetryClientFixture) TestContextCancelledDuringBackoffStopsRetryingPost() {
	f.inner = NewFailingHTTPClient(500, 500, 500)
	ctx, cancel := context.WithCancel(context.Background())

	sleepCount := 0
	cancellingSleeper := func(_ context.Context, d time.Duration) {
		sleepCount++
		if sleepCount == 2 {
			cancel()
		}
	}

	client := NewRetryClient(f.inner, 10, rand.New(rand.NewSource(0)), cancellingSleeper).(*RetryClient)
	request, _ := http.NewRequestWithContext(ctx, "POST", "/", strings.NewReader("body"))
	response, err := client.Do(request)

	f.So(response, should.BeNil)
	f.So(err, should.Equal, context.Canceled)
	f.So(f.inner.call, should.Equal, 2)
}

func (f *RetryClientFixture) TestContextCancelledDuringRequestStopsRetryingGet() {
	ctx, cancel := context.WithCancel(context.Background())

	cancellingClient := &ContextCancellingHTTPClient{
		cancelOnCall: 2,
		cancel:       cancel,
		inner:        NewFailingHTTPClient(500, 500, 500, 500, 500),
	}

	client := NewRetryClient(cancellingClient, 10, rand.New(rand.NewSource(0)), f.sleep).(*RetryClient)
	request, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
	response, err := client.Do(request)

	f.So(response, should.BeNil)
	f.So(err, should.Equal, context.Canceled)
	f.So(cancellingClient.inner.call, should.Equal, 2) // Should not retry after cancellation
}

func (f *RetryClientFixture) TestContextCancelledDuringRequestStopsRetryingPost() {
	ctx, cancel := context.WithCancel(context.Background())

	cancellingClient := &ContextCancellingHTTPClient{
		cancelOnCall: 2,
		cancel:       cancel,
		inner:        NewFailingHTTPClient(500, 500, 500, 500, 500),
	}

	client := NewRetryClient(cancellingClient, 10, rand.New(rand.NewSource(0)), f.sleep).(*RetryClient)
	request, _ := http.NewRequestWithContext(ctx, "POST", "/", strings.NewReader("body"))
	response, err := client.Do(request)

	f.So(response, should.BeNil)
	f.So(err, should.Equal, context.Canceled)
	f.So(cancellingClient.inner.call, should.Equal, 2) // Should not retry after cancellation
}

func (f *RetryClientFixture) TestContextCancelledDuringRateLimitBackoffStopsRetrying() {
	f.inner = NewFailingHTTPClient(429, 429, 429)
	ctx, cancel := context.WithCancel(context.Background())

	sleepCount := 0
	cancellingSleeper := func(_ context.Context, d time.Duration) {
		sleepCount++
		if sleepCount == 2 {
			cancel()
		}
	}

	// 429 handling: request made, then sleep in handleHttpStatusCode (count=1),
	// then backOff sleep (count=2, cancels), then ctx.Err() check exits.
	client := NewRetryClient(f.inner, 10, rand.New(rand.NewSource(0)), cancellingSleeper).(*RetryClient)
	request, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
	response, err := client.Do(request)

	f.So(response, should.BeNil)
	f.So(err, should.Equal, context.Canceled)
	f.So(f.inner.call, should.Equal, 1)
}

/**************************************************************************/

func TestContextSleep(t *testing.T) {
	gunit.Run(new(ContextSleepFixture), t)
}

type ContextSleepFixture struct {
	*gunit.Fixture
}

func (f *ContextSleepFixture) TestSleepsForFullDurationWhenContextNotCancelled() {
	ctx := context.Background()
	start := time.Now()

	ContextSleep(ctx, 50*time.Millisecond)

	elapsed := time.Since(start)
	f.So(elapsed, should.BeGreaterThanOrEqualTo, 50*time.Millisecond)
}

func (f *ContextSleepFixture) TestReturnsImmediatelyWhenContextAlreadyCancelled() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := time.Now()

	ContextSleep(ctx, 1*time.Second)

	elapsed := time.Since(start)
	f.So(elapsed, should.BeLessThan, 100*time.Millisecond)
}

func (f *ContextSleepFixture) TestReturnsEarlyWhenContextCancelledDuringSleep() {
	ctx, cancel := context.WithCancel(context.Background())
	start := time.Now()

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	ContextSleep(ctx, 1*time.Second)

	elapsed := time.Since(start)
	f.So(elapsed, should.BeGreaterThanOrEqualTo, 50*time.Millisecond)
	f.So(elapsed, should.BeLessThan, 200*time.Millisecond)
}
