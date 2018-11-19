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

package service_broker

import (
	"omg-cli/config"
)

var tile = config.Tile{
	config.PivnetMetadata{
		"gcp-service-broker",
		221375,
		255029,
		"beb501dd345b322123d0d608bc4c00dbd8343e066249b87a625e9c3cbc15059e",
	},
	config.OpsManagerMetadata{
		"gcp-service-broker",
		"4.1.0",
		true,
	},
	&config.StemcellMetadata{
		config.PivnetMetadata{
			"stemcells",
			232700,
			266558,
			"5d9a7325c05576b0dffa3dcbb7fd02c78a30c56a465cd0ebf39cbfb52f5ca566",
		},
		"light-bosh-stemcell-3586.56-google-kvm-ubuntu-trusty-go_agent",
	},
}

type Tile struct{}

func (*Tile) Definition(*config.EnvConfig) config.Tile {
	return tile
}

func (*Tile) BuiltIn() bool {
	return false
}
