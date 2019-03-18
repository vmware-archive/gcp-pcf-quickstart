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

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type TerraformConfigSchema struct {
	MgmtSubnetName                string `json:"management_subnet_name"`
	ErtSubnetName                 string `json:"ert_subnet_name"`
	ServicesSubnetName            string `json:"services_subnet_name"`
	JumpboxIP                     string `json:"jumpbox_public_ip"`
	OpsManagerHostname            string `json:"ops_manager_dns"`
	OpsManagerUsername            string `json:"ops_manager_username"`
	OpsManagerPassword            string `json:"ops_manager_password"`
	OpsManagerDecryptionPhrase    string `json:"ops_manager_decryption_phrase"`
	OpsManagerSkipSSLVerification string `json:"ops_manager_skip_ssl_verify"`
	Raw                           map[string]interface{}

	OpsManager OpsManagerCredentials
}

// TerraformFromEnvDirectory creates a Terraform config from a directory.
func TerraformFromEnvDirectory(path string) (*Config, error) {
	config, err := fromTerraform(filepath.Join(path, TerraformOutputFile))
	if err != nil {
		return nil, fmt.Errorf("creating Terraform config from directory %s: %v", path, err)
	}

	return config, nil
}

func decode(input string) string {
	res, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		panic("unable to decode bas64 key")
	}
	return string(res)
}

func fromTerraform(filename string) (*Config, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	flattened, err := flattenTerraform(file)
	if err != nil {
		return nil, err
	}
	flattendStr, err := json.Marshal(flattened)

	if err != nil {
		return nil, err
	}

	hydratedCfg := TerraformConfigSchema{}
	err = json.Unmarshal(flattendStr, &hydratedCfg)
	if err != nil {
		return nil, err
	}
	hydratedCfg.OpsManager.Username = hydratedCfg.OpsManagerUsername
	hydratedCfg.OpsManager.Password = hydratedCfg.OpsManagerPassword
	hydratedCfg.OpsManager.DecryptionPhrase = hydratedCfg.OpsManagerDecryptionPhrase
	if hydratedCfg.OpsManagerSkipSSLVerification == "true" {
		hydratedCfg.OpsManager.SkipSSLVerification = true
	}
	hydratedCfg.Raw = flattened
	cfg := Config(hydratedCfg)

	return &cfg, nil
}

type terraformValue struct {
	Value interface{} `json:"value"`
}

func flattenTerraform(contents []byte) (map[string]interface{}, error) {
	res := map[string]interface{}{}

	tf := map[string]terraformValue{}

	err := json.Unmarshal(contents, &tf)
	if err != nil {
		return nil, err
	}

	for k, v := range tf {
		if str, ok := v.Value.(string); ok {
			if strings.HasSuffix(k, "_base64") {
				res[strings.TrimSuffix(k, "_base64")] = decode(str)
			} else {
				res[k] = str
			}
		} else if arr, ok := v.Value.([]interface{}); ok {
			for i, entry := range arr {
				res[fmt.Sprintf("%s_%d", k, i)] = entry.(string)
			}
		} else {
			return nil, fmt.Errorf("encountered unknown type in terraform config: %v", v.Value)
		}
	}
	// fmt.Println(res)
	return res, nil
}
