/*
 * Copyright 2017 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package stackdriver_nozzle

import (
	"log"

	"omg-cli/config"
)

var tile = config.Tile{
	Pivnet: config.PivnetMetadata{
		Name:      "gcp-stackdriver-nozzle",
		ReleaseID: 53596,
		FileID:    89124,
		Sha256:    "80e137622ca76868693b406114a2c7c1fdf6ce5db91c77a8d848d558d288fe5c",
	},
	Product: config.OpsManagerMetadata{
		Name:         "stackdriver-nozzle",
		Version:      "2.0.1",
		DependsOnPAS: true,
	},
	Stemcell: &config.StemcellMetadata{
		PivnetMetadata: config.PivnetMetadata{
			Name:      "stemcells",
			ReleaseID: 214323,
			FileID:    247292,
			Sha256:    "8c6caeae37711aaf12b4fefba06c348cde5631e872e8892553ddb26514a3953a",
		},
		StemcellName: "light-bosh-stemcell-3468.78-google-kvm-ubuntu-trusty-go_agent",
	},
}

// Tile is the tile for the Stackdriver Nozzle.
type Tile struct {
	Logger *log.Logger
}

// Definition satisfies TileInstaller interface.
func (*Tile) Definition(*config.EnvConfig) config.Tile {
	return tile
}

// BuiltIn satisfies TileInstaller interface.
func (*Tile) BuiltIn() bool {
	return false
}
