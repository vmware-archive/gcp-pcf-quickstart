package stackdriver_nozzle

import "omg-cli/tiles"

var Tile = tiles.Definition{
	tiles.PivnetDefinition{
		"gcp-stackdriver-nozzle",
		"5378",
		"20350",
		"b3156360159dbf20b5ac04b5ebd28c437741bc6d62bcb513587e72ac4e94fc18",
	},
	tiles.ProductDefinition{
		"stackdriver-nozzle",
		"1.0.3",
	},
}
