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

package environment

type OpsManagerQuery interface {
	// Tile returns a TileQuery interface if a tile is installed
	// or an error if it's not found
	Tile(name string) (TileQuery, error)

	// MustGetTile returns a TileQuery or panics if it is not installed
	MustGetTile(name string) TileQuery
}

type TileQuery interface {
	// Property returns the value of the given property set on the tile
	Property(name string) string
}
