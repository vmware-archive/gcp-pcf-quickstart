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

package healthwatch

import (
	"log"

	"omg-cli/config"
)

var tile = config.Tile{
	Pivnet: config.PivnetMetadata{
		Name:      "p-healthwatch",
		ReleaseID: 161272,
		FileID:    194184,
		Sha256:    "75a17ff1c6044de391b327275b944ffb524650e6b40bc2d4f68c219940c37107",
	},
	Product: config.OpsManagerMetadata{
		Name:         "p-healthwatch",
		Version:      "1.3.2-build.9",
		DependsOnPAS: true,
	},
	Stemcell: &config.StemcellMetadata{
		PivnetMetadata: config.PivnetMetadata{
			Name:      "stemcells",
			ReleaseID: 224700,
			FileID:    258680,
			Sha256:    "c4a3be0d143e25e921b090e256ce669f990b10b5ba4181ccacd49338b5200881"},
		StemcellName: "light-bosh-stemcell-3541.59-google-kvm-ubuntu-trusty-go_agent",
	},
}

// Tile is the tile for the Healthwatch.
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
