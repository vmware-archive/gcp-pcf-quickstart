package omg

import (
	"omg-cli/config"
	"omg-cli/ops_manager"
	"omg-cli/pivnet"
	"os"

	"fmt"

	"net"

	"encoding/json"

	"github.com/pivotal-cf/om/commands"
)

const (
	metadataService = "169.254.169.254"
)

type pivnetDefinition struct {
	name      string
	versionId string
	fileId    string
	sha256    string
}

type productDefinition struct {
	name    string
	version string
}

type tileDefinition struct {
	pivnet  pivnetDefinition
	product productDefinition
}

var ertTile = tileDefinition{
	pivnetDefinition{
		"elastic-runtime",
		"5993",
		"24044",
		"a1d248287fff3328459dedb10921394949f818e7b89f017803ac7d23a6c27bf2",
	},
	productDefinition{
		"cf",
		"1.11.2",
	},
}

type SetupService struct {
	cfg    *config.Config
	om     *ops_manager.Sdk
	pivnet *pivnet.Sdk
}

func NewSetupService(cfg *config.Config, omSdk *ops_manager.Sdk, pivnetSdk *pivnet.Sdk) *SetupService {
	return &SetupService{cfg: cfg, om: omSdk, pivnet: pivnetSdk}
}

func (s *SetupService) SetupAuth(decryptionPhrase string) error {
	return s.om.SetupAuth(decryptionPhrase)
}

func (s *SetupService) buildNetwork(name, cidrRange, gateway string) commands.NetworkConfiguration {
	// Reserve .1-.20
	lowerIp, _, err := net.ParseCIDR(cidrRange)
	lowerIp = lowerIp.To4()
	if err != nil {
		panic(err)
	}
	upperIp := make(net.IP, len(lowerIp))
	copy(upperIp, lowerIp)
	upperIp[3] = 20

	return commands.NetworkConfiguration{
		Name: name,
		Subnets: []commands.Subnet{
			{
				IAASIdentifier:    fmt.Sprintf("%s/%s/%s", s.cfg.NetworkName, name, "us-central1"),
				CIDR:              cidrRange,
				Gateway:           gateway,
				ReservedIPRanges:  fmt.Sprintf("%s-%s", lowerIp.String(), upperIp.String()),
				AvailabilityZones: []string{"us-central1-b", "us-central1-c", "us-central1-f"},
				DNS:               metadataService,
			},
		},
	}
}

func (s *SetupService) SetupBosh() error {
	gcp := commands.GCPIaaSConfiguration{
		Project:              s.cfg.ProjectName,
		DefaultDeploymentTag: "omg-opsman",
		AuthJSON:             "",
	}

	director := commands.DirectorConfiguration{
		NTPServers: metadataService, // GCP metadata service
	}

	azs := commands.AvailabilityZonesConfiguration{
		AvailabilityZones: []commands.AvailabilityZone{
			{Name: "us-central1-b"},
			{Name: "us-central1-c"},
			{Name: "us-central1-f"},
		},
	}

	networks := commands.NetworksConfiguration{
		ICMP: false,
		Networks: []commands.NetworkConfiguration{
			s.buildNetwork(s.cfg.MgmtSubnetName, s.cfg.MgmtSubnetCIDR, s.cfg.MgmtSubnetGateway),
			s.buildNetwork(s.cfg.ServicesSubnetName, s.cfg.ServicesSubnetCIDR, s.cfg.ServicesSubnetGateway),
			s.buildNetwork(s.cfg.ErtSubnetName, s.cfg.ErtSubnetCIDR, s.cfg.ErtSubnetGateway),
		},
	}

	networkAssignment := commands.NetworkAssignment{
		UserProvidedNetworkName: s.cfg.MgmtSubnetName,
		UserProvidedAZName:      "us-central1-b",
	}

	if err := s.om.SetupBosh(gcp, director, azs, networks, networkAssignment); err != nil {
		return err
	}

	return nil
}

func (s *SetupService) ApplyChanges() error {
	return s.om.ApplyChanges()
}

func (s *SetupService) productInstalled(name, version string) (bool, error) {
	products, err := s.om.AvaliableProducts()
	if err != nil {
		return false, err
	}

	for _, p := range products {
		if p.Name == name && p.Version == version {
			return true, nil
		}
	}
	return false, nil
}

func (s *SetupService) ensureProductReady(tile tileDefinition) error {
	if i, err := s.productInstalled(tile.product.name, tile.product.version); i == true || err != nil {
		return err
	}

	file, err := s.pivnet.DownloadTile(tile.pivnet.name, tile.pivnet.versionId, tile.pivnet.fileId, tile.pivnet.sha256)
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	if err = s.om.UploadProduct(file.Name()); err != nil {
		return err
	}

	return s.om.StageProduct(tile.product.name, tile.product.version)
}

