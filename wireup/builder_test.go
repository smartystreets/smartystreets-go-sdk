package wireup

import (
	"net/http"
	"net/http/httptest"
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
	if this.So(transport.Protocols, should.NotBeNil) {
		this.So(transport.Protocols.HTTP1(), should.BeTrue)
		this.So(transport.Protocols.HTTP2(), should.BeFalse)
		this.So(transport.Protocols.UnencryptedHTTP2(), should.BeFalse)
	}
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

func (this *BuilderFixture) TestTransport_CleartextHTTP2IgnoresProxy() {
	transport := configure(
		EnableCleartextHTTP2(),
		ViaProxy("http://proxy.example:8080"),
		WithMaxIdleConnections(7),
	).buildTransport()

	this.So(transport.Proxy, should.BeNil)
	this.So(transport.MaxIdleConnsPerHost, should.Equal, 7)
}

func (this *BuilderFixture) TestClient_SuppliedClientUsedVerbatim() {
	supplied := &http.Client{}

	builder := configure(WithHTTPClient(supplied))

	this.So(builder.buildClient(), should.Equal, supplied)
}

func (this *BuilderFixture) TestClient_CleartextHTTP2UsesUnencryptedHTTP2Only() {
	client := configure(EnableCleartextHTTP2()).buildClient()

	transport, ok := client.Transport.(*http.Transport)
	if !this.So(ok, should.BeTrue) {
		return
	}
	if !this.So(transport.Protocols, should.NotBeNil) {
		return
	}
	this.So(transport.Protocols.UnencryptedHTTP2(), should.BeTrue)
	this.So(transport.Protocols.HTTP1(), should.BeFalse)
	this.So(transport.Protocols.HTTP2(), should.BeFalse)
}

func (this *BuilderFixture) TestClient_CleartextHTTP2TakesPrecedenceOverDisableHTTP2() {
	client := configure(DisableHTTP2(), EnableCleartextHTTP2()).buildClient()

	transport, ok := client.Transport.(*http.Transport)
	if !this.So(ok, should.BeTrue) {
		return
	}
	if this.So(transport.Protocols, should.NotBeNil) {
		this.So(transport.Protocols.UnencryptedHTTP2(), should.BeTrue)
		this.So(transport.Protocols.HTTP1(), should.BeFalse)
	}
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
	server := httptest.NewUnstartedServer(handler)
	server.Config.Protocols = new(http.Protocols)
	server.Config.Protocols.SetHTTP1(true)
	server.Config.Protocols.SetUnencryptedHTTP2(true)
	server.Start()
	defer server.Close()

	client := configure(EnableCleartextHTTP2()).buildClient()

	response, err := client.Get(server.URL)
	this.So(err, should.BeNil)
	if !this.So(response, should.NotBeNil) {
		return
	}
	defer func() { _ = response.Body.Close() }()
	this.So(response.ProtoMajor, should.Equal, 2)
}
