package wireup

import (
	"net/http"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestBuilderFixture(t *testing.T) {
	gunit.Run(new(BuilderFixture), t)
}

type BuilderFixture struct {
	*gunit.Fixture
}

func (this *BuilderFixture) TestTransport_HTTP2EnabledByDefault() {
	transport := configure().buildTransport()

	this.So(transport.ForceAttemptHTTP2, should.BeTrue)
	this.So(transport.TLSNextProto, should.BeNil)
	this.So(transport.DialContext, should.NotBeNil)
}

func (this *BuilderFixture) TestTransport_DisableHTTP2OptsOut() {
	transport := configure(DisableHTTP2()).buildTransport()

	this.So(transport.ForceAttemptHTTP2, should.BeFalse)
	this.So(transport.TLSNextProto, should.NotBeNil)
	this.So(transport.TLSNextProto, should.BeEmpty)
	this.So(transport.DialContext, should.NotBeNil)
}

func (this *BuilderFixture) TestTransport_ProxyAndIdleConnectionsWired() {
	transport := configure(
		ViaProxy("http://proxy.example:8080"),
		WithMaxIdleConnections(7),
	).buildTransport()

	this.So(transport.Proxy, should.NotBeNil)
	this.So(transport.MaxIdleConnsPerHost, should.Equal, 7)
}

func (this *BuilderFixture) TestClient_SuppliedClientUsedVerbatim() {
	supplied := &http.Client{}

	builder := configure(WithHTTPClient(supplied))

	this.So(builder.buildClient(), should.Equal, supplied)
}
