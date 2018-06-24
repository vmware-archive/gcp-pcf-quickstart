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
		118424,
		158939,
		"214da7e21a57cf851fd7c842542157af0f26fc3173ca8746377d3c890ec64243",
	},
	product,
	&stemcell,
}

var smallRuntime = config.Tile{
	config.PivnetMetadata{
		"elastic-runtime",
		118424,
		158946,
		"ec5c879be1ca3349aef003e877083a66add34dfd3a86e6dca1dbe99f445cb2b9",
	},
	product,
	&stemcell,
}

var product = config.OpsManagerMetadata{
	"cf",
	"2.1.7",
}

var stemcell = config.StemcellMetadata{
	config.PivnetMetadata{
		"stemcells",
		106151,
		145286,
		"31f0b81919c58f71ba3427b29bdc02099808da378f938ed3d930ff5623c7d0d6"},
	"light-bosh-stemcell-3541.30-google-kvm-ubuntu-trusty-go_agent",
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
