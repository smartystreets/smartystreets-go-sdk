package us_street

import "net/http"

type Sender interface {
	Do(*http.Request) ([]byte, error)
}

type Credential interface {
	Sign(*http.Request) error
}
