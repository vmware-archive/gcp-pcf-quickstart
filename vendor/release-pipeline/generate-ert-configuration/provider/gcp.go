package provider

import (
	"encoding/json"
	"errors"
	"fmt"
)

type GCPMetadata struct {
	AZs                       []string `json:"azs"`
	Network                   string   `json:"network_name"`
	ERTSubnet                 string   `json:"ert_subnet_name"`
	ERTCIDR                   string   `json:"ert_cidr"`
	Project                   string   `json:"project"`
	Region                    string   `json:"region"`
	ServiceKey                string   `json:"service_account_key"`
	HTTPLBBackendName         string   `json:"http_lb_backend_name"`
	WSPool                    string   `json:"ws_router_pool"`
	SSHPool                   string   `json:"ssh_router_pool"`
	TCPPool                   string   `json:"tcp_router_pool"`
	SysDomain                 string   `json:"sys_domain"`
	AppsDomain                string   `json:"apps_domain"`
	PackagesBucket            string   `json:"packages_bucket"`
	DropletsBucket            string   `json:"droplets_bucket"`
	ResourcesBucket           string   `json:"resources_bucket"`
	BuildpacksBucket          string   `json:"buildpacks_bucket"`
	StorageInteropAccessKey   string   `json:"storage_interop_access_key"`
	StorageInteropSecretKey   string   `json:"storage_interop_secret_key"`
	SQLHost                   string   `json:"sql_db_ip"`
	ERTSQLUsername            string   `json:"ert_sql_username"`
	ERTSQLPassword            string   `json:"ert_sql_password"`
	SMTPUsername              string   `json:"smtp_username"`
	SMTPPassword              string   `json:"smtp_password"`
	EnableContainerNetworking bool     `json:"enable_container_networking"`
}

func ParseGCPMetadata(metadata string) (GCPMetadata, error) {
	var gcpMetadata GCPMetadata
	err := json.Unmarshal([]byte(metadata), &gcpMetadata)
	if err != nil {
		return GCPMetadata{}, err
	}

	if len(gcpMetadata.AZs) == 0 {
		return GCPMetadata{}, errors.New("error: expected at least one AZ")
	}

	return gcpMetadata, nil
}

