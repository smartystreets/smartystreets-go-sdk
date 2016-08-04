package us_zipcode

import (
	"encoding/json"
	"strconv"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type BatchFixture struct {
	*gunit.Fixture
}

func (f *BatchFixture) TestBatchKnowsWhenItsFullAndEmpty() {
	batch := NewBatch()

	for x := 0; x < MaxBatchSize; x++ {
		f.So(batch.IsFull(), should.BeFalse)
		batch.Append(&Lookup{})
	}

	f.So(batch.IsFull(), should.BeTrue)
}

func (f *BatchFixture) TestCapacityIsLimitedAt100Inputs() {
	batch := NewBatch()

	f.So(batch.Length(), should.Equal, 0)
	f.So(batch.Records(), should.HaveLength, 0)

	for x := 0; x < MaxBatchSize; x++ {
		f.So(batch.Append(&Lookup{InputID: strconv.Itoa(x)}), should.BeTrue)
	}
	f.So(batch.Length(), should.Equal, MaxBatchSize)
	f.So(batch.Records(), should.HaveLength, MaxBatchSize)

	for x := 100; x < 200; x++ {
		f.So(batch.Append(&Lookup{InputID: strconv.Itoa(x)}), should.BeFalse)
	}

	f.So(batch.Length(), should.Equal, MaxBatchSize)
	f.So(batch.Records(), should.HaveLength, MaxBatchSize)
}

func (f *BatchFixture) TestJSONSerializationShouldNeverFail() {
	batch := NewBatch()
	batch.Append(&Lookup{
		City:          "ensure",
		State:         "serialization",
		ZIPCode:       "always",
		InputID:       "successful",
	})
	serialized, err := json.Marshal(batch.lookups)
	f.So(err, should.BeNil)
	f.So(serialized, should.NotBeEmpty)
}

func (f *BatchFixture) TestClearRemovesAllRecords() {
	batch := NewBatch()
	for x := 0; x < MaxBatchSize; x++ {
		f.So(batch.Append(&Lookup{InputID: strconv.Itoa(x)}), should.BeTrue)
	}

	batch.Clear()

	f.So(batch.Length(), should.Equal, 0)
}
