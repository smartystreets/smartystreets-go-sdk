package wireup

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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

func (this *BuilderFixture) TestClient_CleartextHTTP2UsesH2CTransport() {
	client := configure(EnableCleartextHTTP2()).buildClient()

	transport, ok := client.Transport.(*http2.Transport)
	this.So(ok, should.BeTrue)
	this.So(transport.AllowHTTP, should.BeTrue)
}

func (this *BuilderFixture) TestClient_CleartextHTTP2TakesPrecedenceOverDisableHTTP2() {
	client := configure(DisableHTTP2(), EnableCleartextHTTP2()).buildClient()

	_, ok := client.Transport.(*http2.Transport)
	this.So(ok, should.BeTrue)
}

func (this *BuilderFixture) TestClient_CleartextHTTP2PanicsOnNonHTTPBaseURL() {
	build := func() {
		configure(EnableCleartextHTTP2(), CustomBaseURL("https://example.com")).buildClient()
	}

	this.So(build, should.Panic)
}

func (this *BuilderFixture) TestClient_CleartextHTTP2NegotiatesH2() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(h2c.NewHandler(handler, &http2.Server{}))
	defer server.Close()

	client := configure(EnableCleartextHTTP2()).buildClient()

	response, err := client.Get(server.URL)
	this.So(err, should.BeNil)
	this.So(response, should.NotBeNil)
	if response == nil {
		return
	}
	defer func() { _ = response.Body.Close() }()
	this.So(response.ProtoMajor, should.Equal, 2)
}
