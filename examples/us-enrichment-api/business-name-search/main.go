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

	// Step 1: Send a business summary lookup using a business name and city to find matching businesses
	summaryLookup := us_enrichment.Lookup{
		SmartyKey:    "search",
		BusinessName: "delta air",
		City:         "atlanta",
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
		log.Fatal("No businesses found for this business name search")
	}

	fmt.Printf("Summary results for BusinessName: %s, City: %s\n", summaryLookup.BusinessName, summaryLookup.City)
	for _, biz := range summaryResults[0].Businesses {
		fmt.Printf("  - %s (ID: %s)\n", biz.CompanyName, biz.BusinessID)
	}

	// Step 2: Use the first business ID to get detailed information
	businessID := summaryResults[0].Businesses[0].BusinessID
	fmt.Printf("\nFetching details for business: %s (ID: %s)\n", summaryResults[0].Businesses[0].CompanyName, businessID)

	detailLookup := us_enrichment.Lookup{
		BusinessID: businessID,
		ETag:       "", // optional: check if the record has been updated
	}

	err, detailResults := client.SendBusinessDetail(&detailLookup)
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
