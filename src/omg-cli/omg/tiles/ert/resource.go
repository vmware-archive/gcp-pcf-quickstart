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
	config.PivnetMetadata{
		"elastic-runtime",
		209723,
		242286,
		"2e2bee3cacf75604787854c6db2f8cd25462da23fb796901877df0bbebd70834",
	},
	product,
	&stemcell,
}

var smallRuntime = config.Tile{
	config.PivnetMetadata{
		"elastic-runtime",
		209723,
		242297,
		"9eaa7136576b969dc08806a65d3dab8aec9eda60b0aa71979bcdcf5469fecc1f",
	},
	product,
	&stemcell,
}

var product = config.OpsManagerMetadata{
	"cf",
	"2.3.2", // just the version as returned by ./util/pivnet_meta.sh
}

var stemcell = config.StemcellMetadata{
	config.PivnetMetadata{
		"stemcells-ubuntu-xenial",
		214330,
		247315,
		"c1202c333902e27a5cdaea360ea9ca9006bafb8a5e40d2d305a164fcb31d2e58",
	},
	"light-bosh-stemcell-97.28-google-kvm-ubuntu-xenial-go_agent",
}

type Tile struct{}

func (*Tile) Definition(envConfig *config.EnvConfig) config.Tile {
	if envConfig.SmallFootprint {
		return smallRuntime
	} else {
		return fullRuntime
	}
}

func (*Tile) BuiltIn() bool {
	return false
}
