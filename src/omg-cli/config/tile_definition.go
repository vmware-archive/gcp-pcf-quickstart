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

package config

// PivnetMetadata is metadata about a tile's information on the Pivotal Network.
type PivnetMetadata struct {
	Name      string
	ReleaseID int
	FileID    int
	Sha256    string
}

// StemcellMetadata is metadata about a stemcell.
type StemcellMetadata struct {
	PivnetMetadata
	StemcellName string
}

// OpsManagerMetadata is metadata associated with quickstart tiles.
type OpsManagerMetadata struct {
	Name         string
	Version      string
	DependsOnPAS bool // tiles which depend on PAS but don't specify so in their metadata
}

// Tile represents an Ops Manager tile.
type Tile struct {
	Pivnet   PivnetMetadata
	Product  OpsManagerMetadata
	Stemcell *StemcellMetadata
}
