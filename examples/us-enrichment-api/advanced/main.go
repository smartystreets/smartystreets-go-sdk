package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	sdk "github.com/smartystreets/smartystreets-go-sdk"
	us_enrichment "github.com/smartystreets/smartystreets-go-sdk/us-enrichment-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
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
	// https://www.smarty.com/docs/cloud/us-address-enrichment-api#http-request-input-fields

	smartyKey := "1682393594"

	lookup := us_enrichment.Lookup{
		SmartyKey: smartyKey,
		Include:   "group_structural,sale_date",
		Exclude:   "",
		ETag:      "GU4TINZRHA4TQMY",
	}

	err, results := client.SendPropertyPrincipalWithLookup(&lookup)

	if err != nil {
		//httpError := &sdk.HTTPStatusError{}
		//if errors.As(err, &httpError) && err.(*sdk.HTTPStatusError).StatusCode() == http.StatusNotModified {
		if serr, ok := err.(*sdk.HTTPStatusError); ok && serr.StatusCode() == http.StatusNotModified {
			log.Printf("Record has not been modified since the last request")
			return
		}
		log.Fatal("Error sending lookup:", err)
	}

	fmt.Printf("Results for input: (%s, %s)\n", smartyKey, "principal")
	for s, response := range results {
		fmt.Printf("#%d: %+v\n", s, response)
	}

	log.Println("OK")
}
