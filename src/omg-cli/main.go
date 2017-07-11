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

	"golang.org/x/oauth2/google"
)

const (
	projectName       = "google.com:graphite-test-bosh-cpi-cert"
	username          = "foo"
	password          = "foobar"
	decryptionPhrase  = "foobar"
	skipSSLValidation = true
)

var (
	pivnetAPIToken string
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

	omSdk, err := ops_manager.NewSdk(fmt.Sprintf("https://%s", cfg.OpsManagerIp), username, password, skipSSLValidation)
	if err != nil {
		panic(err)
	}

	pivnetAPIToken = os.Getenv("PIVNET_API_TOKEN")
	if pivnetAPIToken == "" {
		panic("expected environment variable PIVNET_API_TOKEN. Look for 'API TOKEN' at https://network.pivotal.io/users/dashboard/edit-profile")
	}

	stdout := log.New(os.Stdout, "", 0)
	pivnetSdk, err := pivnet.NewSdk(pivnetAPIToken, stdout)
	if err != nil {
		panic(err)
	}

	setup := omg.NewSetupService(cfg, omSdk, pivnetSdk)
	err = setup.SetupAuth(decryptionPhrase)
	if err != nil {
		fmt.Printf("err: %v", err)
	}

	err = setup.SetupBosh()
	if err != nil {
		fmt.Printf("err: %v", err)
	}

	err = setup.ApplyChanges()
	if err != nil {
		fmt.Printf("err: %v", err)
	}

	err = setup.UploadERT()
	if err != nil {
		fmt.Printf("err: %v", err)
	}
}
