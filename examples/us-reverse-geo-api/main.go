package main

import (
	"context"
	"fmt"
	"log"
	"os"

	us_reverse_geo "github.com/smartystreets/smartystreets-go-sdk/us-reverse-geo-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

func main() {
	log.SetFlags(log.Ltime | log.Llongfile)

	client := wireup.BuildUSReverseGeocodingAPIClient(
		wireup.SecretKeyCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")),
		// The appropriate license values to be used for your subscriptions
		// can be found on the Subscriptions page the account dashboard.
		// https://www.smartystreets.com/docs/cloud/licensing
		wireup.WithLicenses("us-reverse-geocoding-cloud"),
		// wireup.DebugHTTPOutput(), // uncomment this line to see detailed HTTP request/response information.
	)

	// Documentation for input fields can be found at:
	// https://smartystreets.com/docs/cloud/us-reverse-geo-api#http-request-input-fields

	lookup := &us_reverse_geo.Lookup{
		Latitude:  40.27644,
		Longitude: -111.65747,
	}

	if err := client.SendLookupWithContext(context.Background(), lookup); err != nil {
		log.Fatal("Error sending batch:", err)
	}

	fmt.Printf("Results for input: (%f, %f)\n", lookup.Latitude, lookup.Longitude)
	for s, address := range lookup.Response.Results {
		fmt.Printf("#%d: %#v\n", s, address)
	}

	log.Println("OK")
}
