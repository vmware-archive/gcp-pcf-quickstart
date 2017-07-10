package omg

import (
	"omg-cli/config"
	"omg-cli/ops_manager"

	"fmt"

	"github.com/pivotal-cf/om/commands"
)

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

func (s *SetupService) fullSubnetName(name string) string {
	return fmt.Sprintf("%s/%s/%s", s.cfg.NetworkName, name, "us-central1")
}

func (s *SetupService) SetupBosh() error {
	gcp := commands.GCPIaaSConfiguration{
		Project:              s.cfg.ProjectName,
		DefaultDeploymentTag: "omg-opsman",
		AuthJSON:             "",
	}

	director := commands.DirectorConfiguration{
		NTPServers: "169.254.169.254", // GCP metadata service
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
			{
				Name: s.cfg.MgmtSubnetName,
				Subnets: []commands.Subnet{
					{
						IAASIdentifier:    s.fullSubnetName(s.cfg.MgmtSubnetName),
						CIDR:              s.cfg.MgmtSubnetCIDR,
						Gateway:           s.cfg.MgmtSubnetGateway,
						ReservedIPRanges:  "10.0.0.0-10.0.0.20", // TODO(jrjohnson): Not true
						AvailabilityZones: []string{"us-central1-b", "us-cetnral1-c", "us-central1-f"},
						DNS:               "169.254.169.254",
					},
				},
			},
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
