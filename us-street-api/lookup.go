package street

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

	Results []*Candidate `json:"results,omitempty"`
}

/**************************************************************************/

type MatchStrategy string

const (
	MatchStrict  = MatchStrategy("strict")
	MatchRange   = MatchStrategy("range")
	MatchInvalid = MatchStrategy("invalid")
)
