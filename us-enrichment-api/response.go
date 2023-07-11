package us_enrichment

type Response struct {
	SmartyKey   string       `json:"smarty-key"`
	DataSetName string       `json:"data-set-name"`
	Attributes  []*Attribute `json:"attributes"`
}

type Attribute struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
}
