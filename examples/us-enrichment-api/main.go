package main

import (
	"fmt"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Ltime | log.Llongfile)

	client := wireup.BuildUSEnrichmentAPIClient(
		//wireup.WebsiteKeyCredential(os.Getenv("SMARTY_AUTH_WEB"), os.Getenv("SMARTY_AUTH_REFERER")),
		wireup.SecretKeyCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")),
		// The appropriate license values to be used for your subscriptions
		// can be found on the Subscriptions page the account dashboard.
		// https://www.smarty.com/docs/cloud/licensing
		wireup.WithLicenses("us-property-data-principal-cloud"),
		// wireup.DebugHTTPOutput(), // uncomment this line to see detailed HTTP request/response information.
	)

	// Documentation for input fields can be found at:
	// https://smartystreets.com/docs/cloud/us-reverse-geo-api#http-request-input-fields

	smartyKey := "1682393594"

	err, results := client.SendPropertyPrincipalLookup(smartyKey)

	if err != nil {
		log.Fatal("Error sending lookup:", err)
	}

	fmt.Printf("Results for input: (%s, %s)\n", smartyKey, "principal")
	for s, response := range results {
		fmt.Printf("#%d: %+v\n", s, response)
	}

	log.Println("OK")
}
