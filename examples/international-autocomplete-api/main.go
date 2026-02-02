package main

import (
	"context"
	"fmt"
	"log"
	"os"

	international_autocomplete "github.com/smartystreets/smartystreets-go-sdk/international-autocomplete-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

func main() {
	log.SetFlags(log.Ltime | log.Llongfile)

	client := wireup.BuildInternationalAutocompleteAPIClient(
		//wireup.WebsiteKeyCredential(os.Getenv("SMARTY_AUTH_WEB"), os.Getenv("SMARTY_AUTH_REFERER")),
		wireup.BasicAuthCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")),
	)

	// Documentation for input fields can be found at:
	// https://smartystreets.com/docs/cloud/us-autocomplete-api#http-request-input-fields

	lookup := &international_autocomplete.Lookup{
		Country:  "FRA",
		Search:   "Louis",
		Locality: "Paris",
	}

	if err := client.SendLookupWithContext(context.Background(), lookup); err != nil {
		log.Fatal("Error sending batch:", err)
	}

	if len(lookup.Search) > 0 {
		fmt.Printf("Results for input: [%s]\n", lookup.Search)
	} else {
		fmt.Printf("Results for input: [%s]\n", lookup.AddressID)
	}
	for s, candidate := range lookup.Result.Candidates {
		fmt.Printf("#%d: %#v\n", s, candidate)
	}

	log.Println("OK")
}
