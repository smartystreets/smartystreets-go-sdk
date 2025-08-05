package contracts

import (
	"flag"
)

const (
	flagAuthIDVariable = "id"
	flagAuthRefererVariable = "referer"
	flagAuthTokenVariable = "token"
	flagAuthWebVariable = "web"
	flagCustomBaseURL = "baseurl"

	defaultAuthIDVariable = ""
	defaultAuthRefererVariable = ""
	defaultAuthTokenVariable = ""
	defaultAuthWebVariable = ""
	defaultCustomBaseURL = ""
)

type Settings struct {
	AuthIDVariable string
	AuthWebVariable string
	AuthTokenVariable string
	AuthRefererVariable string
	Command []string
	CustomBaseURL string
}

func GetSettings(args []string, commandLine *flag.FlagSet) Settings {
	settings := Settings{
		Command: []string{CommandOptions},
	}

	commandLine.StringVar(
		&settings.AuthIDVariable,
		flagAuthIDVariable,
		defaultAuthIDVariable,
		"Name of the environment variable which contains your auth id. Use this option if you are using secret authorization",
	)
	commandLine.StringVar(
		&settings.AuthWebVariable,
		flagAuthWebVariable,
		defaultAuthWebVariable,
		"Name of the environment variable which contains your auth web. Use this option if you are using web authorization",
	)
	commandLine.StringVar(
		&settings.AuthTokenVariable,
		flagAuthTokenVariable,
		defaultAuthTokenVariable,
		"Name of the environment variable which contains your auth token. Use this option if you are using secret authorization",
	)
	commandLine.StringVar(
		&settings.AuthRefererVariable,
		flagAuthRefererVariable,
		defaultAuthRefererVariable,
		"Name of the environment variable which contains your auth referer. Use this option if you are using web authorization",
	)
	commandLine.StringVar(
		&settings.CustomBaseURL,
		flagCustomBaseURL,
		defaultCustomBaseURL,
		"Custom base url.",
	)

	commandLine.Parse(args)

	settings.Command = commandLine.Args()
	return settings
}
