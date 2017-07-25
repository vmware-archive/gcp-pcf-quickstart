package tiles

import (
	"omg-cli/config"
	"omg-cli/ops_manager"
)

type TileInstaller interface {
	Definition() config.Tile
	Configure(cfg *config.Config, om *ops_manager.Sdk) error
	BuiltIn() bool
}