func (s *SetupService) UploadERT() error {
	return s.ensureProductReady(ertTile)
}

// TODO(jrjohnson): Move to it's own ert (sub?)package
type ErtAvalibilityZone struct {
	Name string `json:"name"`
}

type ErtNetworkName struct {
	Name string `json:"name"`
}

type ErtNetwork struct {
	SingletonAvalibilityZone ErtAvalibilityZone   `json:"singleton_availability_zone"`
	OtherAvailabilityZones   []ErtAvalibilityZone `json:"other_availability_zones"`
	Network                  ErtNetworkName       `json:"network"`
}

type ErtProperties struct {
	// Domains
	AppsDomain ErtValue `json:".cloud_controller.apps_domain"`
	SysDomain  ErtValue `json:".cloud_controller.system_domain"`
	// Networking
	NetworkingPointOfEntry    ErtValue `json:".properties.networking_point_of_entry"`
	TcpRouting                ErtValue `json:".properties.tcp_routing"`
	TcpRoutingReservablePorts ErtValue `json:".properties.tcp_routing.enable.reservable_ports"`
	// Application Security Groups
	SecurityAcknowledgement ErtValue `json:".properties.security_acknowledgement"`
	// UAA
	ServiceProviderCredentials ErtRsaCertCredentaial `json:".uaa.service_provider_key_credentials"`
	// MySQL
	MySqlMonitorRecipientEmail ErtValue `json:".mysql_monitor.recipient_email"`
}

type ErtValue struct {
	Value string `json:"value"`
}

type ErtCert struct {
	Cert       string `json:"cert_pem"`
	PrivateKey string `json:"private_key_pem"`
}

type ErtRsaCertCredentaial struct {
	Value ErtCert `json:"value"`
}

type ErtResources struct {
	TcpRouter  ErtResource `json:"tcp_router"`
	Router     ErtResource `json:"router"`
	DiegoBrain ErtResource `json:"diego_brain"`
}

type ErtResource struct {
	RouterNames       []string `json:"elb_names,omitempty"`
	Instances         int      `json:"instances,omitempty"`
	InternetConnected bool     `json:"internet_connected"`
}

func (s *SetupService) ConfigureERT() error {
	ertNetwork := ErtNetwork{
		ErtAvalibilityZone{"us-central1-b"},
		[]ErtAvalibilityZone{{"us-central1-b"}, {"us-central1-c"}, {"us-central1-f"}},
		ErtNetworkName{s.cfg.ErtSubnetName},
	}

	ertNetworkBytes, err := json.Marshal(&ertNetwork)
	if err != nil {
		return err
	}

	ertProperties := ErtProperties{
		AppsDomain:                 ErtValue{fmt.Sprintf("apps.%s", s.cfg.RootDomain)},
		SysDomain:                  ErtValue{fmt.Sprintf("sys.%s", s.cfg.RootDomain)},
		NetworkingPointOfEntry:     ErtValue{"external_non_ssl"},
		TcpRouting:                 ErtValue{"enable"},
		TcpRoutingReservablePorts:  ErtValue{s.cfg.TcpPortRange},
		SecurityAcknowledgement:    ErtValue{"X"},
		ServiceProviderCredentials: ErtRsaCertCredentaial{ErtCert{s.cfg.SslCertificate, s.cfg.SslPrivateKey}},
		MySqlMonitorRecipientEmail: ErtValue{"admin@example.org"},
	}

	ertPropertiesBytes, err := json.Marshal(&ertProperties)
	if err != nil {
		return err
	}

	ertResoruces := ErtResources{
		TcpRouter: ErtResource{
			RouterNames:       []string{fmt.Sprintf("tcp:%s", s.cfg.TcpTargetPoolName)},
			InternetConnected: false,
		},
		Router: ErtResource{
			RouterNames:       []string{fmt.Sprintf("http:%s", s.cfg.HttpBackendServiceName)},
			InternetConnected: false,
		},
		DiegoBrain: ErtResource{
			RouterNames:       []string{fmt.Sprintf("tcp:%s", s.cfg.SshTargetPoolName)},
			InternetConnected: false,
		},
	}
	ertResorucesBytes, err := json.Marshal(&ertResoruces)
	if err != nil {
		return err
	}

	return s.om.ConfigureProduct(ertTile.product.name, string(ertNetworkBytes), string(ertPropertiesBytes), string(ertResorucesBytes))
}
