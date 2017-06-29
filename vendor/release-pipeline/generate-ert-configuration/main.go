package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/pivotal-cf/pcf-releng-ci/tasks/om/generate-ert-configuration/provider"
)

var (
	opsManagerDomain      string
	givenProvider         string
	providerConfiguration string
	sslCert               string
	sslPrivateKey         string
	resources             string
	samlEnabled           bool
)

func main() {
	flag.StringVar(&opsManagerDomain, "ops-manager-domain", "", "The domain that Ops Manager has been deployed to")
	flag.StringVar(&givenProvider, "provider", "", "name of IaaS provider (e.g. gcp, azure)")
	flag.StringVar(&providerConfiguration, "provider-configuration", "", "JSON metadata for your environment (as extracted from a run of terraform)")
	flag.StringVar(&sslCert, "ssl-cert", "", "External SSL cert for CF networking point-of-entry")
	flag.StringVar(&sslPrivateKey, "ssl-private-key", "", "SSL private key for CF networking point-of-entry")
	flag.StringVar(&resources, "resources", "", "JSON containing resource configuration for jobs")
	flag.BoolVar(&samlEnabled, "enable-saml-cert", false, "when true, set SAML cert")
	flag.Parse()

	settings := provider.Settings{
		SSLCertificate:   sslCert,
		SSLPrivateKey:    sslPrivateKey,
		Resources:        resources,
		SAMLEnabled:      samlEnabled,
		OpsManagerDomain: opsManagerDomain,
	}

	switch givenProvider {
	case "gcp":
		metadata, err := provider.ParseGCPMetadata(providerConfiguration)
		if err != nil {
			log.Fatalf("could not parse provider configuration: %s", err)
		}

		fullConfig, err := provider.PrepareGCPForOm(metadata, settings)
		if err != nil {
			log.Fatal(err)
		}

		configurationJSON, err := json.MarshalIndent(fullConfig, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(string(configurationJSON))

	case "azure":
		metadata, err := provider.ParseAzureMetadata(providerConfiguration)
		if err != nil {
			log.Fatalf("could not parse provider configuration: %s", err)
		}

		fullConfig, err := provider.PrepareAzureForOm(metadata, settings)
		if err != nil {
			log.Fatal(err)
		}

		configurationJSON, err := json.MarshalIndent(fullConfig, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(string(configurationJSON))

	case "aws":
		metadata, err := provider.ParseAWSMetadata(providerConfiguration)
		if err != nil {
			log.Fatalf("could not parse provider configuration: %s", err)
		}

		fullConfig, err := provider.PrepareAWSForOm(metadata, settings)
		if err != nil {
			log.Fatal(err)
		}

		configurationJSON, err := json.MarshalIndent(fullConfig, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(string(configurationJSON))

	case "vsphere":
		metadata, err := provider.ParseVSphereMetadata(providerConfiguration)

		if err != nil {
			log.Fatalf("could not parse provider configuration: %s", err)
		}

		fullConfig, err := provider.PrepareVSphereForOm(metadata, settings)
		if err != nil {
			log.Fatal(err)
		}

		configurationJSON, err := json.MarshalIndent(fullConfig, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(string(configurationJSON))
	default:
		log.Fatalf("unsupported provider")
	}
}
