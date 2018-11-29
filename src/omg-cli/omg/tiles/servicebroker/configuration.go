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

package servicebroker

import (
	"encoding/json"

	"omg-cli/config"
	"omg-cli/omg/tiles"
	"omg-cli/opsman"
)

type properties struct {
	ServiceAccountKey tiles.Value       `json:".properties.root_service_account_json"`
	DatabaseHost      tiles.Value       `json:".properties.db_host"`
	DatabaseUsername  tiles.Value       `json:".properties.db_username"`
	DatabasePassword  tiles.SecretValue `json:".properties.db_password"`
}

// Configure satisfies TileInstaller interface.
func (*Tile) Configure(envConfig *config.EnvConfig, cfg *config.Config, om *opsman.Sdk) error {
	if err := om.StageProduct(tile.Product); err != nil {
		return err
	}

	network := tiles.NetworkConfig(cfg.ServicesSubnetName, cfg)

	properties := properties{
		ServiceAccountKey: tiles.Value{Value: cfg.ServiceBrokerServiceAccountKey},
		DatabaseHost:      tiles.Value{Value: cfg.ServiceBrokerDbIP},
		DatabaseUsername:  tiles.Value{Value: cfg.ServiceBrokerDbUsername},
		DatabasePassword:  tiles.SecretValue{Sec: tiles.Secret{Value: cfg.ServiceBrokerDbPassword}},
	}

	productConfig := tiles.ProductConfig{
		ProductName: tile.Product.Name,
		NetworkProperties: network,
		ProductProperties: properties,
		ResourceConfig: struct {}{},
	}

	configBytes, err := json.Marshal(productConfig)
	if err != nil {
		return err
	}

	return om.ConfigureProduct(tile.Product.Name, configBytes)
}
