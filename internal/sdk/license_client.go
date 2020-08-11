package sdk

import (
	"net/http"
	"strings"
)

type LicenseClient struct {
	inner    HTTPClient
	licenses []string
}

func NewLicenseClient(inner HTTPClient, licenses ...string) *LicenseClient {
	return &LicenseClient{inner: inner, licenses: licenses}
}

func (this *LicenseClient) Do(request *http.Request) (*http.Response, error) {
	if this.licenses[0] != "" {
		values := request.URL.Query()
		values.Set("license", strings.Join(this.licenses, ","))
		request.URL.RawQuery = values.Encode()
	}

	return this.inner.Do(request)
}
