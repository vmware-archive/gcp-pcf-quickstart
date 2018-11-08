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

package healthwatch

import (
	"log"
	"omg-cli/config"
)

var tile = config.Tile{
	config.PivnetMetadata{
		"p-healthwatch",
		106732,
		145945,
		"b5d01673c5c911e022a2a725fc27bc2085bc422155277e422ff65d5a5469ad55",
	},
	config.OpsManagerMetadata{
		"p-healthwatch",
		"1.2.2-build.10",
		true,
	},
	&config.StemcellMetadata{
		config.PivnetMetadata{
			"stemcells",
			129476,
			161610,
			"1c35a0d9b8ae5899423bc5b1160600e9608e85bdf5ade9543c954bf3880bbb9b"},
		"light-bosh-stemcell-3468.51-google-kvm-ubuntu-trusty-go_agent",
	},
}

type Tile struct {
	Logger *log.Logger
}

func (*Tile) Definition(*config.EnvConfig) config.Tile {
	return tile
}

func (*Tile) BuiltIn() bool {
	return false
}
