package service_broker

import "omg-cli/tiles"

var Tile = tiles.Definition{
	tiles.PivnetDefinition{
		"gcp-service-broker",
		"5563",
		"21222",
		"81dd57e6a98b62cf27336b84ffac3051feafe23fc28f3e14d2b61dc8982043c1",
	},
	tiles.ProductDefinition{
		"gcp-service-broker",
		"3.4.1",
	},
}
