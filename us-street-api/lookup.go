package street

import "fmt"

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
	MatchStrategy matchStrategy `json:"match,omitempty"`

	Results []*Candidate `json:"-"`
}

/**************************************************************************/

type matchStrategy int

const (
	MatchDefault matchStrategy = iota
	MatchStrict
	MatchRange
	MatchInvalid
)

func (this matchStrategy) MarshalJSON() ([]byte, error) {
	var value string
	switch this {
	case MatchDefault:
		value = matchDefaultValue
	case MatchStrict:
		value = matchStrictValue
	case MatchRange:
		value = matchRangeValue
	case MatchInvalid:
		value = matchInvalidValue
	default:
		value = fmt.Sprintf("\"%d\"", this)
	}
	return []byte(value), nil
}

const (
	matchDefaultValue = `""`
	matchStrictValue  = `"strict"`
	matchRangeValue   = `"range"`
	matchInvalidValue = `"invalid"`
)
