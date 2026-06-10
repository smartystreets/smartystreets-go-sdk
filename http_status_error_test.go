package sdk

import (
	"net/http"
	"testing"

	"github.com/smarty/assertions"
	"github.com/smarty/assertions/should"
)

func TestHTTPStatusErrorUsesAPIErrorMessageWhenPresent(t *testing.T) {
	t.Parallel()

	err := NewHTTPStatusError(http.StatusUnauthorized, []byte(`{"errors":[{"message":"API says no"}]}`))

	assert := assertions.New(t)
	assert.So(err.Error(), should.Equal, "HTTP 401 Unauthorized\nAPI says no")
	assert.So(err.StatusCode(), should.Equal, http.StatusUnauthorized)
	assert.So(err.Content(), should.Equal, `{"errors":[{"message":"API says no"}]}`)
}

func TestHTTPStatusErrorJoinsMultipleAPIErrorMessages(t *testing.T) {
	t.Parallel()

	err := NewHTTPStatusError(http.StatusBadRequest, []byte(`{"errors":[{"message":"first"},{"message":"second"}]}`))

	assertions.New(t).So(err.Error(), should.Equal, "HTTP 400 Bad Request\nfirst second")
}

func TestHTTPStatusErrorFallsBackToStandardMessages(t *testing.T) {
	t.Parallel()

	cases := map[int]string{
		http.StatusNotModified:           "Not Modified: The requested record has not been modified since the previous request with the Etag value.",
		http.StatusBadRequest:            "Bad Request (Malformed Payload): A GET request lacked a required field or the request body of a POST request contained malformed JSON.",
		http.StatusUnauthorized:          "Unauthorized: The credentials were provided incorrectly or did not match any existing, active credentials.",
		http.StatusPaymentRequired:       "Payment Required: There is no active subscription for the account associated with the credentials submitted with the request.",
		http.StatusForbidden:             "Forbidden: The request contained valid data and was understood by the server, but the server is refusing action.",
		http.StatusRequestTimeout:        "Request timeout error.",
		http.StatusRequestEntityTooLarge: "Request Entity Too Large: The request body has exceeded the maximum size.",
		http.StatusUnprocessableEntity:   "GET request lacked required fields.",
		http.StatusTooManyRequests:       "Too Many Requests: The rate limit for your account has been exceeded.",
		http.StatusInternalServerError:   "Internal Server Error.",
		http.StatusBadGateway:            "Bad Gateway error.",
		http.StatusServiceUnavailable:    "Service Unavailable. Try again later.",
		http.StatusGatewayTimeout:        "The upstream data provider did not respond in a timely fashion and the request failed. A serious, yet rare occurrence indeed.",
		http.StatusTeapot:                "The server returned an unexpected HTTP status code: 418",
	}

	assert := assertions.New(t)
	for code, message := range cases {
		err := NewHTTPStatusError(code, nil)
		assert.So(err.Error(), should.Equal, statusText(code)+"\n"+message)
	}
}

func TestHTTPStatusErrorFallsBackWhenContentNotUsable(t *testing.T) {
	t.Parallel()

	fallback := "HTTP 401 Unauthorized\nUnauthorized: The credentials were provided incorrectly or did not match any existing, active credentials."

	assert := assertions.New(t)
	assert.So(NewHTTPStatusError(401, []byte("not json")).Error(), should.Equal, fallback)
	assert.So(NewHTTPStatusError(401, []byte(`{"other":"shape"}`)).Error(), should.Equal, fallback)
	assert.So(NewHTTPStatusError(401, []byte(`{"errors":[]}`)).Error(), should.Equal, fallback)
	assert.So(NewHTTPStatusError(401, []byte(`{"errors":[{"message":""}]}`)).Error(), should.Equal, fallback)
}

func TestHTTPStatusErrorContentStillExposesRawBody(t *testing.T) {
	t.Parallel()

	err := NewHTTPStatusError(http.StatusTeapot, []byte("Hello, World!"))

	assert := assertions.New(t)
	assert.So(err.StatusCode(), should.Equal, http.StatusTeapot)
	assert.So(err.Content(), should.Equal, "Hello, World!")
}

func TestNilHTTPStatusErrorBehavesLikeHTTP200(t *testing.T) {
	t.Parallel()

	var err *HTTPStatusError

	assert := assertions.New(t)
	assert.So(err, should.BeNil)
	assert.So(err.Error(), should.Equal, "HTTP 200 OK")
	assert.So(err.StatusCode(), should.Equal, http.StatusOK)
	assert.So(err.Content(), should.BeEmpty)
}
