package us_street

import "net/http"

type Sender interface {
	Do(*http.Request) ([]byte, error)
}
