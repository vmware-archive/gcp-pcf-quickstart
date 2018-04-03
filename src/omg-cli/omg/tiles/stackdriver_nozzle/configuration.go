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

package stackdriver_nozzle

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
	Endpoint          tiles.Value `json:".properties.firehose_endpoint"`
	SkipSSLValidation tiles.Value `json:".properties.firehose_skip_ssl"`
	ServiceAccount    tiles.Value `json:".properties.service_account"`
	ProjectID         tiles.Value `json:".properties.project_id"`
}

type Resources struct {
	StackdriverNozzle tiles.Resource `json:"stackdriver-nozzle"`
}

func (t *Tile) Configure(envConfig *config.EnvConfig, cfg *config.Config, om *ops_manager.Sdk) error {
	if err := om.StageProduct(tile.Product); err != nil {
		return err
	}

	network := tiles.NetworkConfig(cfg.ServicesSubnetName, cfg)

	networkBytes, err := json.Marshal(&network)
	if err != nil {
		return err
	}

	properties := &Properties{
		Endpoint:          tiles.Value{fmt.Sprintf("https://api.sys.%s", cfg.DnsSuffix)},
		SkipSSLValidation: tiles.Value{skipSSLValidation},
		ServiceAccount:    tiles.Value{cfg.StackdriverNozzleServiceAccountKey},
		ProjectID:         tiles.Value{cfg.ProjectName},
	}

	propertiesBytes, err := json.Marshal(&properties)
	if err != nil {
		return err
	}

	vmType := ""
	if envConfig.SmallFootprint {
		vmType = "micro"
	}
	resoruces := Resources{
		StackdriverNozzle: tiles.Resource{
			InternetConnected: false,
			VmTypeId:          vmType,
		},
	}
	resorucesBytes, err := json.Marshal(&resoruces)
	if err != nil {
		return err
	}

	return om.ConfigureProduct(tile.Product.Name, string(networkBytes), string(propertiesBytes), string(resorucesBytes))
}
