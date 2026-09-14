// Package json is the SDK's single entry point to encoding/json/v2.
//
// Every marshal and unmarshal in the library goes through this package so the
// option set is decided in one place. Unmarshaling keeps the tolerant behaviors
// of encoding/json v1 that matter for data coming back from the API:
// case-insensitive field-name matching, and acceptance of invalid UTF-8 and
// duplicate object names. Marshaling is deterministic so request bodies and
// debug output are stable from run to run.
package json

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
)

var (
	marshalOptions = jsonv2.JoinOptions(
		jsonv2.Deterministic(true),
	)
	unmarshalOptions = jsonv2.JoinOptions(
		jsonv2.MatchCaseInsensitiveNames(true),
		jsontext.AllowInvalidUTF8(true),
		jsontext.AllowDuplicateNames(true),
	)
)

// Marshal returns the JSON encoding of in using the SDK's marshal options.
func Marshal(in any) ([]byte, error) {
	return jsonv2.Marshal(in, marshalOptions)
}

// Unmarshal decodes in into out using the SDK's unmarshal options.
func Unmarshal(in []byte, out any) error {
	return jsonv2.Unmarshal(in, out, unmarshalOptions)
}
