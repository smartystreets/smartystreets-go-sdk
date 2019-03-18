package main

import (
	"fmt"
	"log"
	"os"

	"github.com/smartystreets/smartystreets-go-sdk/us-autocomplete-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

func main() {
	log.SetFlags(log.Ltime | log.Llongfile)

	client := wireup.BuildUSAutocompleteAPIClient(
		wireup.SecretKeyCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")),
		// wireup.DebugHTTPOutput(), // uncomment this line to see detailed HTTP request/response information.
	)

	// Documentation for input fields can be found at:
	// https://smartystreets.com/docs/cloud/us-autocomplete-api#http-request-input-fields

	lookup := &autocomplete.Lookup{
		Prefix:         "4770 Lincoln Ave O",
		MaxSuggestions: 5,
		CityFilter:     []string{"Ogden"},
		StateFilter:    []string{"IL"},
		Preferences:    []string{"Fallon, IL"},
		Geolocation:    autocomplete.GeolocateNone,
		PreferRatio:    0.3333333333,
	}

	if err := client.SendLookup(lookup); err != nil {
		log.Fatal("Error sending batch:", err)
	}

	fmt.Printf("Results for input: [%s]\n", lookup.Prefix)
	for s, suggestion := range lookup.Results {
		fmt.Printf("#%d: %#v\n", s, suggestion)
	}

	log.Println("OK")
}
