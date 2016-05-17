package us_street

import "net/http"

type requestSender interface {
	Send(*http.Request) ([]byte, error)
}
