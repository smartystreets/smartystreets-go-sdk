package sdk

import (
	"net/http"
)

type headerCredential struct {
	authID    string
	authToken string
}

func NewHeaderCredential(authID, authToken string) *headerCredential {
	return &headerCredential{
		authID:    authID,
		authToken: authToken,
	}
}

func (c headerCredential) Sign(request *http.Request) error {
	SignRequest(request, c.authID, c.authToken)
	return nil
}

func SignRequest(request *http.Request, authID string, authToken string) {
	request.SetBasicAuth(authID, authToken)
}
