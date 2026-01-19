package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	usenrich "github.com/smartystreets/smartystreets-go-sdk/us-enrichment-api"
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

	//smartyKey := "1270119982"
	smartyKey := "search" // SmartyKey to use when performing an address search

	lookup := usenrich.Lookup{
		SmartyKey: smartyKey,
		Freeform:  "1 E Main Carnegie OK",
		ETag:      "", // optional: check if the record has been updated
	}

	err, results := client.SendRisk(&lookup)

	if err != nil {
		// If ETag was supplied in the lookup, this status will be returned if the ETag value for the record is current
		if client.IsHTTPErrorCode(err, http.StatusNotModified) {
			log.Printf("Record has not been modified since the last request")
			return
		}
		log.Fatal("Error sending lookup:", err)
	}

	fmt.Printf("Results for input: (%s, %s)\n", smartyKey, "risk")
	for i, response := range results {
		prettyPrinted, err := json.MarshalIndent(response, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("Response %d: %s", i, prettyPrinted)
	}
}
