package main

import (
	"flag"
	"log"
	"omg-cli/config"
	"omg-cli/omg/app"
	"os"
)

//TODO(jrjohnson): These constants should be detected, generated, or flags
const (
	username          = "foo"
	password          = "foobar"
	decryptionPhrase  = "foobar"
	skipSSLValidation = true
)

func main() {
	mode := app.Mode(app.ConfigureOpsManager)
	flag.Var(&mode, "mode", "BakeImage, ConfigureOpsManager")
	var pivnetApiToken = flag.String("pivnet-api-token", "", "Required for BakeImage. Look for 'API TOKEN' at https://network.pivotal.io/users/dashboard/edit-profile.")
	var terraformState = flag.String("terraform-state-path", "env.json", "Path to terraform output")

	flag.Parse()

	logger := log.New(os.Stderr, "[ONG] ", 0)

	creds := config.OpsManagerCredentials{username, password, decryptionPhrase, skipSSLValidation}
	app, err := app.New(logger, mode, *terraformState, *pivnetApiToken, creds)
	if err != nil {
		logger.Fatal(err)
	}

	if err := app.Run(mode); err != nil {
		logger.Fatal(err)
	}
}
