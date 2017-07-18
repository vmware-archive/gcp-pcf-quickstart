package ert

import "omg-cli/tiles"

var Tile = tiles.Definition{
	tiles.PivnetDefinition{
		"elastic-runtime",
		"5993",
		"24044",
		"a1d248287fff3328459dedb10921394949f818e7b89f017803ac7d23a6c27bf2",
	},
	tiles.ProductDefinition{
		"cf",
		"1.11.2",
	},
}
