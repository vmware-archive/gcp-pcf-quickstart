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

var tile = config.Tile{
	config.PivnetMetadata{
		"elastic-runtime",
		"7295",
		"32099",
		"869c0b8cdae06fc6bc1d8b708d27e88d6425163db48d3f8b867b1f9c13253d87",
	},
	config.OpsManagerMetadata{
		"cf",
		"1.12.3",
	},
	&config.StemcellMetadata{
		config.PivnetMetadata{"stemcells",
			"7069",
			"30701",
			"978683f9243fce813a2b4fe0a9e62395e93df8cbde9258bb676bc3c199df398d"},
		"light-bosh-stemcell-3445.11-google-kvm-ubuntu-trusty-go_agent",
	},
}

type Tile struct{}

func (*Tile) Definition() config.Tile {
	return tile
}

func (*Tile) BuiltIn() bool {
	return false
}
