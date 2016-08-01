package smarty_sdk

import "net/http"

type Credential interface {
	Sign(*http.Request) error
}

