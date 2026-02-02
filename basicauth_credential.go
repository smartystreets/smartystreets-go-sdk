package sdk

import (
	"errors"
	"net/http"
)

var ErrCredentialsRequired = errors.New("credentials (auth id, auth token) required")

type basicAuthCredential struct {
	authID    string
	authToken string
}

func NewBasicAuthCredential(authID, authToken string) *basicAuthCredential {
	if len(authID) == 0 || len(authToken) == 0 {
		panic(ErrCredentialsRequired)
	}
	return &basicAuthCredential{
		authID:    authID,
		authToken: authToken,
	}
}

func (c basicAuthCredential) Sign(request *http.Request) error {
	request.SetBasicAuth(c.authID, c.authToken)
	return nil
}
