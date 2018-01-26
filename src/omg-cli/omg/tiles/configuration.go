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

package tiles

import (
	"omg-cli/config"
)

type AvalibilityZone struct {
	Name string `json:"name"`
}

type NetworkName struct {
	Name string `json:"name"`
}

type Network struct {
	SingletonAvalibilityZone AvalibilityZone   `json:"singleton_availability_zone"`
	OtherAvailabilityZones   []AvalibilityZone `json:"other_availability_zones"`
	Network                  NetworkName       `json:"network"`
}

func NetworkConfig(subnetName string, cfg *config.Config) Network {
	return Network{
		AvalibilityZone{cfg.Zone1},
		[]AvalibilityZone{{cfg.Zone1}, {cfg.Zone2}, {cfg.Zone3}},
		NetworkName{subnetName},
	}
}

type Value struct {
	Value string `json:"value"`
}

type IntegerValue struct {
	Value int `json:"value"`
}

type BooleanValue struct {
	Value bool `json:"value"`
}

type Secret struct {
	Value string `json:"secret"`
}

type SecretValue struct {
	Sec Secret `json:"value"`
}

type Certificate struct {
	PublicKey  string `json:"cert_pem"`
	PrivateKey string `json:"private_key_pem"`
}

type CertificateValue struct {
	Value Certificate `json:"value"`
}

type Resource struct {
	RouterNames       []string `json:"elb_names,omitempty"`
	Instances         *int     `json:"instances,omitempty"`
	InternetConnected bool     `json:"internet_connected"`
	VmTypeId          string   `json:"vm_type_id,omitempty"`
}
