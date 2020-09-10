package us_reverse_geo

import "github.com/smartystreets/smartystreets-go-sdk"

// Response structure defined here: https://smartystreets.com/docs/cloud/us-reverse-geo-api#http-response-output
type Response struct {
	Results []Result `json:"results"`
}

type Result struct {
	Coordinate Coordinate `json:"coordinate"`
	Address    Address    `json:"address"`
	Distance   float64    `json:"distance"`
}

type Coordinate struct {
	Latitude  float64               `json:"latitude"`
	Longitude float64               `json:"longitude"`
	Accuracy  string                `json:"accuracy"`
	License   sdk.CoordinateLicense `json:"license"`
}

type Address struct {
	Street            string `json:"street"`
	City              string `json:"city"`
	StateAbbreviation string `json:"state_abbreviation"`
	ZIPCode           string `json:"zipcode"`
}
