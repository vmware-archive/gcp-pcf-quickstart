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
		191048,
		221334,
		"38b7eea0437af30901803f2ed147e0a167df54e3f4e43096eea11112505efc35",
	},
	config.OpsManagerMetadata{
		"gcp-service-broker",
		"4.0.0",
		true,
	},
	&config.StemcellMetadata{
		config.PivnetMetadata{"stemcells",
			183841,
			214146,
			"93fe455d2ab111cdb5f6902cf8f8552e406e6665f4b46bba5fca37ac47aa0ecd"},
		"light-bosh-stemcell-3586.42-google-kvm-ubuntu-trusty-go_agent",
	},
}

type Tile struct{}

func (*Tile) Definition(*config.EnvConfig) config.Tile {
	return tile
}

func (*Tile) BuiltIn() bool {
	return false
}
