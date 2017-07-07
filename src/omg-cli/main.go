package main

import (
	"context"
	"fmt"
	"omg-cli/config"

	"golang.org/x/oauth2/google"
	runtimeconfig "google.golang.org/api/runtimeconfig/v1beta1"
)

const configName = "projects/google.com:graphite-test-bosh-cpi-cert/configs/omgConfig"

func main() {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, runtimeconfig.CloudruntimeconfigScope, runtimeconfig.CloudPlatformScope)
	if err != nil {
		panic(err)
	}

	envCfg, err := config.FromEnvironment(ctx, client, configName)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", envCfg)
}
