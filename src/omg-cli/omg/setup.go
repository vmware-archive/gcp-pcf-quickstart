package omg

import (
	"omg-cli/config"
	"omg-cli/ops_manager"
	"omg-cli/pivnet"
	"os"

	"fmt"

	"net"

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
