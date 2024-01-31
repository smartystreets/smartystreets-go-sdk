// sdk is a top-level package containing elements common to all SmartyStreets APIs.
package sdk

import "net/http"

type RequestSender interface {
	Send(*http.Request) ([]byte, error)
	SendAndReturnHeaders(*http.Request) ([]byte, http.Header, error)
}
