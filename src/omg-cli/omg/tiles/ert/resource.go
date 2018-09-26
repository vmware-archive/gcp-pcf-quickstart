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
		188503,
		218723,
		"20906dcb6352e69b2e0fc056fdeefcd6f2b6ac45f7dca6204abf0aea87c35e2f",
	},
	product,
	&stemcell,
}

var smallRuntime = config.Tile{
	config.PivnetMetadata{
		"elastic-runtime",
		188503,
		218738,
		"aabe0a5b6297ce4375db39559bb631a551917f11f93fb5f1b8756112a05a0858",
	},
	product,
	&stemcell,
}

var product = config.OpsManagerMetadata{
	"cf",
	"2.3.0",
}

var stemcell = config.StemcellMetadata{
	config.PivnetMetadata{
		"stemcells-ubuntu-xenial",
		194743,
		225464,
		"c47d062781073c64c3b978262331e6e3d45bd7074a57d437c50e7b99930c4581"},
	"light-bosh-stemcell-97.18-google-kvm-ubuntu-xenial-go_agent",
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
