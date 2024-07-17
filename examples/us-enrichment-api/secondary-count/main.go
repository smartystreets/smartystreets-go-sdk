package main

import (
	"encoding/json"
	"fmt"
	usenrich "github.com/smartystreets/smartystreets-go-sdk/us-enrichment-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetFlags(log.Ltime | log.Llongfile)

	client := wireup.BuildUSEnrichmentAPIClient(
		//wireup.WebsiteKeyCredential(os.Getenv("SMARTY_AUTH_WEB"), os.Getenv("SMARTY_AUTH_REFERER")),
		wireup.SecretKeyCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")),
	)

	// Documentation for input fields can be found at:
	// https://www.smarty.com/docs/cloud/us-address-enrichment-api#http-request-input-fields

	smartyKey := "1270119982"

	lookup := usenrich.Lookup{
		SmartyKey: smartyKey,
		ETag:      "", // optional: check if the record has been updated
	}

	err, results := client.SendSecondaryCountLookup(&lookup)

	if err != nil {
		// If ETag was supplied in the lookup, this status will be returned if the ETag value for the record is current
		if client.IsHTTPErrorCode(err, http.StatusNotModified) {
			log.Printf("Record has not been modified since the last request")
			return
		}
		log.Fatal("Error sending lookup:", err)
	}

	fmt.Printf("Results for input: (%s, %s)\n", smartyKey, "secondary")
	for i, response := range results {
		prettyPrinted, err := json.MarshalIndent(response, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("Response %d: %s", i, prettyPrinted)
	}
}
