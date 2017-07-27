package ert

import (
	"omg-cli/config"
)

var tile = config.Tile{
	config.PivnetMetadata{
		"elastic-runtime",
		"5993",
		"24044",
		"a1d248287fff3328459dedb10921394949f818e7b89f017803ac7d23a6c27bf2",
	},
	config.OpsManagerMetadata{
		"cf",
		"1.11.2",
	},
	nil,
}

type Tile struct{}

func (Tile) Definition() config.Tile {
	return tile
}

func (Tile) BuiltIn() bool {
	return false
}
