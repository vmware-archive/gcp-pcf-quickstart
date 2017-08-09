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
		"5378",
		"20350",
		"b3156360159dbf20b5ac04b5ebd28c437741bc6d62bcb513587e72ac4e94fc18",
	},
	config.OpsManagerMetadata{
		"stackdriver-nozzle",
		"1.0.3",
	},
	&config.StemcellMetadata{
		config.PivnetMetadata{"stemcells",
			"5884",
			"23445",
			"9b3175baf9d0b8b0bb1f37b029298e88cf352011aa632472a637d023bf928832"},
		"light-bosh-stemcell-3363.26-google-kvm-ubuntu-trusty-go_agent",
	},
}

type Tile struct {
	Logger *log.Logger
}

func (*Tile) Definition() config.Tile {
	return tile
}

func (*Tile) BuiltIn() bool {
	return false
}
