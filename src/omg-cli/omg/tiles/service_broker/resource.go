package service_broker

import (
	"omg-cli/config"
	"omg-cli/ops_manager"
)

var tile = config.Tile{
	config.PivnetMetadata{
		"gcp-service-broker",
		"5563",
		"21222",
		"81dd57e6a98b62cf27336b84ffac3051feafe23fc28f3e14d2b61dc8982043c1",
	},
	config.OpsManagerMetadata{
		"gcp-service-broker",
		"3.4.1",
	},

	&config.StemcellMetadata{
		config.PivnetMetadata{"stemcells",
			"5884",
			"23445",
			"9b3175baf9d0b8b0bb1f37b029298e88cf352011aa632472a637d023bf928832"},
		"light-bosh-stemcell-3363.26-google-kvm-ubuntu-trusty-go_agent",
	},
}

type Tile struct{}

func (Tile) Configure(cfg *config.Config, om *ops_manager.Sdk) error {
	println("TODO: Configure Service Broker. Skipping.")
	return nil
}

func (Tile) Definition() config.Tile {
	return tile
}

func (Tile) BuiltIn() bool {
	return false
}
