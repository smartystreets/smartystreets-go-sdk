package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/smartystreets/smartystreets-go-sdk/us-autocomplete-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

/*
This example is for us-autocomplete (V2). It has the same name as a previous product
which has been deprecated since 2022 which we refer to as US Autocomplete Basic.

If you are still using US Autocomplete Basic, this SDK will not work.
*/

func main() {
	log.SetFlags(log.Ltime | log.Llongfile)

	client := wireup.BuildUSAutocompleteAPIClient(
		//wireup.WebsiteKeyCredential(os.Getenv("SMARTY_AUTH_WEB"), os.Getenv("SMARTY_AUTH_REFERER")),
		wireup.BasicAuthCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")),
		// wireup.DebugHTTPOutput(), // uncomment this line to see detailed HTTP request/response information.
	)

	// Documentation for input fields can be found at:
	// https://smartystreets.com/docs/cloud/us-autocomplete-api#http-request-input-fields

	lookup := &autocomplete.Lookup{
		Search:      "200",
		MaxResults:  10,
		CityFilter:  []string{"Denver,Aurora,CO", "Provo,UT"},
		PreferState: []string{"CO"},
		PreferRatio: 3,
		Source:      "all",
	}

	if err := client.SendLookupWithContext(context.Background(), lookup); err != nil {
		log.Fatal("Error sending lookup:", err)
	}

	fmt.Printf("Results for input: [%s]\n", lookup.Search)
	entryID := ""
	addressWithSecondaries := ""
	for s, suggestion := range lookup.Results {
		fmt.Printf("#%d: %#v\n", s, suggestion)
		if suggestion.EntryID != "" {
			addressWithSecondaries = suggestion.StreetLine + " " + suggestion.City + " " + suggestion.State
			entryID = suggestion.EntryID
		}
	}

	// expand the secondaries of a result that has an entryID
	if len(entryID) > 0 {
		lookup.Selected = entryID
		if err := client.SendLookupWithContext(context.Background(), lookup); err != nil {
			log.Fatal("Error sending lookup:", err)
		}

		fmt.Printf("\nSecondaries for: [%s]\n", addressWithSecondaries)
		for s, suggestion := range lookup.Results {
			fmt.Printf("#%d: %#v\n", s, suggestion)
			if suggestion.EntryID != "" {
				entryID = suggestion.EntryID
			}
		}
	}

	log.Println("OK")
}
