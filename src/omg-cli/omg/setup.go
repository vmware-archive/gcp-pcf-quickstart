package omg

import (
	"omg-cli/config"
	"omg-cli/ops_manager"

	"fmt"

	"net"

	"github.com/pivotal-cf/om/commands"
)

const metadataService = "169.254.169.254"

type SetupService struct {
	cfg *config.Config
	sdk *ops_manager.Sdk
}

func NewSetupService(cfg *config.Config, sdk *ops_manager.Sdk) *SetupService {
	return &SetupService{cfg: cfg, sdk: sdk}
}

func (s *SetupService) SetupAuth(decryptionPhrase string) error {
	return s.sdk.SetupAuth(decryptionPhrase)
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

	if err := s.sdk.SetupBosh(gcp, director, azs, networks, networkAssignment); err != nil {
		return err
	}

	return nil
}

func (s *SetupService) ApplyChanges() error {
	return s.sdk.ApplyChanges()
}
