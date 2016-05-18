package sdk

import "net/http"

type SecretKeyCredential struct {
	AuthID    string
	AuthToken string
}

func (c SecretKeyCredential) Sign(request *http.Request) error {
	query := request.URL.Query()
	query.Set("auth-id", c.AuthID)
	query.Set("auth-token", c.AuthToken)
	request.URL.RawQuery = query.Encode()
	return nil
}
