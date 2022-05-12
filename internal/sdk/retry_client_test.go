package sdk

import (
	"errors"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
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
	client := NewRetryClient(f.inner, 10, rand.New(rand.NewSource(0)), f.sleep).(*RetryClient)
	request, _ := http.NewRequest("POST", "/", &ErrorProneReadCloser{readError: errors.New("GOPHERS!")})
	return client.Do(request)
}
func (f *RetryClientFixture) sleep(duration time.Duration) {
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
	f.So(f.naps, should.Resemble,
		[]time.Duration{0 * time.Second, 0 * time.Second, 1 * time.Second, 2 * time.Second})
}

/**************************************************************************/

func (f *RetryClientFixture) TestRetryOnBadResponseUntilSuccess() {
	f.inner = NewFailingHTTPClient(500, 501, 502, 522, 200)

	f.response, f.err = f.sendPostWithRetry(4)

	f.assertRequestWasSuccessful()
	f.assertBackOffStrategyWasObserved()
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
	f.So(f.naps[0], should.Equal, 0)
	for i := 1; i < len(f.naps); i++ {
		f.So(f.naps[i], should.BeBetweenOrEqual, 0, time.Second*time.Duration(min(i, maxBackOffDuration)))
	}
}

func (f *RetryClientFixture) TestBackOffRateLimitedGet() {
	retries := 5
	x := http.StatusTooManyRequests
	f.inner = NewFailingHTTPClient(x, x, x, x, x, x, x, x, x, x, http.StatusOK) //x10 rate limits > retries
	f.inner.responses[10].Body = io.NopCloser(strings.NewReader("Alohomora"))

	_, f.err = f.sendGetWithRetry(retries - 1)

	f.So(f.err, should.BeNil)
	if f.So(f.inner.call, should.Equal, 11) {
		var napTotal time.Duration
		for i := 0; i < 10; i++ {
			napTotal += f.naps[i]
			f.So(f.naps[i], should.BeBetweenOrEqual, 0, backOffRateLimit*time.Second)
		}
		f.So(napTotal, should.BeGreaterThan, time.Second*5)
	}
}

func (f *RetryClientFixture) TestBackOffRateLimitedPost() {
	retries := 5
	x := http.StatusTooManyRequests
	f.inner = NewFailingHTTPClient(x, x, x, x, x, x, x, x, x, x, http.StatusOK) //x10 rate limits > retries
	f.inner.responses[10].Body = io.NopCloser(strings.NewReader("Alohomora"))

	_, f.err = f.sendPostWithRetry(retries - 1)

	f.So(f.err, should.BeNil)
	if f.So(f.inner.call, should.Equal, 11) {
		var napTotal time.Duration
		for i := 0; i < 10; i++ {
			napTotal += f.naps[i]
			f.So(f.naps[i], should.BeBetweenOrEqual, 0, backOffRateLimit*time.Second)
		}
		f.So(napTotal, should.BeGreaterThan, time.Second*5)
	}
}

/**************************************************************************/

func (f *RetryClientFixture) sendGetWithRetry(retries int) (*http.Response, error) {
	client := NewRetryClient(f.inner, retries, rand.New(rand.NewSource(0)), f.sleep).(*RetryClient)
	request, _ := http.NewRequest("GET", "/?body=request", nil)
	return client.Do(request)
}
func (f *RetryClientFixture) sendPostWithRetry(retries int) (*http.Response, error) {
	client := NewRetryClient(f.inner, retries, rand.New(rand.NewSource(0)), f.sleep).(*RetryClient)
	request, _ := http.NewRequest("POST", "/", strings.NewReader("request"))
	return client.Do(request)
}
