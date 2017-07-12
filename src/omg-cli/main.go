package main

import (
	"context"
	"fmt"

	"google.golang.org/api/runtimeconfig/v1beta1"

	"omg-cli/config"
	"omg-cli/omg"
	"omg-cli/ops_manager"

	"os"

	"omg-cli/pivnet"

	"log"

	"errors"

	"golang.org/x/oauth2/google"
)

//TODO(jrjohnson): These constants should be detected, generated, or flags
const (
	projectName       = "google.com:graphite-test-bosh-cpi-cert"
	username          = "foo"
	password          = "foobar"
	decryptionPhrase  = "foobar"
	skipSSLValidation = true
)

func main() {
	logger := log.New(os.Stderr, "[ONG] ", 0)

	setup, err := NewApp(logger)
	if err != nil {
		logger.Fatal(err)
	}

	steps := []struct {
		fn func() error
	}{
		{func() error { return setup.SetupAuth(decryptionPhrase) }},
		{setup.SetupBosh},
		{setup.ApplyChanges},
		{setup.UploadERT},
		{setup.ConfigureERT},
	}

	for _, v := range steps {
		if err := v.fn(); err != nil {
			logger.Fatal(err)
		}
	}
}

func NewApp(logger *log.Logger) (*omg.SetupService, error) {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, runtimeconfig.CloudruntimeconfigScope, runtimeconfig.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	cfg, err := config.FromEnvironment(ctx, client, projectName)
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
