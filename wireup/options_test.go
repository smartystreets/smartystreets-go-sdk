package wireup

import (
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestOptionsFixture(t *testing.T) {
	gunit.Run(new(OptionsFixture), t)
}

type OptionsFixture struct {
	*gunit.Fixture
}

func (this *OptionsFixture) TestConfigure_NilOptionIgnored() {
	this.So(func() { configure(nil) }, should.NotPanic)
}

func (this *OptionsFixture) TestConfigure_WithCustomCommaSeparatedQuery() {
	builder := configure(
		WithCustomCommaSeparatedQuery("test", "first"),
		WithCustomCommaSeparatedQuery("test", "second"),
		WithCustomCommaSeparatedQuery("test", "third"),
	)
	this.So(builder.customQueries.Get("test"), should.Equal, "first,second,third")
}

func (this *OptionsFixture) TestConfigure_WithFeatureComponentAnalysis() {
	builder := configure(
		WithFeatureComponentAnalysis(),
	)
	this.So(builder.customQueries.Get("features"), should.Equal, "component-analysis")
}

func (this *OptionsFixture) TestConfigure_WithFeatureComponentAnalysisAndCustom_ShouldAppend() {
	builder := configure(
		WithFeatureComponentAnalysis(),
		WithCustomCommaSeparatedQuery("features", "test-feature"),
	)
	this.So(builder.customQueries.Get("features"), should.Equal, "component-analysis,test-feature")
}

func (this *OptionsFixture) TestConfigure_WithCustomQuery() {
	builder := configure(
		WithCustomQuery("test", "first"),
	)
	this.So(builder.customQueries.Get("test"), should.Equal, "first")
}

func (this *OptionsFixture) TestConfigure_WithCustomQuery_ShouldOverwrite() {
	builder := configure(
		WithCustomQuery("test", "first"),
		WithCustomQuery("test", "second"),
		WithCustomQuery("test", "third"),
	)
	this.So(builder.customQueries.Get("test"), should.Equal, "third")
}
