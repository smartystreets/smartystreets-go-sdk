package main

import (
	"fmt"
	"log"
	"os"

	"github.com/smartystreets/smartystreets-go-sdk/international-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

func main() {
	log.SetFlags(log.Ltime | log.Llongfile)

	client := wireup.BuildInternationalStreetAPIClient(
		wireup.SecretKeyCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")),
		//wireup.DebugHTTPOutput(), // uncomment this line to see detailed HTTP request/response information.
	)

	lookup := &street.Lookup{
		Address1:           "Rua Padre Antonio D'Angelo 121",
		Address2:           "Casa Verde",
		Locality:           "Sao Paulo",
		AdministrativeArea: "SP",
		Country:            "Brazil",
	}

	if err := client.SendLookup(lookup); err != nil {
		log.Fatal("Error sending batch:", err)
	}

	fmt.Println("Results:")
	fmt.Println()

	for c, candidate := range lookup.Results {
		fmt.Println("Candidate:", c)
		display(candidate.Address1)
		display(candidate.Address2)
		display(candidate.Address3)
		display(candidate.Address4)
		display(candidate.Address5)
		display(candidate.Address6)
		display(candidate.Address7)
		display(candidate.Address8)
		display(candidate.Address9)
		display(candidate.Address10)
		display(candidate.Address11)
		display(candidate.Address12)
		fmt.Println()
	}

	log.Println("OK")
}

func display(value string) {
	if len(value) > 0 {
		fmt.Printf("  %s\n", value)
	}
}
