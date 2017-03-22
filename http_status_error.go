package sdk

import (
	"fmt"
	"net/http"
)

func NewHTTPStatusError(statusCode int, content []byte) *HTTPStatusError {
	return &HTTPStatusError{
		statusCode: statusCode,
		content:    string(content),
	}
}

// HTTPStatusError stands in for the error type but also provides convenience methods
// for accessing the status code and content of the request that caused the error.
// Instances of this type are returned by sdk.HTTPSender.Send(). When nil, the methods
// of this type behave as if called on a non-nil instance instantiated with http.StatusOK (200).
type HTTPStatusError struct {
	statusCode int
	content    string
}

func (e *HTTPStatusError) Error() string {
	if e == nil {
		return statusText(http.StatusOK)
	}
	return statusText(e.statusCode) + "\n" + e.content
}

func statusText(code int) string {
	return fmt.Sprintf("HTTP %d %s", code, http.StatusText(code))
}

func (e *HTTPStatusError) StatusCode() int {
	if e == nil {
		return http.StatusOK
	}
	return e.statusCode
}

func (e *HTTPStatusError) Content() string {
	return e.content
}
