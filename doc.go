package smarty_sdk

import "net/http"

type Credential interface {
	Sign(*http.Request) error
}

func NewHTTPStatusError(statusCode int, content []byte) HTTPStatusError {
	return HTTPStatusError{
		statusCode: statusCode,
		content:    content,
	}
}

// HTTPStatusError stands for for the error type but also provides convenience methods
// for accessing the status code and content of the request that caused the error.
// Instances of this type are returned by sdk.HTTPSender.Send().
type HTTPStatusError struct {
	statusCode int
	content    []byte
}

func (e HTTPStatusError) Error() string {
	return http.StatusText(e.statusCode)
}

func (e HTTPStatusError) StatusCode() int {
	return e.statusCode
}

func (e HTTPStatusError) Content() []byte {
	return e.content
}
