package bosh_director

import (
	"fmt"
	"net"
	"omg-cli/config"

	"github.com/pivotal-cf/om/commands"
)

const (
	metadataService = "169.254.169.254"
)

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

func Network(cfg *config.Config) (networks commands.NetworksConfiguration, networkAssignment commands.NetworkAssignment) {
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

func Director() (director commands.DirectorConfiguration) {
	director = commands.DirectorConfiguration{
		NTPServers: metadataService, // GCP metadata service
	}

	return
}

func Resources() commands.ResourceConfiguration {
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

func GCP(cfg *config.Config) commands.GCPIaaSConfiguration {
	return commands.GCPIaaSConfiguration{
		Project:              cfg.ProjectName,
		DefaultDeploymentTag: cfg.DeploymentTargetTag,
		AuthJSON:             "",
	}
}

func AvalibilityZones(cfg *config.Config) commands.AvailabilityZonesConfiguration {
	return commands.AvailabilityZonesConfiguration{
		AvailabilityZones: []commands.AvailabilityZone{
			{Name: cfg.Zone1},
			{Name: cfg.Zone2},
			{Name: cfg.Zone3},
		},
	}
}
