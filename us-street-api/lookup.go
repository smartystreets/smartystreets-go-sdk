package us_street

type Lookup struct {
	Street        string `json:"street,omitempty"`
	Street2       string `json:"street2,omitempty"`
	Secondary     string `json:"secondary,omitempty"`
	City          string `json:"city,omitempty"`
	State         string `json:"state,omitempty"`
	ZIPCode       string `json:"zipcode,omitempty"`
	LastLine      string `json:"lastline,omitempty"`
	Addressee     string `json:"addressee,omitempty"`
	Urbanization  string `json:"urbanization,omitempty"`
	InputID       string `json:"input_id,omitempty"`
	MaxCandidates int    `json:"candidates,omitempty"`

	Results []*Candidate `json:"-"`
}
