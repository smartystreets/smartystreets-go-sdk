package main

import (
	"cmd/main/contracts"
	"cmd/main/runners"
	"cmd/main/tap14"
	"flag"
	"log"
	"os"

	street "github.com/smartystreets/smartystreets-go-sdk/us-street-api"
)

func main() {
	log.SetFlags(0)

	emitter := tap14.Emitter{}
	emitter.StartTestGroupStreamed(os.Stdout)
	defer emitter.CompleteTestGroup()

	settings := contracts.GetSettings(os.Args[1:], flag.CommandLine)
	emitter.TestBool(true, tap14.AddLabel("get settings"), tap14.SetData(settings))

	requests := []*street.Lookup{
		{
			InputID:       "24601", // Optional ID from your system
			Addressee:     "John Doe",
			Street:        "1 Rosedale",
			Street2:       "closet under the stairs",
			Secondary:     "APT 2",
			Urbanization:  "", // Only applies to Puerto Rico addresses
			City:          "Baltimore",
			State:         "MD",
			ZIPCode:       "21229",
			MaxCandidates: 3,
			MatchStrategy: street.MatchInvalid,
		},
		{
			Street:        "1600 Pennsylvania Avenue",
			LastLine:      "Washington, DC",
			MaxCandidates: 5,
		},
		{
			InputID:       "8675309",
			Street:        "1600 Amphitheatre Parkway Mountain View, CA 94043",
			MaxCandidates: 1,
		},
	}

	runner := &runners.USStreetRunner{}
	emitter.TestError(runner.Setup(settings), tap14.AddLabel("setup runner"))
	emitter.TestError(runner.Run(requests), tap14.AddLabel("run runner"))
}
