package us_street

import (
	"testing"
	"net/http"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func TestSecretKeySigning(t *testing.T) {
	t.Parallel()

	credential := SecretKeyCredential{
		AuthID: "my id",
		AuthToken: "my token",
	}
	request, _ := http.NewRequest("GET", "http://google.com", nil)
	err := credential.Sign(request)
	assertions.New(t).So(err, should.BeNil)
	assertions.New(t).So(request.URL.String(), should.Equal, "http://google.com?auth-id=my+id&auth-token=my+token")
}
