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
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type EnvConfig struct {
	DnsZoneName    string
	ProjectID      string
	BaseImageURL   string
	EnvName        string
	Region         string
	PivnetApiToken string
	Zone1          string
	Zone2          string
	Zone3          string
	SmallFootprint bool
}

func DefaultEnvConfig() (*EnvConfig, error) {
	c := &EnvConfig{
		DnsZoneName:    "pcf-zone",
		BaseImageURL:   "https://storage.cloud.google.com/ops-manager-us/pcf-gcp-2.1-build.204.tar.gz",
		EnvName:        "pcf",
		Region:         "us-east1",
		Zone1:          "us-east1-b",
		Zone2:          "us-east1-c",
		Zone3:          "us-east1-d",
		SmallFootprint: true,
	}

	return c, nil
}

func ConfigFromEnvDirectory(path string) (*EnvConfig, error) {
	return fromEnvironment(filepath.Join(path, EnvConfigFile))
}

func fromEnvironment(filename string) (*EnvConfig, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cfg := &EnvConfig{}
	if err = json.Unmarshal(file, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
