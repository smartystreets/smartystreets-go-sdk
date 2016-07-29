package smarty_sdk

import (
	"errors"
	"fmt"
	"net/http"
)

type Credential interface {
	Sign(*http.Request) error
}

// Common HTTP status codes returned by the SmartyStreets APIs:
var (
	StatusUnauthorized          = errors.New("401 Unauthorized")
	StatusPaymentRequired       = errors.New("402 Payment Required")
	StatusBadRequest            = errors.New("400 Bad Request")
	StatusRequestEntityTooLarge = errors.New("413 Request entity too large")
	StatusTooManyRequests       = errors.New("429 Too many requests")
)

func StatusOtherError(status string, content []byte) error {
	return fmt.Errorf("Non-200 status: %s\n%s", status, string(content))
}
