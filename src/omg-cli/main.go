package main

import (
	"fmt"

	"omg-cli/config"
	"omg-cli/omg"
	"omg-cli/ops_manager"

	"os"

	"omg-cli/pivnet"

	"log"

	"errors"
)

//TODO(jrjohnson): These constants should be detected, generated, or flags
const (
	username          = "foo"
	password          = "foobar"
	decryptionPhrase  = "foobar"
	skipSSLValidation = true
	terraformState    = "env.json"
)

type step func() error

func main() {
	logger := log.New(os.Stderr, "[ONG] ", 0)

	setup, err := NewApp(logger)
	if err != nil {
		logger.Fatal(err)
	}

	run([]step{
		func() error { return setup.SetupAuth(decryptionPhrase) },
		setup.SetupBosh,
		setup.ApplyChanges,
		//setup.UploadERT,
		//setup.UploadNozzle,
		//setup.UploadServiceBroker,
		//setup.ConfigureERT,
		//TODO(jrjohnson): ConfigureNozzle
		//TODO(jrjohnson): ConfigureServiceBroker
	}, logger)
}

func run(steps []step, logger *log.Logger) {
	for _, v := range steps {
		if err := v(); err != nil {
			logger.Fatal(err)
		}
	}
}

func LoadTerraformConfig() (*config.Config, error) {
	return config.FromTerraform(terraformState)
}

func NewApp(logger *log.Logger) (*omg.SetupService, error) {
	cfg, err := LoadTerraformConfig()
	if err != nil {
		return nil, err
	}

	omSdk, err := ops_manager.NewSdk(fmt.Sprintf("https://%s", cfg.OpsManagerIp), username, password, skipSSLValidation, *logger)
	if err != nil {
		return nil, err
	}

	pivnetAPIToken := os.Getenv("PIVNET_API_TOKEN")
	if pivnetAPIToken == "" {
		return nil, errors.New("expected environment variable PIVNET_API_TOKEN. Look for 'API TOKEN' at https://network.pivotal.io/users/dashboard/edit-profile")
	}

	pivnetSdk, err := pivnet.NewSdk(pivnetAPIToken, logger)
	if err != nil {
		return nil, err
	}

	return omg.NewSetupService(cfg, omSdk, pivnetSdk), nil
}
