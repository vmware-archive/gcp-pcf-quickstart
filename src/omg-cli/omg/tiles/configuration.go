package tiles

import (
	"omg-cli/config"
)

type AvalibilityZone struct {
	Name string `json:"name"`
}

type NetworkName struct {
	Name string `json:"name"`
}

type Network struct {
	SingletonAvalibilityZone AvalibilityZone   `json:"singleton_availability_zone"`
	OtherAvailabilityZones   []AvalibilityZone `json:"other_availability_zones"`
	Network                  NetworkName       `json:"network"`
}

func NetworkConfig(subnetName string, cfg *config.Config) Network {
	return Network{
		AvalibilityZone{cfg.Zone1},
		[]AvalibilityZone{{cfg.Zone1}, {cfg.Zone2}, {cfg.Zone3}},
		NetworkName{subnetName},
	}
}
