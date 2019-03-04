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
		ReleaseID: 309043,
		FileID:    317727,
		Sha256:    "4ccd8080c3e383d4359208ca1849322a857dbb5d03cc96cfb16ec0fe7b90a6f5",
	},
	Product:  product,
	Stemcell: &stemcell,
}

var smallRuntime = config.Tile{
	Pivnet: config.PivnetMetadata{
		Name:      "elastic-runtime",
		ReleaseID: 309043,
		FileID:    317700,
		Sha256:    "04b381f83736e3d3af398bd6efbbf06dff986c6706371aac6b7f64863ee90d1f",
	},
	Product:  product,
	Stemcell: &stemcell,
}

var product = config.OpsManagerMetadata{
	Name:    "cf",
	Version: "2.4.4",
}

var stemcell = config.StemcellMetadata{
	PivnetMetadata: config.PivnetMetadata{
		Name:      "stemcells-ubuntu-xenial",
		ReleaseID: 301761,
		FileID:    313921,
		Sha256:    "b28e52be92ba3bba929807d395c81159cbb99b6feb0524aa41ab548cfd77b85b",
	},
	StemcellName: "light-bosh-stemcell-170.30-google-kvm-ubuntu-xenial-go_agent",
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
