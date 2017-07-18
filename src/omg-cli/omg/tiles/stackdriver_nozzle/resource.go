package stackdriver_nozzle

import (
	"errors"
	"omg-cli/config"
	"omg-cli/ops_manager"
)

var tile = config.Tile{
	config.PivnetMetadata{
		"gcp-stackdriver-nozzle",
		"5378",
		"20350",
		"b3156360159dbf20b5ac04b5ebd28c437741bc6d62bcb513587e72ac4e94fc18",
	},
	config.OpsManagerMetadata{
		"stackdriver-nozzle",
		"1.0.3",
	},
}

type Tile struct{}

func (Tile) Configure(cfg *config.Config, om *ops_manager.Sdk) error {
	return errors.New("stackdriver_nozzle: Tile: Conifgure NYI")
}

func (Tile) Definition() config.Tile {
	return tile
}
