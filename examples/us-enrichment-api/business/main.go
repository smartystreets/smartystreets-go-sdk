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
		wireup.BasicAuthCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")),
	)

	// Documentation for input fields can be found at:
	// https://www.smarty.com/docs/cloud/us-address-enrichment-api#http-request-input-fields

	smartyKey := "1962995076"

	// Step 1: Send a business summary lookup to get the list of businesses at this address
	summaryLookup := us_enrichment.Lookup{
		SmartyKey: smartyKey,
	}

	err, summaryResults := client.SendBusinessSummary(&summaryLookup)
	if err != nil {
		if client.IsHTTPErrorCode(err, http.StatusNotModified) {
			log.Printf("Record has not been modified since the last request")
			return
		}
		log.Fatal("Error sending summary lookup:", err)
	}

	if len(summaryResults) == 0 || len(summaryResults[0].Businesses) == 0 {
		log.Fatal("No businesses found for this SmartyKey")
	}

	fmt.Printf("Summary results for SmartyKey: %s\n", smartyKey)
	for _, biz := range summaryResults[0].Businesses {
		fmt.Printf("  - %s (ID: %s)\n", biz.CompanyName, biz.BusinessID)
	}

	// Step 2: Use the first business ID to get detailed information
	businessID := summaryResults[0].Businesses[0].BusinessID
	fmt.Printf("\nFetching details for business: %s (ID: %s)\n", summaryResults[0].Businesses[0].CompanyName, businessID)

	detailLookup := us_enrichment.Lookup{
		SmartyKey: smartyKey,
		ETag:      "", // optional: check if the record has been updated
	}

	err, detailResults := client.SendBusinessDetail(&detailLookup, businessID)
	if err != nil {
		if client.IsHTTPErrorCode(err, http.StatusNotModified) {
			log.Printf("Record has not been modified since the last request")
			return
		}
		log.Fatal("Error sending detail lookup:", err)
	}

	fmt.Println("\nDetail results:")
	for s, response := range detailResults {
		jsonResponse, _ := json.MarshalIndent(response, "", "     ")
		fmt.Printf("#%d: %s\n", s, string(jsonResponse))
	}

	log.Println("OK")
}
