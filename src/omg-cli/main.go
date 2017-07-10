package main

import (
	"context"
	"fmt"
	"omg-cli/config"
	"omg-cli/ops_manager"

	"golang.org/x/oauth2/google"
	runtimeconfig "google.golang.org/api/runtimeconfig/v1beta1"
)

const projectName = "google.com:graphite-test-bosh-cpi-cert"

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

	fmt.Printf("%#v", cfg)
	om := ops_manager.New(cfg)
	err = om.SetupAuth()
	if err != nil {
		fmt.Printf("err: %v", err)
	}

	err = om.SetupBosh()
	if err != nil {
		fmt.Printf("err: %v", err)
	}
}
