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

import (
	"fmt"

	"omg-cli/ops_manager"

	"github.com/pivotal-cf/om/api"
)

type DirectorProperties map[string]map[string]interface{}

type OpsManagerQuery interface {
	// Tile returns a TileQuery interface if a tile is installed
	// or an error if it's not found
	Tile(name string) (TileQuery, error)

	// MustGetTile returns a TileQuery or panics if it is not installed
	MustGetTile(name string) TileQuery

	// Director returns information about the deployed BOSH director
	Director() DirectorProperties
}

type TileQuery interface {
	// Property returns the value of the given property set on the tile
	Property(name string) api.ResponseProperty

	// Resource returns information about a given resource for a job declared by the tile
	Resource(name string) api.JobProperties
}

type liveOpsManager struct {
	sdk *ops_manager.Sdk
}

func (lom *liveOpsManager) Tile(name string) (TileQuery, error) {
	props, err := lom.sdk.GetProduct(name)
	if err != nil {
		return nil, fmt.Errorf("getting product propeties: %v", err)
	}

	return &liveTileQuery{name, props, lom.sdk}, nil
}

func (lom *liveOpsManager) MustGetTile(name string) TileQuery {
	tile, err := lom.Tile(name)
	if err != nil {
		panic(fmt.Errorf("expected tile: %v", err))
	}

	return tile
}

func (lom *liveOpsManager) Director() DirectorProperties {
	prop, err := lom.sdk.GetDirector()
	if err != nil {
		panic(fmt.Errorf("retreving director: %v", err))
	}

	return prop
}

type liveTileQuery struct {
	name  string
	props *ops_manager.ProductProperties
	sdk   *ops_manager.Sdk
}

func (ltq *liveTileQuery) Property(name string) api.ResponseProperty {
	return ltq.props.Properties[name]
}

func (ltq *liveTileQuery) Resource(jobID string) api.JobProperties {
	resource, err := ltq.sdk.GetResource(ltq.name, jobID)
	if err != nil {
		panic(err)
	}
	return *resource
}
