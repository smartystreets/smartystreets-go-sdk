package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	us_enrichment "github.com/smartystreets/smartystreets-go-sdk/us-enrichment-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

func main() {
	log.SetFlags(log.Ltime | log.Llongfile)

	client := wireup.BuildUSEnrichmentAPIClient(
		//wireup.WebsiteKeyCredential(os.Getenv("SMARTY_AUTH_WEB"), os.Getenv("SMARTY_AUTH_REFERER")),
		wireup.HeaderCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")),
	)

	// Documentation for input fields can be found at:
	// https://www.smarty.com/docs/cloud/us-address-enrichment-api#http-request-input-fields

	smartyKey := "87844267"

	lookup := us_enrichment.Lookup{
		SmartyKey: smartyKey, //smartyKey,
		Include:   "",        // optional: only include these attributes in the returned data. e.g. "group_structural,sale_date"
		Exclude:   "",        // optional: exclude attributes from the returned data
		ETag:      "",        // optional: check if the record has been updated
		Features:  "financial",
	}

	err, results := client.SendPropertyPrincipal(&lookup)

	if err != nil {
		// If ETag was supplied in the lookup, this status will be returned if the ETag value for the record is current
		if client.IsHTTPErrorCode(err, http.StatusNotModified) {
			log.Printf("Record has not been modified since the last request")
			return
		}
		log.Fatal("Error sending lookup:", err)
	}

	fmt.Printf("Results for input: (%s, %s)\n", smartyKey, "principal")
	for s, response := range results {
		jsonResponse, _ := json.MarshalIndent(response, "", "     ")
		fmt.Printf("#%d: %s\n", s, string(jsonResponse))
	}
}
