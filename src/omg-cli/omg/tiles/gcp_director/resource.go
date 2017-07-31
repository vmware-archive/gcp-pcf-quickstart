package gcp_director

import (
	"omg-cli/config"
)

var tile = config.Tile{
	Product: config.OpsManagerMetadata{
		Name: "BOSH Director",
	},
}

type Tile struct{}

func (*Tile) Definition() config.Tile {
	return tile
}
func (*Tile) BuiltIn() bool {
	return true
}
