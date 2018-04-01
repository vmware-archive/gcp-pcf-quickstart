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
		69794,
		104972,
		"4eb5940e0711f8de4ebf81f1031f1ff0c4c38704b07ecaa2f707962100296d0f",
	},
	product,
	&stemcell,
}

var smallRuntime = config.Tile{
	config.PivnetMetadata{
		"elastic-runtime",
		69794,
		104982,
		"0c9892f64b7a0314f3066aaf3e45f030736a36ebbb72907d0c0368164f91a11c",
	},
	product,
	&stemcell,
}

var product = config.OpsManagerMetadata{
	"cf",
	"2.1.0",
}

var stemcell = config.StemcellMetadata{
	config.PivnetMetadata{"stemcells",
		61940,
		96300,
		"b8aa5baed1d4cba84740a800ca465e29fe4469143e275044ac1b84e9ef71a326"},
	"light-bosh-stemcell-3541.8-google-kvm-ubuntu-trusty-go_agent",
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
