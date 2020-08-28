package sdk

import (
	"errors"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/clock"
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

	sleeper *clock.Sleeper
}

func (f *RetryClientFixture) Setup() {
	f.sleeper = clock.StayAwake()
}

func (f *RetryClientFixture) TestRequestBodyCannotBeBuffered_ErrorReturnedImmediately() {
	f.response, f.err = f.sendErrorProneRequest()
	f.assertReadErrorReturnedAndRequestNotSent()
}
func (f *RetryClientFixture) sendErrorProneRequest() (*http.Response, error) {
	f.inner = &FakeMultiHTTPClient{}
	client := NewRetryClient(f.inner, 10, rand.New(rand.NewSource(0))).(*RetryClient)
	client.sleeper = f.sleeper
	request, _ := http.NewRequest("POST", "/", &ErrorProneReadCloser{readError: errors.New("GOPHERS!")})
	return client.Do(request)
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
	f.So(f.sleeper.Naps, should.Resemble,
		[]time.Duration{2 * time.Second, 2 * time.Second, 3 * time.Second, 6 * time.Second})
}

/**************************************************************************/

func (f *RetryClientFixture) TestRetryOnBadResponseUntilSuccess() {
	f.inner = NewFailingHTTPClient(400, 401, 402, 422, 200)

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
	client := NewRetryClient(inner, 0, rand.New(rand.NewSource(0)))
	f.So(client, should.Equal, inner)
}

/**************************************************************************/

func (f *RetryClientFixture) TestBackOffNeverToExceedHardCodedMaximum() {
	f.inner = NewFailingHTTPClient(make([]int, 20)...)

	_, f.err = f.sendPostWithRetry(19)

	f.So(f.err, should.BeNil)
	f.So(f.inner.call, should.Equal, 20)
	f.So(f.sleeper.Naps, should.Resemble,

		[]time.Duration{
			time.Second * 2, // randomly between 0-2
			time.Second * 2, // randomly between 0-4
			time.Second * 3, // randomly between 0-8
			// the rest are randomly between 0-10 (capped)
			6 * time.Second, 5 * time.Second, 6 * time.Second, 7 * time.Second,
			7 * time.Second, 8 * time.Second, 8 * time.Second, 8 * time.Second,
			7 * time.Second, 9 * time.Second, 8 * time.Second, 2 * time.Second,
			6 * time.Second, 1 * time.Second, 0 * time.Second, 0 * time.Second})
}

/**************************************************************************/

func (f *RetryClientFixture) sendGetWithRetry(retries int) (*http.Response, error) {
	client := NewRetryClient(f.inner, retries, rand.New(rand.NewSource(0))).(*RetryClient)
	client.sleeper = f.sleeper
	request, _ := http.NewRequest("GET", "/?body=request", nil)
	return client.Do(request)
}
func (f *RetryClientFixture) sendPostWithRetry(retries int) (*http.Response, error) {
	client := NewRetryClient(f.inner, retries, rand.New(rand.NewSource(0))).(*RetryClient)
	client.sleeper = f.sleeper
	request, _ := http.NewRequest("POST", "/", strings.NewReader("request"))
	return client.Do(request)
}
