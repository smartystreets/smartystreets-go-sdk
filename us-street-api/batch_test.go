package us_street

import (
	"strconv"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type BatchFixture struct {
	*gunit.Fixture
}

func (this *BatchFixture) TestCapacityIsLimitedAt100Inputs() {
	batch := NewBatch()
	for x := 0; x < 100; x++ {
		this.So(batch.Append(&Input{InputID: strconv.Itoa(x)}), should.BeTrue)
	}
	this.So(batch.Length(), should.Equal, 100)

	for x := 100; x < 200; x++ {
		this.So(batch.Append(&Input{InputID: strconv.Itoa(x)}), should.BeFalse)
	}

	this.So(batch.Length(), should.Equal, 100)
}
