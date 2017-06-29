package provider

import (
	"encoding/json"
	"fmt"
)

type VSphereMetadata struct {
	AZs          []string `json:"azs"`
	Network      string   `json:"network_name"`
	SysDomain    string   `json:"sys_domain"`
	AppsDomain   string   `json:"apps_domain"`
	SMTPUsername string   `json:"smtp_username"`
	SMTPPassword string   `json:"smtp_password"`
}

func ParseVSphereMetadata(metadata string) (VSphereMetadata, error) {
	var vsphereMetadata VSphereMetadata
	err := json.Unmarshal([]byte(metadata), &vsphereMetadata)
	if err != nil {
		return VSphereMetadata{}, err
	}

	return vsphereMetadata, nil
}

func PrepareVSphereForOm(vsphereMetadata VSphereMetadata, settings Settings) (FullConfiguration, error) {
	fullConfig := FullConfiguration{
		ERTConfig: ERTConfiguration{
			OpsManagerDomain: settings.OpsManagerDomain,
			ProductProperties: ERTProperties{
				SysDomain:      &Property{Value: vsphereMetadata.SysDomain},
				AppsDomain:     &Property{Value: vsphereMetadata.AppsDomain},
				HAProxySkipSSL: &Property{Value: true},
				NetworkingPOE:  &Property{Value: "haproxy"},
				HAProxySSLCert: &Property{Value: map[string]string{
					"cert_pem":        settings.SSLCertificate,
					"private_key_pem": settings.SSLPrivateKey,
				}},
				MySQLRecipientEmail: &Property{Value: "cf-release-engineering@pivotal.io"},
				SecurityCheckBox:    &Property{Value: "X"},
				SystemBlobstore:     &Property{Value: "internal"},
				TCPRouting:          &Property{Value: "disable"},
				SMTPFrom:            &Property{Value: "identitystaging@mailinator.com"},
				SMTPAddress:         &Property{Value: "smtp.sendgrid.net"},
				SMTPPort:            &Property{Value: "2525"},
				SMTPCredentials: &Property{Value: map[string]string{
					"identity": vsphereMetadata.SMTPUsername,
					"password": vsphereMetadata.SMTPPassword,
				}},
				SMTPEnableStartTLSAuto: &Property{Value: true},
				SMTPAuthMechanism:      &Property{Value: "plain"},
			},
			ProductNetwork: ERTNetworks{
				SingletonAZ: AZ{Name: vsphereMetadata.AZs[0]},
				OtherAZs:    generateAZNameMap(vsphereMetadata.AZs),
				Network:     Network{Name: vsphereMetadata.Network},
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
