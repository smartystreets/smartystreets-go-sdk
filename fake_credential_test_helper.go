package sdk

import "net/http"

// FakeCredential is a test helper that implements Credential.
// It returns the configured error from Sign without modifying the request.
type FakeCredential struct {
	Err error
}

func (f *FakeCredential) Sign(*http.Request) error {
	return f.Err
}
