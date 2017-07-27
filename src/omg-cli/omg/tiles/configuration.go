package tiles

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
