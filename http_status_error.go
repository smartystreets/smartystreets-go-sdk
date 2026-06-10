package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	if message := extractAPIErrorMessage(e.content); message != "" {
		return statusText(e.statusCode) + "\n" + message
	}
	return statusText(e.statusCode) + "\n" + fallbackMessage(e.statusCode)
}

func statusText(code int) string {
	return fmt.Sprintf("HTTP %d %s", code, http.StatusText(code))
}

// extractAPIErrorMessage pulls the message(s) from the API's JSON error body
// ({"errors":[{"message":"..."}]}), returning "" when the body is empty,
// unparseable, or missing the expected fields.
func extractAPIErrorMessage(content string) string {
	var body struct {
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if json.Unmarshal([]byte(content), &body) != nil {
		return ""
	}
	var messages []string
	for _, item := range body.Errors {
		if message := strings.TrimSpace(item.Message); message != "" {
			messages = append(messages, message)
		}
	}
	return strings.Join(messages, " ")
}

func fallbackMessage(code int) string {
	switch code {
	case http.StatusNotModified:
		return "Not Modified: The requested record has not been modified since the previous request with the Etag value."
	case http.StatusBadRequest:
		return "Bad Request (Malformed Payload): A GET request lacked a required field or the request body of a POST request contained malformed JSON."
	case http.StatusUnauthorized:
		return "Unauthorized: The credentials were provided incorrectly or did not match any existing, active credentials."
	case http.StatusPaymentRequired:
		return "Payment Required: There is no active subscription for the account associated with the credentials submitted with the request."
	case http.StatusForbidden:
		return "Forbidden: The request contained valid data and was understood by the server, but the server is refusing action."
	case http.StatusRequestTimeout:
		return "Request timeout error."
	case http.StatusRequestEntityTooLarge:
		return "Request Entity Too Large: The request body has exceeded the maximum size."
	case http.StatusUnprocessableEntity:
		return "GET request lacked required fields."
	case http.StatusTooManyRequests:
		return "Too Many Requests: The rate limit for your account has been exceeded."
	case http.StatusInternalServerError:
		return "Internal Server Error."
	case http.StatusBadGateway:
		return "Bad Gateway error."
	case http.StatusServiceUnavailable:
		return "Service Unavailable. Try again later."
	case http.StatusGatewayTimeout:
		return "The upstream data provider did not respond in a timely fashion and the request failed. A serious, yet rare occurrence indeed."
	default:
		return fmt.Sprintf("The server returned an unexpected HTTP status code: %d", code)
	}
}

func (e *HTTPStatusError) StatusCode() int {
	if e == nil {
		return http.StatusOK
	}
	return e.statusCode
}

func (e *HTTPStatusError) Content() string {
	if e == nil {
		return ""
	}
	return e.content
}
