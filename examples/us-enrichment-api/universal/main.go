package main

import (
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
		wireup.SecretKeyCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")),
	)

	// Documentation for input fields and available datasets can be found at:
	// https://www.smarty.com/docs/cloud/us-address-enrichment-api

	smartyKey := "1682393594"

	lookup := us_enrichment.Lookup{
		SmartyKey: smartyKey,
		//Include:   "group_structural,sale_date", // optional: only include these attributes in the returned data
		Exclude: "", // optional: exclude attributes from the returned data
		ETag:    "", // optional: check if the record has been updated
	}

	// Universal lookup works with all datasets and optional subsets.
	// Returns JSON bytes
	// Note: The DataSubset field can be an empty string for datasets that have no subsets.
	err, results := client.SendUniversalLookup(&lookup, "property", "financial")

	if err != nil {
		// If ETag was supplied in the lookup, this status will be returned if the ETag value for the record is current
		if client.IsHTTPErrorCode(err, http.StatusNotModified) {
			log.Printf("Record has not been modified since the last request")
			return
		}
		log.Fatal("Error sending lookup:", err)
	}

	fmt.Printf("Results for input: (%s, %s)\n", smartyKey, "principal")
	fmt.Println(string(results))
}
