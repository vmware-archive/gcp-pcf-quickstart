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
	cfGuid            = "cf-96162609e45ab17f57bd"
	uaaCredential     = ".uaa.admin_credentials"
	skipSSLValidation = "true"
)

type Properties struct {
	Endpoint          Value       `json:".properties.firehose_endpoint"`
	Username          Value       `json:".properties.firehose_username"`
	Password          SecretValue `json:".properties.firehose_password"`
	SkipSSLValidation Value       `json:".properties.firehose_skip_ssl"`
	ServiceAccount    Value       `json:".properties.service_account"`
	ProjectID         Value       `json:".properties.project_id"`
}

type Value struct {
	Value string `json:"value"`
}

type BoolValue struct {
	Value bool `json:"value"`
}

type Secret struct {
	Value string `json:"secret"`
}

type SecretValue struct {
	Sec Secret `json:"value"`
}

func (t *Tile) Configure(cfg *config.Config, om *ops_manager.Sdk) error {
	// TODO(jrjohnson): Should we create a new, scoped UAA user here?
	// What about `.uaa.stackdriver_nozzle_credentials` Is that a real user?
	cred, err := om.GetCredentials(cfGuid, uaaCredential)
	if err != nil {
		// TODO(jrjohnson): We should check if ERT has been successfully deployed
		t.Logger.Printf("stackdriver nozzle: error getting credentials. skipping configuration. If ERT isn't deployed yet then this ignore this error: %v", err)
		return nil
	}

	if err := om.StageProduct(tile.Product); err != nil {
		return err
	}

	network := tiles.NetworkConfig(cfg.ServicesSubnetName, cfg)

	networkBytes, err := json.Marshal(&network)
	if err != nil {
		return err
	}

	properties := &Properties{
		Username:          Value{cred.Identity},
		Password:          SecretValue{Secret{cred.Password}},
		Endpoint:          Value{fmt.Sprintf("https://api.sys.%s", cfg.DnsSuffix)},
		SkipSSLValidation: Value{skipSSLValidation},
		ServiceAccount:    Value{cfg.StackdriverNozzleServiceAccountKey},
		ProjectID:         Value{cfg.ProjectName},
	}

	propertiesBytes, err := json.Marshal(&properties)
	if err != nil {
		return err
	}

	resoruces := "{}"

	return om.ConfigureProduct(tile.Product.Name, string(networkBytes), string(propertiesBytes), resoruces)
}
