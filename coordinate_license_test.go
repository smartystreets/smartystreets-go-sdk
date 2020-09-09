package sdk

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestCoordinateLicenseFixture(t *testing.T) {
    gunit.Run(new(CoordinateLicenseFixture), t)
}

type CoordinateLicenseFixture struct {
    *gunit.Fixture
}

func (this *CoordinateLicenseFixture) TestLicenseString() {
	this.So(CoordinateLicenseSmartyStreets.String(), should.Equal, "SmartyStreets")
	this.So(CoordinateLicenseGatewaySpatial.String(), should.Equal, "Gateway Spatial, LLC")
	this.So(CoordinateLicense(42).String(), should.Equal, "SmartyStreets")
}
