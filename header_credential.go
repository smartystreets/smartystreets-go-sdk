package sdk

import (
	"encoding/base64"
	"fmt"
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
	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte(authID+":"+authToken))))
}
