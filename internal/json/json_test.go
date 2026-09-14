package json

import (
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestFixture(t *testing.T) {
	gunit.Run(new(Fixture), t)
}

type Fixture struct {
	*gunit.Fixture
}

type record struct {
	DeliveryLine string `json:"delivery_line"`
	Notes        string `json:"notes"`
}

func (this *Fixture) TestFieldNamesMatchCaseInsensitively() {
	var r record
	err := Unmarshal([]byte(`{"Delivery_Line":"1 Main St"}`), &r)
	this.So(err, should.BeNil)
	this.So(r.DeliveryLine, should.Equal, "1 Main St")
}

func (this *Fixture) TestInvalidUTF8IsAccepted() {
	var r record
	err := Unmarshal([]byte("{\"notes\":\"caf\xe9\"}"), &r)
	this.So(err, should.BeNil)
	this.So(r.Notes, should.StartWith, "caf")
}

func (this *Fixture) TestDuplicateNamesTakeTheLastValue() {
	var r record
	err := Unmarshal([]byte(`{"notes":"first","notes":"second"}`), &r)
	this.So(err, should.BeNil)
	this.So(r.Notes, should.Equal, "second")
}

func (this *Fixture) TestMarshalIsDeterministic() {
	m := map[string]int{"zulu": 1, "alpha": 2, "mike": 3}
	first, err := Marshal(m)
	this.So(err, should.BeNil)
	this.So(string(first), should.Equal, `{"alpha":2,"mike":3,"zulu":1}`)
	for range 20 {
		again, _ := Marshal(m)
		this.So(string(again), should.Equal, string(first))
	}
}
