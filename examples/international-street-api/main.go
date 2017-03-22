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

	client := wireup.NewClientBuilder().
		WithSecretKeyCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")).
		//WithDebugHTTPOutput(). // uncomment this line to see detailed HTTP request/response information.
		BuildInternationalStreetAPIClient()

	lookup := &international.Lookup{
		Address1:           "Rua Padre Antônio D'angelo 121",
		Locality:           "Sao Paulo",
		AdministrativeArea: "SP",
		Country:            "Brazil",
	}

	if err := client.SendLookup(lookup); err != nil {
		log.Fatal("Error sending batch:", err)
	}

	for s, suggestion := range lookup.Results {
		fmt.Println("Result:       ", s)
		Println("Organization: ", suggestion.Organization)
		Println("Address1:     ", suggestion.Address1)
		Println("Address2:     ", suggestion.Address2)
		Println("Address3:     ", suggestion.Address3)
		Println("Address4:     ", suggestion.Address4)
		Println("Address5:     ", suggestion.Address5)
		Println("Address6:     ", suggestion.Address6)
		Println("Address7:     ", suggestion.Address7)
		Println("Address8:     ", suggestion.Address8)
		Println("Address9:     ", suggestion.Address9)
		Println("Address10:    ", suggestion.Address10)
		Println("Address11:    ", suggestion.Address11)
		Println("Address12:    ", suggestion.Address12)
	}

	log.Println("OK")
}

func Println(name, value string) {
	if len(value) > 0 {
		fmt.Println(name, value)
	}
}
