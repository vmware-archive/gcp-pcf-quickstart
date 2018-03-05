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
		24491,
		54209,
		"e192cff35a92e9b78ee8c626bf0a5b59e86946d93111d953a39e211c9fb649b8",
	},
	config.OpsManagerMetadata{
		"gcp-service-broker",
		"3.6.0",
	},
	&config.StemcellMetadata{
		config.PivnetMetadata{"stemcells",
			28003,
			58587,
			"af238d0d9d94b18da32d302512831d83aec00312a18bb528b351144e0f281f0e"},
		"light-bosh-stemcell-3468.17-google-kvm-ubuntu-trusty-go_agent",
	},
}

type Tile struct{}

func (*Tile) Definition(*config.EnvConfig) config.Tile {
	return tile
}

func (*Tile) BuiltIn() bool {
	return false
}
