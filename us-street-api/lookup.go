package street

import (
	"net/url"
	"strconv"
)

// Lookup contains all input fields defined here:
// https://smartystreets.com/docs/us-street-api#input-fields
type Lookup struct {
	Street        string        `json:"street,omitempty"`
	Street2       string        `json:"street2,omitempty"`
	Secondary     string        `json:"secondary,omitempty"`
	City          string        `json:"city,omitempty"`
	State         string        `json:"state,omitempty"`
	ZIPCode       string        `json:"zipcode,omitempty"`
	LastLine      string        `json:"lastline,omitempty"`
	Addressee     string        `json:"addressee,omitempty"`
	Urbanization  string        `json:"urbanization,omitempty"`
	InputID       string        `json:"input_id,omitempty"`
	MaxCandidates int           `json:"candidates,omitempty"` // Default value: 1
	MatchStrategy MatchStrategy `json:"match,omitempty"`
	OutputFormat  OutputFormat  `json:"format,omitempty"`

	Results []*Candidate `json:"results,omitempty"`
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
	if l.MaxCandidates > 0 {
		encode(query, strconv.Itoa(l.MaxCandidates), "candidates")
	} else if l.MatchStrategy == MatchEnhanced {
		encode(query, "5", "candidates")
	}
	if l.MatchStrategy != MatchStrict {
		encode(query, string(l.MatchStrategy), "match")
	}
	if l.OutputFormat != FormatDefault {
		encode(query, string(l.OutputFormat), "format")
	}
}
func encode(query url.Values, source string, target string) {
	if source != "" {
		query.Set(target, source)
	}
}

/**************************************************************************/

type MatchStrategy string

const (
	MatchStrict   = MatchStrategy("strict")
	MatchRange    = MatchStrategy("range") // Deprecated
	MatchInvalid  = MatchStrategy("invalid")
	MatchEnhanced = MatchStrategy("enhanced")
)

type OutputFormat string

const (
	FormatDefault    = OutputFormat("default")
	FormatProjectUSA = OutputFormat("project-usa")
	FormatCASS       = OutputFormat("cass")
)
