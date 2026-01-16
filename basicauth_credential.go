package sdk

import (
	"net/http"
)

type basicAuthCredential struct {
	authID    string
	authToken string
}

func NewBasicAuthCredential(authID, authToken string) *basicAuthCredential {
	return &basicAuthCredential{
		authID:    authID,
		authToken: authToken,
	}
}

func (c basicAuthCredential) Sign(request *http.Request) error {
	SignRequest(request, c.authID, c.authToken)
	return nil
}

func SignRequest(request *http.Request, authID string, authToken string) {
	request.SetBasicAuth(authID, authToken)
}
