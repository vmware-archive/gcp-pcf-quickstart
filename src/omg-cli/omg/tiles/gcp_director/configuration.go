package gcp_director

import (
	"fmt"
	"net"
	"omg-cli/config"

	"omg-cli/ops_manager"

	"github.com/pivotal-cf/om/commands"
)

const (
	metadataService = "169.254.169.254"
)

func (Tile) Configure(cfg *config.Config, om *ops_manager.Sdk) error {
	networks, networkAssignment := networkCfg(cfg)

	return om.SetupBosh(gcp(cfg), director(), avalibilityZones(cfg), networks, networkAssignment, resources())
}

func buildNetwork(cfg *config.Config, name, cidrRange, gateway string) commands.NetworkConfiguration {
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
				IAASIdentifier:    fmt.Sprintf("%s/%s/%s", cfg.NetworkName, name, cfg.Region),
				CIDR:              cidrRange,
				Gateway:           gateway,
				ReservedIPRanges:  fmt.Sprintf("%s-%s", lowerIp.String(), upperIp.String()),
				AvailabilityZones: []string{cfg.Zone1, cfg.Zone2, cfg.Zone3},
				DNS:               metadataService,
			},
		},
	}
}

func networkCfg(cfg *config.Config) (networks commands.NetworksConfiguration, networkAssignment commands.NetworkAssignment) {
	networks = commands.NetworksConfiguration{
		ICMP: false,
		Networks: []commands.NetworkConfiguration{
			buildNetwork(cfg, cfg.MgmtSubnetName, cfg.MgmtSubnetCIDR, cfg.MgmtSubnetGateway),
			buildNetwork(cfg, cfg.ServicesSubnetName, cfg.ServicesSubnetCIDR, cfg.ServicesSubnetGateway),
			buildNetwork(cfg, cfg.ErtSubnetName, cfg.ErtSubnetCIDR, cfg.ErtSubnetGateway),
		},
	}

	networkAssignment = commands.NetworkAssignment{
		UserProvidedNetworkName: cfg.MgmtSubnetName,
		UserProvidedAZName:      cfg.Zone1,
	}

	return
}

func director() (director commands.DirectorConfiguration) {
	director = commands.DirectorConfiguration{
		NTPServers: metadataService, // gcp metadata service
	}

	return
}

func resources() commands.ResourceConfiguration {
	f := false
	return commands.ResourceConfiguration{
		DirectorResourceConfiguration: commands.DirectorResourceConfiguration{
			InternetConnected: &f,
		},
		CompilationResourceConfiguration: commands.CompilationResourceConfiguration{
			InternetConnected: &f,
		},
	}
}

func gcp(cfg *config.Config) commands.GCPIaaSConfiguration {
	return commands.GCPIaaSConfiguration{
		Project:              cfg.ProjectName,
		DefaultDeploymentTag: cfg.DeploymentTargetTag,
		AuthJSON:             "",
	}
}

func avalibilityZones(cfg *config.Config) commands.AvailabilityZonesConfiguration {
	return commands.AvailabilityZonesConfiguration{
		AvailabilityZones: []commands.AvailabilityZone{
			{Name: cfg.Zone1},
			{Name: cfg.Zone2},
			{Name: cfg.Zone3},
		},
	}
}
