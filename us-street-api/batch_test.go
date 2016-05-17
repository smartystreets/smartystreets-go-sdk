package us_street

import (
	"strconv"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type BatchFixture struct {
	*gunit.Fixture
}

func (f *BatchFixture) TestCapacityIsLimitedAt100Inputs() {
	batch := NewBatch()
	for x := 0; x < 100; x++ {
		f.So(batch.Append(&Input{InputID: strconv.Itoa(x)}), should.BeTrue)
	}
	f.So(batch.Length(), should.Equal, 100)

	for x := 100; x < 200; x++ {
		f.So(batch.Append(&Input{InputID: strconv.Itoa(x)}), should.BeFalse)
	}

	f.So(batch.Length(), should.Equal, 100)
}

func (f *BatchFixture) TestJSONSerializationShouldNeverFail() {
	batch := NewBatch()
	batch.Append(&Input{
		Street:        "This",
		Street2:       "test",
		Secondary:     "exists",
		City:          "to",
		State:         "ensure",
		ZIPCode:       "the",
		LastLine:      "input",
		Addressee:     "always",
		Urbanization:  "serializes",
		InputID:       "successfully",
		MaxCandidates: 7,
	})
	serialized, err := batch.marshalJSON()
	f.So(err, should.BeNil)
	f.So(serialized, should.NotBeEmpty)
}
