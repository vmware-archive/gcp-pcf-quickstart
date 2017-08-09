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
		"5993",
		"24044",
		"a1d248287fff3328459dedb10921394949f818e7b89f017803ac7d23a6c27bf2",
	},
	config.OpsManagerMetadata{
		"cf",
		"1.11.2",
	},
	nil,
}

type Tile struct{}

func (*Tile) Definition() config.Tile {
	return tile
}

func (*Tile) BuiltIn() bool {
	return false
}
