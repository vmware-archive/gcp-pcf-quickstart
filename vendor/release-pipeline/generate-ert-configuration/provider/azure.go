package provider

import (
	"encoding/json"
	"fmt"
)

type AzureMetadata struct {
	Network                      string `json:"network_name"`
	ERTSubnet                    string `json:"ert_subnet_name"`
	SysDomain                    string `json:"sys_domain"`
	AppsDomain                   string `json:"apps_domain"`
	WebLBName                    string `json:"web_lb_name"`
	MySQLDNS                     string `json:"mysql_dns"`
	MySQLLBName                  string `json:"mysql_lb_name"`
	TCPLBName                    string `json:"tcp_lb_name"`
	SMTPUsername                 string `json:"smtp_username"`
	SMTPPassword                 string `json:"smtp_password"`
	CFStorageAccountName         string `json:"cf_storage_account_name"`
	CFStorageAccountAccessKey    string `json:"cf_storage_account_access_key"`
	CFBuildpacksStorageContainer string `json:"cf_buildpacks_storage_container"`
	CFDropletsStorageContainer   string `json:"cf_droplets_storage_container"`
	CFPackagesStorageContainer   string `json:"cf_packages_storage_container"`
	CFResourcesStorageContainer  string `json:"cf_resources_storage_container"`
}

func ParseAzureMetadata(metadata string) (AzureMetadata, error) {
	var azureMetadata AzureMetadata
	err := json.Unmarshal([]byte(metadata), &azureMetadata)
	if err != nil {
		return AzureMetadata{}, err
	}

	return azureMetadata, nil
}

func PrepareAzureForOm(azureMetadata AzureMetadata, settings Settings) (FullConfiguration, error) {
	fullConfig := FullConfiguration{
		ERTConfig: ERTConfiguration{
			OpsManagerDomain: settings.OpsManagerDomain,
			ProductProperties: ERTProperties{
				SysDomain:           &Property{Value: azureMetadata.SysDomain},
				AppsDomain:          &Property{Value: azureMetadata.AppsDomain},
				HAProxySkipSSL:      &Property{Value: true},
				NetworkingPOE:       &Property{Value: "haproxy"},
				MySQLRecipientEmail: &Property{Value: "cf-release-engineering@pivotal.io"},
				HAProxySSLCert: &Property{Value: map[string]string{
					"cert_pem":        settings.SSLCertificate,
					"private_key_pem": settings.SSLPrivateKey,
				}},
				SecurityCheckBox:   &Property{Value: "X"},
				TCPRouting:         &Property{Value: "enable"},
				TCPReservablePorts: &Property{Value: "1024-1173"},
				SMTPFrom:           &Property{Value: "identitystaging@mailinator.com"},
				SMTPAddress:        &Property{Value: "smtp.sendgrid.net"},
				SMTPPort:           &Property{Value: "2525"},
				SMTPCredentials: &Property{Value: map[string]string{
					"identity": azureMetadata.SMTPUsername,
					"password": azureMetadata.SMTPPassword,
				}},
				SMTPEnableStartTLSAuto: &Property{Value: true},
				SMTPAuthMechanism:      &Property{Value: "plain"},
			},
			ProductNetwork: ERTNetworks{
				SingletonAZ: AZ{Name: "null"},
				OtherAZs:    []AZ{AZ{Name: "null"}},
				Network:     Network{Name: azureMetadata.ERTSubnet},
			},
			ProductResources: ERTResources{
				TCPRouter: &JobResourceConfig{
					ELBNames: []string{azureMetadata.TCPLBName},
				},
				HAProxy: &JobResourceConfig{
					ELBNames: []string{azureMetadata.WebLBName},
				},
				MySQLProxy: &JobResourceConfig{
					ELBNames: []string{azureMetadata.MySQLLBName},
				},
			},
		},
	}

	if settings.SAMLEnabled {
		fullConfig.ERTConfig.ProductProperties.SAMLCert = &Property{
			Value: map[string]string{
				"cert_pem":        settings.SSLCertificate,
				"private_key_pem": settings.SSLPrivateKey,
			},
		}
	}

	if azureMetadata.CFStorageAccountName != "" {
		fullConfig.ERTConfig.ProductProperties.SystemBlobstore = &Property{Value: "external_azure"}
		fullConfig.ERTConfig.ProductProperties.CFStorageAccountName = &Property{Value: azureMetadata.CFStorageAccountName}
		fullConfig.ERTConfig.ProductProperties.CFStorageAccountAccessKey = &Property{Value: map[string]string{
			"secret": azureMetadata.CFStorageAccountAccessKey,
		}}
		fullConfig.ERTConfig.ProductProperties.CFBuildpacksStorageContainer = &Property{Value: azureMetadata.CFBuildpacksStorageContainer}
		fullConfig.ERTConfig.ProductProperties.CFDropletsStorageContainer = &Property{Value: azureMetadata.CFDropletsStorageContainer}
		fullConfig.ERTConfig.ProductProperties.CFPackagesStorageContainer = &Property{Value: azureMetadata.CFPackagesStorageContainer}
		fullConfig.ERTConfig.ProductProperties.CFResourcesStorageContainer = &Property{Value: azureMetadata.CFResourcesStorageContainer}
	}

	if settings.Resources != "" {
		err := json.Unmarshal([]byte(settings.Resources), &fullConfig.ERTConfig.ProductResources)
		if err != nil {
			return fullConfig, fmt.Errorf("could not parse resource configuration: %s", err)
		}
	}

	return fullConfig, nil
}
