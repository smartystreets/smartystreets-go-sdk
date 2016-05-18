package sdk

import "net/http"

type NopCredential struct {}

func (c NopCredential) Sign(request *http.Request) error {
	return nil
}
