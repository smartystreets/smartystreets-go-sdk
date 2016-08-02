package smarty_sdk

import (
	"testing"
	"net/http"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func TestHTTPStatusError(t *testing.T) {
	t.Parallel()

	err := NewHTTPStatusError(http.StatusTeapot, []byte("Hello, World!"))

	assert := assertions.New(t)
	assert.So(err.Error(), should.Equal, http.StatusText(http.StatusTeapot))
	assert.So(err.StatusCode(), should.Equal, http.StatusTeapot)
	assert.So(err.Content(), should.Resemble, []byte("Hello, World!"))
}

func TestNilHTTPStatusErrorBehavesLikeHTTP200(t *testing.T) {
	t.Parallel()

	var err *HTTPStatusError

	assert := assertions.New(t)
	assert.So(err, should.BeNil)
	assert.So(err.Error(), should.Equal, http.StatusText(http.StatusOK))
	assert.So(err.StatusCode(), should.Equal, http.StatusOK)
	assert.So(err.Content(), should.BeNil)
}
