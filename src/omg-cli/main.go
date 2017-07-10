package main

import (
	"context"
	"fmt"
	"omg-cli/config"
	"omg-cli/omg"
	"omg-cli/ops_manager"

	"golang.org/x/oauth2/google"
	runtimeconfig "google.golang.org/api/runtimeconfig/v1beta1"
)

const (
	projectName      = "google.com:graphite-test-bosh-cpi-cert"
	username         = "foo"
	password         = "foobar"
	decryptionPhrase = "foobar"
)

func main() {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, runtimeconfig.CloudruntimeconfigScope, runtimeconfig.CloudPlatformScope)
	if err != nil {
		panic(err)
	}

	cfg, err := config.FromEnvironment(ctx, client, projectName)
	if err != nil {
		panic(err)
	}

	sdk, err := ops_manager.NewSdk(fmt.Sprintf("https://%s", cfg.OpsManagerIp), username, password, config.SkipSSLValidation)
	if err != nil {
		panic(err)
	}

	setup := omg.NewSetupService(cfg, sdk)
	err = setup.SetupAuth(decryptionPhrase)
	if err != nil {
		fmt.Printf("err: %v", err)
	}

	err = setup.SetupBosh()
	if err != nil {
		fmt.Printf("err: %v", err)
	}
}
