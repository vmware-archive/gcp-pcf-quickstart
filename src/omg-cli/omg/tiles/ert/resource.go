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

var tile = config.Tile{
	config.PivnetMetadata{
		"elastic-runtime",
		"6434",
		"26790",
		"bdfdd6e7fce47de2d90377a50ed78bb11ef45bc4cf7b8c55f5d37046ff481533",
	},
	config.OpsManagerMetadata{
		"cf",
		"1.11.8",
	},
	&config.StemcellMetadata{
		config.PivnetMetadata{"stemcells",
			"6666",
			"27817",
			"742f683a87405af46ef9eb7546f12d1461da91898ac23b63b8cc65624ecebea4"},
		"light-bosh-stemcell-3421.20-google-kvm-ubuntu-trusty-go_agent",
	},
}

type Tile struct{}

func (*Tile) Definition() config.Tile {
	return tile
}

func (*Tile) BuiltIn() bool {
	return false
}
