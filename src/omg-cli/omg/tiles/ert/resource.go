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
		52010,
		85259,
		"3b462800f3d1478f2fec8ef7230c7c657c5e540765d3489639ccc2345786c1cb",
	},
	product,
	&stemcell,
}

var smallRuntime = config.Tile{
	config.PivnetMetadata{
		"elastic-runtime",
		52010,
		85268,
		"098e07101c8233d041f643b83d0430efe6102216a27e016355c9a9f6a3e1aa18",
	},
	product,
	&stemcell,
}

var product = config.OpsManagerMetadata{
	"cf",
	"2.0.6",
}

var stemcell = config.StemcellMetadata{
	config.PivnetMetadata{"stemcells",
		36314,
		67872,
		"6c966883018e34edc8c0c61a48b7aa07582571e39e37f9065ff58eff4f4b4423"},
	"light-bosh-stemcell-3468.21-google-kvm-ubuntu-trusty-go_agent",
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
