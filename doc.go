// sdk is a top-level package containing elements common to all SmartyStreets APIs.
package sdk

import "net/http"

type RequestSender interface { // TODO: make this public, top-level
	Send(*http.Request) ([]byte, error)
}

