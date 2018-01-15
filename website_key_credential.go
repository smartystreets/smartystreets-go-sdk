package sdk

import (
	"net/http"
	"strings"
)

func NewWebsiteKeyCredential(key, hostNameOrIP string) *websiteKeyCredential {
	if !strings.HasPrefix(hostNameOrIP, httpsScheme) && !strings.HasPrefix(hostNameOrIP, httpScheme) {
		hostNameOrIP = httpScheme + hostNameOrIP
	}
	return &websiteKeyCredential{
		key:  key,
		host: hostNameOrIP,
	}
}

type websiteKeyCredential struct {
	key  string
	host string
}

func (c websiteKeyCredential) Sign(request *http.Request) error {
	query := request.URL.Query()
	query.Set("auth-id", c.key)
	request.URL.RawQuery = query.Encode()
	request.Header.Set("Referer", c.host)
	return nil
}

const (
	schemeTerminator = "://"
	httpScheme       = "http" + schemeTerminator
	httpsScheme      = "https" + schemeTerminator
)
