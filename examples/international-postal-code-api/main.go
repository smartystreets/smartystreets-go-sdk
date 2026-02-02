package main

import (
	"context"
	"fmt"
	"log"
	"os"

	international_postal_code "github.com/smartystreets/smartystreets-go-sdk/international-postal-code-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

func main() {
	log.SetFlags(log.Ltime | log.Llongfile)

	client := wireup.BuildInternationalPostalCodeAPIClient(
		//wireup.WebsiteKeyCredential(os.Getenv("SMARTY_AUTH_WEB"), os.Getenv("SMARTY_AUTH_REFERER")),
		wireup.BasicAuthCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")),
		// wireup.DebugHTTPOutput(), // uncomment this line to see detailed HTTP request/response information.
	)

	// For complete list of input fields, refer to:
	// https://smartystreets.com/docs/cloud/international-postal-code-api

	lookup := &international_postal_code.Lookup{
		InputID:            "ID-8675309", // Optional ID from your system
		Locality:           "Sao Paulo",
		AdministrativeArea: "SP",
		Country:            "Brazil",
		PostalCode:         "02516",
	}

	if err := client.SendLookupWithContext(context.Background(), lookup); err != nil {
		log.Fatal("Error sending batch:", err)
	}

	fmt.Println("Results:")
	fmt.Println()

	for c, candidate := range lookup.Results {
		fmt.Println("Candidate:", c)
		display(candidate.InputID)
		display(candidate.CountryIso3)
		display(candidate.Locality)
		display(candidate.DependentLocality)
		display(candidate.DoubleDependentLocality)
		display(candidate.SubAdministrativeArea)
		display(candidate.AdministrativeArea)
		display(candidate.SuperAdministrativeArea)
		display(candidate.PostalCode)
		fmt.Println()
	}

	log.Println("OK")
}

func display(value string) {
	if len(value) > 0 {
		fmt.Printf("  %s\n", value)
	}
}
