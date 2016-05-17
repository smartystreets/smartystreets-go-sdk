package us_street

import "net/http"

type Sender interface {
	Send(*http.Request) ([]byte, error)
}
