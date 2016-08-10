package street

import (
	"encoding/json"
	"testing"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func TestMatchStrategy_MarshalJSON(t *testing.T) {
	assert := assertions.New(t)
	assert.So(serialize(MatchDefault), should.Equal, `{}`)
	assert.So(serialize(MatchInvalid), should.Equal, `{"match":"invalid"}`)
	assert.So(serialize(MatchStrict), should.Equal, `{"match":"strict"}`)
	assert.So(serialize(MatchRange), should.Equal, `{"match":"range"}`)
}

func serialize(strategy matchStrategy) string {
	lookup := Lookup{MatchStrategy: strategy}
	content, _ := json.Marshal(lookup)
	return string(content)
}
