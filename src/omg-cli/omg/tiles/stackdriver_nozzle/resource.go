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

package stackdriver_nozzle

import (
	"log"
	"omg-cli/config"
)

var tile = config.Tile{
	config.PivnetMetadata{
		"gcp-stackdriver-nozzle",
		"6586",
		"27540",
		"8aabde09128db9f7ec0d26cd7e05d6e56530fff3aaded0e88d4fbde83b3e55b4",
	},
	config.OpsManagerMetadata{
		"stackdriver-nozzle",
		"1.0.5",
	},
	&config.StemcellMetadata{
		config.PivnetMetadata{"stemcells",
			"7071",
			"30703",
			"6b11ee6dcc59fba486c0bb624b3c2431eb94c7b293d1af7fb98040b4a22aab33"},
		"light-bosh-stemcell-3421.26-google-kvm-ubuntu-trusty-go_agent",
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
