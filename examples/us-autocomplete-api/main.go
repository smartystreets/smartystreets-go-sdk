package main

import (
	"fmt"
	"log"
	"os"

	"github.com/smartystreets/smartystreets-go-sdk/us-autocomplete-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

func main() {
	log.SetFlags(log.Ltime)

	client := wireup.NewClientBuilder().
		WithSecretKeyCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")).
		WithDebugHTTPOutput(). // uncomment this line to see detailed HTTP request/response information.
		BuildUSAutocompleteAPIClient()

	lookup := &autocomplete.Lookup{Prefix: "123 main"}

	if err := client.SendLookup(lookup); err != nil {
		log.Fatal("Error sending batch:", err)
	}

	fmt.Println("Results for input:")
	fmt.Println()
	for s, suggestion := range lookup.Results {
		fmt.Println("  Suggestion:", s)
		fmt.Println(" ", suggestion.Text)
		fmt.Println(" ", suggestion.StreetLine)
		fmt.Println(" ", suggestion.City)
		fmt.Println(" ", suggestion.State)
		fmt.Println()
	}
}
