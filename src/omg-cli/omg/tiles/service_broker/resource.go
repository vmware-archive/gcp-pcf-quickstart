package service_broker

import (
	"errors"
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
}

type Tile struct{}

func (Tile) Configure(cfg *config.Config, om *ops_manager.Sdk) error {
	return errors.New("service_broker: Tile: Conifgure NYI")
}

func (Tile) Definition() config.Tile {
	return tile
}
