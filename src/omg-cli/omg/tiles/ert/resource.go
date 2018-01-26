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
		"28181",
		"58794",
		"e7d3d2826f11d0e5d1a42c92d1fb41e7842ca1921ccebff7f1e56defb330b08d",
	},
	product,
	&stemcell,
}

var smallRuntime = config.Tile{
	config.PivnetMetadata{
		"elastic-runtime",
		"28181",
		"58801",
		"2ab0cb6959933332fc44cca452f1dd1a35461f4c181cc16c54aae2392c177424",
	},
	product,
	&stemcell,
}

var product = config.OpsManagerMetadata{
	"cf",
	"1.12.11",
}

var stemcell = config.StemcellMetadata{
	config.PivnetMetadata{"stemcells",
		"28002",
		"58592",
		"fa6f3b8fe7e64987b628b17812989524550fea45a297fb7ead469c72d10f3b87"},
	"light-bosh-stemcell-3445.22-google-kvm-ubuntu-trusty-go_agent",
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
