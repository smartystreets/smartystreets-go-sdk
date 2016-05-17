package sdk

import "net/http"

type SecretKeyCredential struct {
	AuthID    string
	AuthToken string
}

func (this SecretKeyCredential) Sign(request *http.Request) error {
	query := request.URL.Query()
	query.Set("auth-id", this.AuthID)
	query.Set("auth-token", this.AuthToken)
	request.URL.RawQuery = query.Encode()
	return nil
}
