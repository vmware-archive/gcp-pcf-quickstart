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

package healthwatch

import (
	"encoding/json"
	"fmt"
	"omg-cli/config"
	"omg-cli/omg/tiles"
	"omg-cli/ops_manager"
)

const (
	skipSSLValidation = "true"
)

type Properties struct {
	OpsManagerURL     tiles.Value `json:".properties.opsman.enable.url"`
	BoshHealthCheckAZ tiles.Value `json:".healthwatch-forwarder.health_check_az"`
}

type Resources struct {
}

func (t *Tile) Configure(envConfig *config.EnvConfig, cfg *config.Config, om *ops_manager.Sdk) error {
	if err := om.StageProduct(tile.Product); err != nil {
		return err
	}

	network := tiles.NetworkODBConfig(cfg.ServicesSubnetName, cfg, cfg.DynamicServicesSubnetName)

	networkBytes, err := json.Marshal(&network)
	if err != nil {
		return err
	}

	properties := &Properties{
		OpsManagerURL:     tiles.Value{fmt.Sprintf("https://opsman.%s", cfg.DnsSuffix)},
		BoshHealthCheckAZ: tiles.Value{cfg.Zone1},
	}

	propertiesBytes, err := json.Marshal(&properties)
	if err != nil {
		return err
	}

	resources := Resources{}
	resourcesBytes, err := json.Marshal(&resources)
	if err != nil {
		return err
	}

	return om.ConfigureProduct(tile.Product.Name, string(networkBytes), string(propertiesBytes), string(resourcesBytes))
}
