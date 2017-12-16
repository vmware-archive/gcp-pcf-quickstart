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
		"5563",
		"21222",
		"81dd57e6a98b62cf27336b84ffac3051feafe23fc28f3e14d2b61dc8982043c1",
	},
	config.OpsManagerMetadata{
		"gcp-service-broker",
		"3.4.1",
	},
	&config.StemcellMetadata{
		config.PivnetMetadata{"stemcells",
			"5884",
			"23445",
			"9b3175baf9d0b8b0bb1f37b029298e88cf352011aa632472a637d023bf928832"},
		"light-bosh-stemcell-3363.26-google-kvm-ubuntu-trusty-go_agent",
	},
}

type Tile struct{}

func (*Tile) Definition(*config.EnvConfig) config.Tile {
	return tile
}

func (*Tile) BuiltIn() bool {
	return false
}