func PrepareGCPForOm(gcpMetadata GCPMetadata, settings Settings) (FullConfiguration, error) {
	fullConfig := FullConfiguration{
		ERTConfig: ERTConfiguration{
			OpsManagerDomain: settings.OpsManagerDomain,
			ProductProperties: ERTProperties{
				SysDomain:           &Property{Value: gcpMetadata.SysDomain},
				AppsDomain:          &Property{Value: gcpMetadata.AppsDomain},
				HAProxySkipSSL:      &Property{Value: true},
				NetworkingPOE:       &Property{Value: "external_ssl"},
				MySQLRecipientEmail: &Property{Value: "cf-release-engineering@pivotal.io"},
				ExternalSSLCert: &Property{Value: map[string]string{
					"cert_pem":        settings.SSLCertificate,
					"private_key_pem": settings.SSLPrivateKey,
				}},
				SecurityCheckBox:    &Property{Value: "X"},
				SystemBlobstore:     &Property{Value: "external_gcs"},
				GCSPackagesBucket:   &Property{Value: gcpMetadata.PackagesBucket},
				GCSDropletsBucket:   &Property{Value: gcpMetadata.DropletsBucket},
				GCSResourcesBucket:  &Property{Value: gcpMetadata.ResourcesBucket},
				GCSBuildpacksBucket: &Property{Value: gcpMetadata.BuildpacksBucket},
				GCSAccessKey:        &Property{Value: gcpMetadata.StorageInteropAccessKey},
				GCSSecretKey: &Property{Value: map[string]string{
					"secret": gcpMetadata.StorageInteropSecretKey,
				}},
				TCPRouting:         &Property{Value: "enable"},
				TCPReservablePorts: &Property{Value: "1024-1123"},
				SMTPFrom:           &Property{Value: "identitystaging@mailinator.com"},
				SMTPAddress:        &Property{Value: "smtp.sendgrid.net"},
				SMTPPort:           &Property{Value: "2525"},
				SMTPCredentials: &Property{Value: map[string]string{
					"identity": gcpMetadata.SMTPUsername,
					"password": gcpMetadata.SMTPPassword,
				}},
				SMTPEnableStartTLSAuto: &Property{Value: true},
				SMTPAuthMechanism:      &Property{Value: "plain"},
			},
			ProductNetwork: ERTNetworks{
				SingletonAZ: AZ{Name: gcpMetadata.AZs[0]},
				OtherAZs:    generateAZNameMap(gcpMetadata.AZs),
				Network:     Network{Name: gcpMetadata.ERTSubnet},
			},
			ProductResources: ERTResources{
				TCPRouter: &JobResourceConfig{
					ELBNames: []string{fmt.Sprintf("tcp:%s", gcpMetadata.TCPPool)},
				},
				Router: &JobResourceConfig{
					ELBNames: []string{fmt.Sprintf("http:%s", gcpMetadata.HTTPLBBackendName), fmt.Sprintf("tcp:%s", gcpMetadata.WSPool)},
				},
				DiegoBrain: &JobResourceConfig{
					ELBNames: []string{fmt.Sprintf("tcp:%s", gcpMetadata.SSHPool)},
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

	// fully understand this is duplicated, this whole config generation program/
	// is going to change.
	if gcpMetadata.ERTSQLUsername != "" {
		fullConfig.ERTConfig.ProductProperties.SystemDB = &Property{Value: "external"}
		fullConfig.ERTConfig.ProductProperties.ExtDBPort = &Property{Value: 3306}
		fullConfig.ERTConfig.ProductProperties.ExtDBHost = &Property{Value: gcpMetadata.SQLHost}
		fullConfig.ERTConfig.ProductProperties.AppUsageDBUsername = &Property{Value: gcpMetadata.ERTSQLUsername}
		fullConfig.ERTConfig.ProductProperties.AppUsageDBPassword = &Property{Value: map[string]string{
			"secret": gcpMetadata.ERTSQLPassword,
		}}
		fullConfig.ERTConfig.ProductProperties.AutoscaleDBUsername = &Property{Value: gcpMetadata.ERTSQLUsername}
		fullConfig.ERTConfig.ProductProperties.AutoscaleDBPassword = &Property{Value: map[string]string{
			"secret": gcpMetadata.ERTSQLPassword,
		}}
		fullConfig.ERTConfig.ProductProperties.AccountDBUsername = &Property{Value: gcpMetadata.ERTSQLUsername}
		fullConfig.ERTConfig.ProductProperties.AccountDBPassword = &Property{Value: map[string]string{
			"secret": gcpMetadata.ERTSQLPassword,
		}}
		fullConfig.ERTConfig.ProductProperties.DiegoDBUsername = &Property{Value: gcpMetadata.ERTSQLUsername}
		fullConfig.ERTConfig.ProductProperties.DiegoDBPassword = &Property{Value: map[string]string{
			"secret": gcpMetadata.ERTSQLPassword,
		}}
		fullConfig.ERTConfig.ProductProperties.CCDBUsername = &Property{Value: gcpMetadata.ERTSQLUsername}
		fullConfig.ERTConfig.ProductProperties.CCDBPassword = &Property{Value: map[string]string{
			"secret": gcpMetadata.ERTSQLPassword,
		}}
		fullConfig.ERTConfig.ProductProperties.NetworkPolicyDBUsername = &Property{Value: gcpMetadata.ERTSQLUsername}
		fullConfig.ERTConfig.ProductProperties.NetworkPolicyDBPassword = &Property{Value: map[string]string{
			"secret": gcpMetadata.ERTSQLPassword,
		}}
		fullConfig.ERTConfig.ProductProperties.NFSVolumeDBUsername = &Property{Value: gcpMetadata.ERTSQLUsername}
		fullConfig.ERTConfig.ProductProperties.NFSVolumeDBPassword = &Property{Value: map[string]string{
			"secret": gcpMetadata.ERTSQLPassword,
		}}
		fullConfig.ERTConfig.ProductProperties.NotificationsDBUsername = &Property{Value: gcpMetadata.ERTSQLUsername}
		fullConfig.ERTConfig.ProductProperties.NotificationsDBPassword = &Property{Value: map[string]string{
			"secret": gcpMetadata.ERTSQLPassword,
		}}
		fullConfig.ERTConfig.ProductProperties.RoutingDBUsername = &Property{Value: gcpMetadata.ERTSQLUsername}
		fullConfig.ERTConfig.ProductProperties.RoutingDBPassword = &Property{Value: map[string]string{
			"secret": gcpMetadata.ERTSQLPassword,
		}}
		fullConfig.ERTConfig.ProductProperties.UAADBUsername = &Property{Value: gcpMetadata.ERTSQLUsername}
		fullConfig.ERTConfig.ProductProperties.UAADBPassword = &Property{Value: map[string]string{
			"secret": gcpMetadata.ERTSQLPassword,
		}}
	}

	if gcpMetadata.EnableContainerNetworking {
		fullConfig.ERTConfig.ProductProperties.ContainerNetworking = &Property{Value: "enable"}
	}

	return fullConfig, nil
}

func generateAZNameMap(azs []string) []AZ {
	var output []AZ
	for _, az := range azs {
		output = append(output, AZ{Name: az})
	}
	return output
}
