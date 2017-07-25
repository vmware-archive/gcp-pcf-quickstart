package gcp_director

import (
	"omg-cli/config"
)

type Tile struct{}

func (Tile) Definition() config.Tile {
	panic("not applicable to built in tiles")
	return config.Tile{}
}
func (Tile) BuiltIn() bool {
	return true
}
