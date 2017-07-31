package stackdriver_nozzle

import (
	"log"
	"omg-cli/config"
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
	&config.StemcellMetadata{
		config.PivnetMetadata{"stemcells",
			"5884",
			"23445",
			"9b3175baf9d0b8b0bb1f37b029298e88cf352011aa632472a637d023bf928832"},
		"light-bosh-stemcell-3363.26-google-kvm-ubuntu-trusty-go_agent",
	},
}

type Tile struct {
	Logger *log.Logger
}

func (*Tile) Definition() config.Tile {
	return tile
}

func (*Tile) BuiltIn() bool {
	return false
}
