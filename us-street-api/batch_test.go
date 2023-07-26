package street

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestBatchFixture(t *testing.T) {
	gunit.Run(new(BatchFixture), t)
}

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
