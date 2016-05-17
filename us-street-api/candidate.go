package us_street

type (
	Candidate struct {
		InputID              string     `json:"input_id"`
		InputIndex           int        `json:"input_index"`
		CandidateIndex       int        `json:"candidate_index"`
		Addressee            string     `json:"addressee"`
		DeliveryLine1        string     `json:"delivery_line_1"`
		DeliveryLine2        string     `json:"delivery_line_2"`
		LastLine             string     `json:"last_line"`
		DeliveryPointBarcode string     `json:"delivery_point_barcode"`
		Components           Components `json:"components"`
		Metadata             Metadata   `json:"metadata"`
		Analysis             Analysis   `json:"analysis"`
	}

	Components struct {
		PrimaryNumber            string `json:"primary_number"`
		StreetPredirection       string `json:"street_predirection"`
		StreetName               string `json:"street_name"`
		StreetPostdirection      string `json:"street_postdirection"`
		StreetSuffix             string `json:"street_suffix"`
		SecondaryNumber          string `json:"secondary_number"`
		SecondaryDesignator      string `json:"secondary_designator"`
		ExtraSecondaryNumber     string `json:"extra_secondary_number"`
		ExtraSecondaryDesignator string `json:"extra_secondary_designator"`
		PMBNumber                string `json:"pmb_number"`
		PMBDesignator            string `json:"pmb_designator"`
		CityName                 string `json:"city_name"`
		DefaultCityName          string `json:"default_city_name"`
		StateAbbreviation        string `json:"state_abbreviation"`
		ZIPCode                  string `json:"zipcode"`
		Plus4Code                string `json:"plus4_code"`
		DeliveryPoint            string `json:"delivery_point"`
		DeliveryPointCheckDigit  string `json:"delivery_point_check_digit"`
		Urbanization             string `json:"urbanization"`
	}

	Metadata struct {
		RecordType               string  `json:"record_type"`
		ZIPType                  string  `json:"zip_type"`
		CountyFIPS               string  `json:"county_fips"`
		CountyName               string  `json:"county_name"`
		CarrierRoute             string  `json:"carrier_route"`
		CongressionalDistrict    string  `json:"congressional_district"`
		BuildingDefaultIndicator string  `json:"building_default_indicator"`
		RDI                      string  `json:"rdi"`
		ELOTSequence             string  `json:"elot_sequence"`
		ELOTSort                 string  `json:"elot_sort"`
		Latitude                 float64 `json:"latitude"`
		Longitude                float64 `json:"longitude"`
		Precision                string  `json:"precision"`
		TimeZone                 string  `json:"time_zone"`
		UTCOffset                float32 `json:"utc_offset"`
		DST                      bool    `json:"dst"`
	}

	Analysis struct {
		DPVMatchCode      string `json:"dpv_match_code"`
		DPVFootnotes      string `json:"dpv_footnotes"`
		DPVCMRACode       string `json:"dpv_cmra"`
		DPVVacantCode     string `json:"dpv_vacant"`
		Active            string `json:"active"`
		Footnotes         string `json:"footnotes"`
		LACSLinkCode      string `json:"lacslink_code"`
		LACSLinkIndicator string `json:"lacslink_indicator"`
		SuiteLinkMatch    bool   `json:"suitelink_match"`
		EWSMatch          bool   `json:"ews_match"`
	}
)
