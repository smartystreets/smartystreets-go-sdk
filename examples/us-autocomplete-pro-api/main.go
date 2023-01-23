package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/smartystreets/smartystreets-go-sdk/us-autocomplete-pro-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

func main() {
	log.SetFlags(log.Ltime | log.Llongfile)

	client := wireup.BuildUSAutocompleteProAPIClient(
		wireup.WebsiteKeyCredential(os.Getenv("SMARTY_AUTH_WEB"), os.Getenv("SMARTY_AUTH_REFERER")),
		//wireup.SecretKeyCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")),
		// The appropriate license values to be used for your subscriptions
		// can be found on the Subscriptions page the account dashboard.
		// https://www.smartystreets.com/docs/cloud/licensing
		wireup.WithLicenses("us-autocomplete-pro-cloud"),
		// wireup.DebugHTTPOutput(), // uncomment this line to see detailed HTTP request/response information.
	)

	// Documentation for input fields can be found at:
	// https://smartystreets.com/docs/cloud/us-autocomplete-api#http-request-input-fields

	lookup := &autocomplete_pro.Lookup{
		Search:      "1042 W Center",
		MaxResults:  5,
		CityFilter:  []string{"Denver", "Orem"},
		StateFilter: []string{"CO", "UT"},
		PreferState: []string{"CO"},
		PreferRatio: 3,
		Geolocation: autocomplete_pro.GeolocateCity,
		Source:      "all",
	}

	if err := client.SendLookupWithContext(context.Background(), lookup); err != nil {
		log.Fatal("Error sending batch:", err)
	}

	fmt.Printf("Results for input: [%s]\n", lookup.Search)
	for s, suggestion := range lookup.Results {
		fmt.Printf("#%d: %#v\n", s, suggestion)
	}

	log.Println("OK")
}
