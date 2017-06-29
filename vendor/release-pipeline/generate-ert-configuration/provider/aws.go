package provider

import (
	"encoding/json"
	"errors"
	"fmt"
)

type AWSMetadata struct {
	AZs              []string `json:"azs"`
	Network          string   `json:"network_name"`
	SysDomain        string   `json:"sys_domain"`
	AppsDomain       string   `json:"apps_domain"`
	PackagesBucket   string   `json:"ert_packages_bucket"`
	DropletsBucket   string   `json:"ert_droplets_bucket"`
	ResourcesBucket  string   `json:"ert_resources_bucket"`
	BuildpacksBucket string   `json:"ert_buildpacks_bucket"`
	AccessKey        string   `json:"iam_user_access_key"`
	SecretKey        string   `json:"iam_user_secret_access_key"`
	Region           string   `json:"region"`
	SMTPUsername     string   `json:"smtp_username"`
	SMTPPassword     string   `json:"smtp_password"`
	TCPELB           string   `json:"tcp_elb_name"`
	WebELB           string   `json:"web_elb_name"`
	SSHELB           string   `json:"ssh_elb_name"`
}

func ParseAWSMetadata(metadata string) (AWSMetadata, error) {
	var awsMetadata AWSMetadata
	err := json.Unmarshal([]byte(metadata), &awsMetadata)
	if err != nil {
		return AWSMetadata{}, err
	}

	if len(awsMetadata.AZs) == 0 {
		return AWSMetadata{}, errors.New("error: expected at least one AZ")
	}

	return awsMetadata, nil
}

func PrepareAWSForOm(awsMetadata AWSMetadata, settings Settings) (FullConfiguration, error) {
	fullConfig := FullConfiguration{
		ERTConfig: ERTConfiguration{
			OpsManagerDomain: settings.OpsManagerDomain,
			ProductProperties: ERTProperties{
				SysDomain:      &Property{Value: awsMetadata.SysDomain},
				AppsDomain:     &Property{Value: awsMetadata.AppsDomain},
				HAProxySkipSSL: &Property{Value: true},
				NetworkingPOE:  &Property{Value: "external_ssl"},
				ExternalSSLCert: &Property{Value: map[string]string{
					"cert_pem":        settings.SSLCertificate,
					"private_key_pem": settings.SSLPrivateKey,
				}},
				MySQLRecipientEmail:      &Property{Value: "cf-release-engineering@pivotal.io"},
				SecurityCheckBox:         &Property{Value: "X"},
				SystemBlobstore:          &Property{Value: "external"},
				ExternalPackagesBucket:   &Property{Value: awsMetadata.PackagesBucket},
				ExternalDropletsBucket:   &Property{Value: awsMetadata.DropletsBucket},
				ExternalResourcesBucket:  &Property{Value: awsMetadata.ResourcesBucket},
				ExternalBuildpacksBucket: &Property{Value: awsMetadata.BuildpacksBucket},
				StorageIAMEndpoint:       &Property{Value: fmt.Sprintf("https://s3-%s.amazonaws.com", awsMetadata.Region)},
				StorageIAMAccessKey:      &Property{Value: awsMetadata.AccessKey},
				StorageIAMSSecretAccessKey: &Property{Value: map[string]string{
					"secret": awsMetadata.SecretKey,
				}},
				StorageIAMRegion:   &Property{Value: awsMetadata.Region},
				TCPRouting:         &Property{Value: "enable"},
				TCPReservablePorts: &Property{Value: "1024-1123"},
				SMTPFrom:           &Property{Value: "identitystaging@mailinator.com"},
				SMTPAddress:        &Property{Value: "smtp.sendgrid.net"},
				SMTPPort:           &Property{Value: "2525"},
				SMTPCredentials: &Property{Value: map[string]string{
					"identity": awsMetadata.SMTPUsername,
					"password": awsMetadata.SMTPPassword,
				}},
				SMTPEnableStartTLSAuto: &Property{Value: true},
				SMTPAuthMechanism:      &Property{Value: "plain"},
				LoggerEndpointPort:     &Property{Value: "4443"},
			},
			ProductNetwork: ERTNetworks{
				SingletonAZ: AZ{Name: awsMetadata.AZs[0]},
				OtherAZs:    generateAZNameMap(awsMetadata.AZs),
				Network:     Network{Name: awsMetadata.Network},
			},
			ProductResources: ERTResources{
				TCPRouter: &JobResourceConfig{
					ELBNames: []string{awsMetadata.TCPELB},
				},
				Router: &JobResourceConfig{
					ELBNames: []string{awsMetadata.WebELB},
				},
				DiegoBrain: &JobResourceConfig{
					ELBNames: []string{awsMetadata.SSHELB},
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

	if settings.Resources != "" {
		err := json.Unmarshal([]byte(settings.Resources), &fullConfig.ERTConfig.ProductResources)
		if err != nil {
			return fullConfig, fmt.Errorf("could not parse resource configuration: %s", err)
		}
	}

	return fullConfig, nil
}
