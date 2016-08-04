package sdk

import "net/http"

// NopCredential - See https://en.wikipedia.org/wiki/Null_Object_pattern
type NopCredential struct{}

func (c NopCredential) Sign(request *http.Request) error {
	return nil
}
