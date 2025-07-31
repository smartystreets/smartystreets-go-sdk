package main

import (
	"flag"
	"fmt"
)

const (
	flagAuthIDVariable = "id"
	flagAuthTokenVariable = "token"
	flagCustomBaseURL = "baseurl"

	defaultAuthIDVariable = "SMARTY_AUTH_ID"
	defaultAuthTokenVariable = "SMARTY_AUTH_TOKEN"
	defaultCustomBaseURL = ""
)

type Settings struct {
	AuthIDVariable string
	AuthTokenVariable string
	CustomBaseURL string
}

func getSettings(args []string, commandLine *flag.FlagSet) Settings {
	settings := Settings{}

	commandLine.StringVar(
		&settings.AuthIDVariable,
		flagAuthIDVariable,
		defaultAuthIDVariable,
		fmt.Sprintf(
			"Name of the environment variable which contains your auth id. If left blank, will default to %q",
			defaultAuthIDVariable,
		),
	)
	commandLine.StringVar(
		&settings.AuthTokenVariable,
		flagAuthTokenVariable,
		defaultAuthTokenVariable,
		fmt.Sprintf(
			"Name of the environment variable which contains your auth token. If left blank, will default to %q",
			defaultAuthTokenVariable,
		),
	)
	commandLine.StringVar(
		&settings.CustomBaseURL,
		flagCustomBaseURL,
		defaultCustomBaseURL,
		"Custom base url.",
	)
	
	commandLine.Parse(args)
	return settings
}
