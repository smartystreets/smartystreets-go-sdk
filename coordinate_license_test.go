package sdk

import (
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestCoordinateLicenseFixture(t *testing.T) {
	gunit.Run(new(CoordinateLicenseFixture), t)
}

type CoordinateLicenseFixture struct {
	*gunit.Fixture
}

func (this *CoordinateLicenseFixture) TestLicenseString() {
	this.So(CoordinateLicenseSmartyStreets.String(), should.Equal, "SmartyStreets")
	this.So(CoordinateLicenseSmartyStreetsProprietary.String(), should.Equal, "SmartyStreets Proprietary")
	this.So(CoordinateLicense(42).String(), should.Equal, "SmartyStreets")
}
