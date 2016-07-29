package sdk

import (
	"errors"
	"net/http"

	"github.com/smartystreets/gunit"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/clock"
)

type RetryClientFixture struct {
	*gunit.Fixture
	sleeper *clock.Sleeper
}

func (this *RetryClientFixture) Setup() {
	this.sleeper = clock.StayAwake()
}

func (f *RetryClientFixture) TestRetryOnClientErrorUntilSuccess() {
	clientError := errors.New("Simulating Network Outage")
	inner := NewErringHTTPClient(clientError, clientError, clientError, clientError, nil)
	response, err := f.sendWithRetry(4, inner)

	if f.So(response, should.NotBeNil) {
		f.So(response.StatusCode, should.Equal, 200)
	}
	f.So(err, should.BeNil)
	f.So(inner.call, should.Equal, 5)
	f.So(f.sleeper.Naps, should.HaveLength, 4)
}

func (f *RetryClientFixture) TestRetryOnBadResponseUntilSuccess() {
	inner := NewFailingHTTPClient(400, 401, 402, 422, 200)
	response, err := f.sendWithRetry(4, inner)

	if f.So(response, should.NotBeNil) {
		f.So(response.StatusCode, should.Equal, 200)
	}
	f.So(err, should.BeNil)
	f.So(inner.call, should.Equal, 5)
	f.So(f.sleeper.Naps, should.HaveLength, 4)
}

func (f *RetryClientFixture) TestFailureReturnedIfRetryExceeded() {
	inner := NewFailingHTTPClient(500, 500, 500, 500, 500)
	response, err := f.sendWithRetry(4, inner)

	if f.So(response, should.NotBeNil) {
		f.So(response.StatusCode, should.Equal, 500)
	}
	f.So(err, should.BeNil)
	f.So(inner.call, should.Equal, 5)
	f.So(f.sleeper.Naps, should.HaveLength, 4)
}

func (f *RetryClientFixture) sendWithRetry(retries int, inner *FakeMultiHTTPClient) (*http.Response, error) {
	client := NewRetryClient(inner, retries)
	client.sleeper = f.sleeper
	request, _ := http.NewRequest("GET", "/", nil)
	return client.Do(request)
}