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

package ert

import (
	"omg-cli/config"
)

var fullRuntime = config.Tile{
	Pivnet: config.PivnetMetadata{
		Name:      "elastic-runtime",
		ReleaseID: 259105,
		FileID:    279697,
		Sha256:    "30a827ec8106f359ee49448707c2304b87e73ad4a422baff15038b8acb1525c7",
	},
	Product:  product,
	Stemcell: &stemcell,
}

var smallRuntime = config.Tile{
	Pivnet: config.PivnetMetadata{
		Name:      "elastic-runtime",
		ReleaseID: 259105,
		FileID:    279697,
		Sha256:    "9756289d1b4f7c9ad565755395cad98dce77917432221ea6c789719696521887",
	},
	Product:  product,
	Stemcell: &stemcell,
}

var product = config.OpsManagerMetadata{
	Name:    "cf",
	Version: "2.4.1",
}

var stemcell = config.StemcellMetadata{
	PivnetMetadata: config.PivnetMetadata{
		Name:      "stemcells-ubuntu-xenial",
		ReleaseID: 276954,
		FileID:    288597,
		Sha256:    "9f9d1fab2b5165c5cee52c6b0a84df96c6861355135fa8a99748f59832466b4e",
	},
	StemcellName: "light-bosh-stemcell-170.19-google-kvm-ubuntu-xenial-go_agent",
}

// Tile is the tile for the Pivotal Application Service.
type Tile struct{}

// Definition satisfies TileInstaller interface.
func (*Tile) Definition(envConfig *config.EnvConfig) config.Tile {
	if envConfig.SmallFootprint {
		return smallRuntime
	}

	return fullRuntime
}

// BuiltIn satisfies TileInstaller interface.
func (*Tile) BuiltIn() bool {
	return false
}
