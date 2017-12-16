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
		"9710",
		"37986",
		"398b1bf6bff84fef27d4e14e1d7693f37f3a338aec5beb10becce9570512e9f9",
	},
	product,
	&stemcell,
}

var smallRuntime = config.Tile{
	config.PivnetMetadata{
		"elastic-runtime",
		"9710",
		"37994",
		"c423b4b5983aaf324842d5d32f26a15ea5cd864f3c0e7d414ccf909a7140114d",
	},
	product,
	&stemcell,
}

var product = config.OpsManagerMetadata{
	"cf",
	"1.12.9",
}

var stemcell = config.StemcellMetadata{
	config.PivnetMetadata{"stemcells",
		"9610",
		"37837",
		"6c44e71b4eabc9665bcaa5db753a7ddaba909383a8c07001bef78882b42e8784"},
	"light-bosh-stemcell-3445.19-google-kvm-ubuntu-trusty-go_agent",
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
