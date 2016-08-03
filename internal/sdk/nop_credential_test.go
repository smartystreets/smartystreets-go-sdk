package sdk

import (
	"testing"
	"net/http"
	"strings"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/assertions"
)

func TestNopCredentialDoesNothing(t *testing.T) {
	t.Parallel()
	assert := assertions.New(t)

	credential := &NopCredential{}
	request, _ := http.NewRequest("GET", "https://www.google.com", strings.NewReader("Hello, World!"))

	before := dumpRequest(request)
	err := credential.Sign(request)
	after := dumpRequest(request)

	assert.So(err, should.BeNil)
	assert.So(before, should.Equal, after)
}
