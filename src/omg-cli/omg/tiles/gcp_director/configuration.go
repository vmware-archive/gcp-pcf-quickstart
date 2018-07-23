/*
 * Copyright 2017 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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

func (*Tile) Configure(envConfig *config.EnvConfig, cfg *config.Config, om *ops_manager.Sdk) error {
	networks, networkAssignment := networkCfg(cfg)

	return om.SetupBosh(gcp(cfg), director(cfg), avalibilityZones(cfg), networks, networkAssignment, resources(envConfig))
}

func buildNetwork(cfg *config.Config, name, cidrRange, gateway string, serviceNetwork bool) commands.NetworkConfiguration {
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
			buildNetwork(cfg, cfg.MgmtSubnetName, cfg.MgmtSubnetCIDR, cfg.MgmtSubnetGateway, false),
			buildNetwork(cfg, cfg.ServicesSubnetName, cfg.ServicesSubnetCIDR, cfg.ServicesSubnetGateway, false),
			buildNetwork(cfg, cfg.DynamicServicesSubnetName, cfg.DynamicServicesSubnetCIDR, cfg.DynamicServicesSubnetGateway, true),
			buildNetwork(cfg, cfg.ErtSubnetName, cfg.ErtSubnetCIDR, cfg.ErtSubnetGateway, false),
		},
	}

	networkAssignment = commands.NetworkAssignment{
		UserProvidedNetworkName: cfg.MgmtSubnetName,
		UserProvidedAZName:      cfg.Zone1,
	}

	return
}

func director(cfg *config.Config) (director commands.DirectorConfiguration) {
	t := true
	director = commands.DirectorConfiguration{
		NTPServers:                metadataService, // gcp metadata service
		EnableBoshDeployRetries:   &t,
		EnableVMResurrectorPlugin: &t,
		DatabaseType:              "external",
		ExternalDatabaseOptions: commands.ExternalDatabaseOptions{
			Host:     cfg.ExternalSqlIp,
			Database: cfg.OpsManagerSqlDbName,
			Username: cfg.OpsManagerSqlUsername,
			Password: cfg.OpsManagerSqlPassword,
			Port:     &cfg.ExternalSqlPort,
		},
	}

	return
}

func resources(envConfig *config.EnvConfig) commands.ResourceConfiguration {
	var instanceCount *int
	var compilation commands.CompilationInstanceType

	if envConfig.SmallFootprint {
		one := 1
		instanceCount = &one

		medium := "medium.mem"
		compilation.ID = &medium
	}

	// Healthwatch includes a C++ package that requires a large
	// ephemeral disk for compilation.
	if envConfig.IncludeHealthwatch {
		large := "large.disk"
		compilation.ID = &large
	}

	f := false
	return commands.ResourceConfiguration{
		DirectorResourceConfiguration: commands.DirectorResourceConfiguration{
			InternetConnected: &f,
		},
		CompilationResourceConfiguration: commands.CompilationResourceConfiguration{
			Instances:               instanceCount,
			CompilationInstanceType: compilation,
			InternetConnected:       &f,
		},
	}
}

func gcp(cfg *config.Config) commands.GCPIaaSConfiguration {
	return commands.GCPIaaSConfiguration{
		Project:              cfg.ProjectName,
		DefaultDeploymentTag: cfg.DeploymentTargetTag,
		AuthJSON:             cfg.OpsManagerServiceAccountKey,
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
