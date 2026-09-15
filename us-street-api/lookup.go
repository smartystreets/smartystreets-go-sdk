package street

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"maps"
	"net/url"
	"slices"
	"strconv"
)

// Lookup contains all input fields defined here:
// https://smartystreets.com/docs/us-street-api#input-fields
type Lookup struct {
	Street           string            `json:"street,omitempty"`
	Street2          string            `json:"street2,omitempty"`
	Secondary        string            `json:"secondary,omitempty"`
	City             string            `json:"city,omitempty"`
	State            string            `json:"state,omitempty"`
	ZIPCode          string            `json:"zipcode,omitempty"`
	LastLine         string            `json:"lastline,omitempty"`
	Addressee        string            `json:"addressee,omitempty"`
	Urbanization     string            `json:"urbanization,omitempty"`
	InputID          string            `json:"input_id,omitempty"`
	MaxCandidates    int               `json:"candidates,omitzero"` // Default value: 1, if MatchStrategy is "enhanced" default value: 5
	MatchStrategy    MatchStrategy     `json:"match,omitempty"`     // Default value: "enhanced"
	OutputFormat     OutputFormat      `json:"format,omitempty"`
	CountySource     CountySource      `json:"county_source,omitempty"`
	CustomParameters map[string]string `json:"-"`

	Results []*Candidate `json:"results,omitempty"`
}

// AddCustomParameter adds custom query parameters/json properties to the request, it will overwrite
// existing key value pairs.
func (l *Lookup) AddCustomParameter(name, value string) {
	if l.CustomParameters == nil {
		l.CustomParameters = make(map[string]string)
	}
	l.CustomParameters[name] = value
}

func (l *Lookup) encodeQueryString(query url.Values) {
	encode(query, l.Street, "street")
	encode(query, l.Street2, "street2")
	encode(query, l.Secondary, "secondary")
	encode(query, l.City, "city")
	encode(query, l.State, "state")
	encode(query, l.ZIPCode, "zipcode")
	encode(query, l.LastLine, "lastline")
	encode(query, l.Addressee, "addressee")
	encode(query, l.Urbanization, "urbanization")
	encode(query, l.InputID, "input_id")
	encode(query, string(l.OutputFormat), "format")
	encode(query, string(l.CountySource), "county_source")

	matchStrategy, maxCandidates := l.defaultValues()
	encode(query, string(matchStrategy), "match")
	if maxCandidates > 0 {
		encode(query, strconv.Itoa(maxCandidates), "candidates")
	}

	for k, v := range l.CustomParameters {
		encode(query, v, k)
	}
}

func encode(query url.Values, source string, target string) {
	if source != "" {
		query.Set(target, source)
	}
}

// MarshalJSONTo implements [json.MarshalerTo]. The body carries the same defaults
// as the query string (see defaultValues), and each custom parameter is written as
// a string member after the named fields. A custom parameter whose name matches a
// named field replaces it, mirroring url.Values.Set in encodeQueryString.
func (l *Lookup) MarshalJSONTo(enc *jsontext.Encoder) error {
	type alias Lookup
	lc := alias(*l)
	lc.MatchStrategy, lc.MaxCandidates = l.defaultValues()
	if len(l.CustomParameters) == 0 {
		return json.MarshalEncode(enc, &lc)
	}

	fields, err := json.Marshal(&lc, enc.Options())
	if err != nil {
		return err
	}
	dec := jsontext.NewDecoder(bytes.NewReader(fields))
	if _, err := dec.ReadToken(); err != nil { // consume '{'
		return err
	}
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	for dec.PeekKind() != '}' {
		name, err := dec.ReadToken()
		if err != nil {
			return err
		}
		if _, replaced := l.CustomParameters[name.String()]; replaced {
			if err := dec.SkipValue(); err != nil {
				return err
			}
			continue
		}
		if err := enc.WriteToken(name); err != nil {
			return err
		}
		value, err := dec.ReadValue()
		if err != nil {
			return err
		}
		if err := enc.WriteValue(value); err != nil {
			return err
		}
	}
	for _, name := range slices.Sorted(maps.Keys(l.CustomParameters)) {
		if err := enc.WriteToken(jsontext.String(name)); err != nil {
			return err
		}
		if err := enc.WriteToken(jsontext.String(l.CustomParameters[name])); err != nil {
			return err
		}
	}
	return enc.WriteToken(jsontext.EndObject)
}

func (l *Lookup) defaultValues() (matchStrategy MatchStrategy, maxCandidates int) {
	matchStrategy = l.MatchStrategy
	if matchStrategy == "" {
		matchStrategy = MatchEnhanced
	}
	maxCandidates = l.MaxCandidates
	if maxCandidates == 0 && matchStrategy == MatchEnhanced {
		maxCandidates = 5
	}
	return matchStrategy, maxCandidates
}

/**************************************************************************/

type MatchStrategy string

const (
	MatchStrict   = MatchStrategy("strict")
	MatchInvalid  = MatchStrategy("invalid")
	MatchEnhanced = MatchStrategy("enhanced")
)

type OutputFormat string

const (
	FormatDefault    = OutputFormat("default")
	FormatProjectUSA = OutputFormat("project-usa")
)

type CountySource string

const (
	PostalCounty     = CountySource("postal")
	GeographicCounty = CountySource("geographic")
)
