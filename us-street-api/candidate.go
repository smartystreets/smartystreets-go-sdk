package street

import sdk "github.com/smartystreets/smartystreets-go-sdk/v2"

type (
	// Candidate contains all output fields defined here:
	// https://smartystreets.com/docs/us-street-api#http-response-output
	Candidate struct {
		InputID              string     `json:"input_id,omitempty"`
		InputIndex           int        `json:"input_index"`
		CandidateIndex       int        `json:"candidate_index"`
		Addressee            string     `json:"addressee,omitempty"`
		DeliveryLine1        string     `json:"delivery_line_1,omitempty"`
		DeliveryLine2        string     `json:"delivery_line_2,omitempty"`
		LastLine             string     `json:"last_line,omitempty"`
		DeliveryPointBarcode string     `json:"delivery_point_barcode,omitempty"`
		SmartyKey            string     `json:"smarty_key,omitempty"`
		SmartyKeyEXT         string     `json:"smarty_key_ext,omitempty"`
		Components           Components `json:"components"`
		Metadata             Metadata   `json:"metadata"`
		Analysis             Analysis   `json:"analysis"`
	}

	// Components contains all output fields defined here:
	// https://smartystreets.com/docs/us-street-api#components
	Components struct {
		PrimaryNumber            string `json:"primary_number,omitempty"`
		StreetPredirection       string `json:"street_predirection,omitempty"`
		StreetName               string `json:"street_name,omitempty"`
		StreetPostdirection      string `json:"street_postdirection,omitempty"`
		StreetSuffix             string `json:"street_suffix,omitempty"`
		SecondaryNumber          string `json:"secondary_number,omitempty"`
		SecondaryDesignator      string `json:"secondary_designator,omitempty"`
		ExtraSecondaryNumber     string `json:"extra_secondary_number,omitempty"`
		ExtraSecondaryDesignator string `json:"extra_secondary_designator,omitempty"`
		PMBNumber                string `json:"pmb_number,omitempty"`
		PMBDesignator            string `json:"pmb_designator,omitempty"`
		CityName                 string `json:"city_name,omitempty"`
		DefaultCityName          string `json:"default_city_name,omitempty"`
		StateAbbreviation        string `json:"state_abbreviation,omitempty"`
		ZIPCode                  string `json:"zipcode,omitempty"`
		Plus4Code                string `json:"plus4_code,omitempty"`
		DeliveryPoint            string `json:"delivery_point,omitempty"`
		DeliveryPointCheckDigit  string `json:"delivery_point_check_digit,omitempty"`
		Urbanization             string `json:"urbanization,omitempty"`
	}

	// Metadata contains all output fields defined here:
	// https://smartystreets.com/docs/us-street-api#metadata
	Metadata struct {
		RecordType               string                `json:"record_type,omitempty"`
		ZIPType                  string                `json:"zip_type,omitempty"`
		CountyFIPS               string                `json:"county_fips,omitempty"`
		CountyName               string                `json:"county_name,omitempty"`
		CarrierRoute             string                `json:"carrier_route,omitempty"`
		CongressionalDistrict    string                `json:"congressional_district,omitempty"`
		BuildingDefaultIndicator string                `json:"building_default_indicator,omitempty"`
		RDI                      string                `json:"rdi,omitempty"`
		ELOTSequence             string                `json:"elot_sequence,omitempty"`
		ELOTSort                 string                `json:"elot_sort,omitempty"`
		Latitude                 float64               `json:"latitude,omitzero"`
		Longitude                float64               `json:"longitude,omitzero"`
		CoordinateLicense        sdk.CoordinateLicense `json:"coordinate_license,omitzero"`
		Precision                string                `json:"precision,omitempty"`
		TimeZone                 string                `json:"time_zone,omitempty"`
		UTCOffset                float32               `json:"utc_offset,omitzero"`
		DST                      bool                  `json:"dst,omitzero"`
		IANATimeZone             string                `json:"iana_time_zone,omitempty"`
		IANAUTCOffset            float32               `json:"iana_utc_offset,omitzero"`
		IANADST                  bool                  `json:"iana_dst,omitzero"`
		EWSMatch                 bool                  `json:"ews_match,omitzero"`
	}

	// Analysis contains all output fields defined here:
	// https://smartystreets.com/docs/us-street-api#analysis
	Analysis struct {
		DPVMatchCode      string            `json:"dpv_match_code,omitempty"`
		DPVFootnotes      string            `json:"dpv_footnotes,omitempty"`
		DPVCMRACode       string            `json:"dpv_cmra,omitempty"`
		DPVVacantCode     string            `json:"dpv_vacant,omitempty"`
		DPVNoStat         string            `json:"dpv_no_stat,omitempty"`
		Active            string            `json:"active,omitempty"`
		Footnotes         string            `json:"footnotes,omitempty"` // https://smartystreets.com/docs/us-street-api#footnotes
		LACSLinkCode      string            `json:"lacslink_code,omitempty"`
		LACSLinkIndicator string            `json:"lacslink_indicator,omitempty"`
		SuiteLinkMatch    bool              `json:"suitelink_match,omitzero"`
		EWSMatch          bool              `json:"ews_match,omitzero"`       // deprecated
		EnhancedMatch     string            `json:"enhanced_match,omitempty"` //v2 integration
		Components        ComponentAnalysis `json:"components"`
	}

	ComponentAnalysis struct {
		PrimaryNumber            MatchInfo `json:"primary_number"`
		StreetPredirection       MatchInfo `json:"street_predirection"`
		StreetName               MatchInfo `json:"street_name"`
		StreetPostdirection      MatchInfo `json:"street_postdirection"`
		StreetSuffix             MatchInfo `json:"street_suffix"`
		SecondaryNumber          MatchInfo `json:"secondary_number"`
		SecondaryDesignator      MatchInfo `json:"secondary_designator"`
		ExtraSecondaryNumber     MatchInfo `json:"extra_secondary_number"`
		ExtraSecondaryDesignator MatchInfo `json:"extra_secondary_designator"`
		CityName                 MatchInfo `json:"city_name"`
		StateAbbreviation        MatchInfo `json:"state_abbreviation"`
		ZIPCode                  MatchInfo `json:"zipcode"`
		Plus4Code                MatchInfo `json:"plus4_code"`
		Urbanization             MatchInfo `json:"urbanization"`
	}

	MatchInfo struct {
		Status string   `json:"status,omitempty"`
		Change []string `json:"change,omitempty"`
	}
)
