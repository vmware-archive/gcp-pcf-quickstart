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
		53596,
		89124,
		"80e137622ca76868693b406114a2c7c1fdf6ce5db91c77a8d848d558d288fe5c",
	},
	config.OpsManagerMetadata{
		"stackdriver-nozzle",
		"2.0.1",
	},
	&config.StemcellMetadata{
		config.PivnetMetadata{"stemcells",
			28004,
			58600,
			"160e6ebc0f34af53543803eed4b7772fb7beaaa25cb3262ab170884e06feb721"},
		"light-bosh-stemcell-3421.36-google-kvm-ubuntu-trusty-go_agent",
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
