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
		ReleaseID: 220833,
		FileID:    254457,
		Sha256:    "5540900a3626b092bffdb01b530791a116cf5f1022fd1b048edaeea4424318fd",
	},
	Product:  product,
	Stemcell: &stemcell,
}

var smallRuntime = config.Tile{
	Pivnet: config.PivnetMetadata{
		Name:      "elastic-runtime",
		ReleaseID: 220833,
		FileID:    254473,
		Sha256:    "1ab242bff8f95598193b0c742b7d6a520628ebeb682fd949d18da5ef6c8e5c7a",
	},
	Product:  product,
	Stemcell: &stemcell,
}

var product = config.OpsManagerMetadata{
	Name:    "cf",
	Version: "2.3.3",
}

var stemcell = config.StemcellMetadata{
	PivnetMetadata: config.PivnetMetadata{
		Name:      "stemcells-ubuntu-xenial",
		ReleaseID: 226360,
		FileID:    260592,
		Sha256:    "a23bd96427043afd34a62223f11aebf3177c7f8f0c2c46006c2168120effa099",
	},
	StemcellName: "light-bosh-stemcell-97.32-google-kvm-ubuntu-xenial-go_agent",
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
