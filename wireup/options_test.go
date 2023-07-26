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
