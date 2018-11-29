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

type ProductConfig struct {
	ProductName       string      `json:"product-name,omitempty"`
	ProductProperties interface{} `json:"product-properties,omitempty"`
	NetworkProperties interface{} `json:"network-properties,omitempty"`
	ResourceConfig    interface{} `json:"resource-config,omitempty"`
	ErrandConfig      interface{} `json:"errand-config,omitempty"`
}

type Errand struct {
	PostDeployState interface{} `json:"post-deploy-state,omitempty"`
	PreDeleteState interface{} `json:"pre-delete-state,omitempty"`
}

// AvailabilityZone is a shared config struct used for tile configuration.
type AvailabilityZone struct {
	Name string `json:"name"`
}

// NetworkName is a shared config struct used for tile configuration.
type NetworkName struct {
	Name string `json:"name"`
}

// Network is a shared config struct used for tile configuration.
type Network struct {
	SingletonAvailabilityZone AvailabilityZone   `json:"singleton_availability_zone"`
	OtherAvailabilityZones    []AvailabilityZone `json:"other_availability_zones"`
	Network                   NetworkName        `json:"network"`
	ODBNetwork                NetworkName        `json:"service_network"`
}

// NetworkConfig creates a Network.
func NetworkConfig(subnetName string, cfg *config.Config) Network {
	return Network{
		SingletonAvailabilityZone: AvailabilityZone{cfg.Zone1},
		OtherAvailabilityZones:    []AvailabilityZone{{cfg.Zone1}, {cfg.Zone2}, {cfg.Zone3}},
		Network:                   NetworkName{subnetName},
	}
}

// NetworkODBConfig creates a Network for an ODB network.
func NetworkODBConfig(subnetName string, cfg *config.Config, odbNetworkName string) Network {
	return Network{
		SingletonAvailabilityZone: AvailabilityZone{cfg.Zone1},
		OtherAvailabilityZones:    []AvailabilityZone{{cfg.Zone1}, {cfg.Zone2}, {cfg.Zone3}},
		Network:                   NetworkName{subnetName},
		ODBNetwork:                NetworkName{odbNetworkName},
	}
}

// Value is a shared config struct used for tile configuration.
type Value struct {
	Value string `json:"value"`
}

// IntegerValue is a shared config struct used for tile configuration.
type IntegerValue struct {
	Value int `json:"value"`
}

// BooleanValue is a shared config struct used for tile configuration.
type BooleanValue struct {
	Value bool `json:"value"`
}

// Secret is a shared config struct used for tile configuration.
type Secret struct {
	Value string `json:"secret"`
}

// SecretValue is a shared config struct used for tile configuration.
type SecretValue struct {
	Sec Secret `json:"value"`
}

// Certificate is a shared config struct used for tile configuration.
type Certificate struct {
	PublicKey  string `json:"cert_pem"`
	PrivateKey string `json:"private_key_pem"`
}

// CertificateConstruct is a shared config struct used for tile configuration.
type CertificateConstruct struct {
	Certificate Certificate `json:"certificate"`
	Name        string      `json:"name"`
}

// CertificateValue is a shared config struct used for tile configuration.
type CertificateValue struct {
	Value []CertificateConstruct `json:"value"`
}

// OldCertificateValue is a shared config struct used for tile configuration.
type OldCertificateValue struct {
	Value Certificate `json:"value"`
}

// KeyStruct is a shared config struct used for tile configuration.
type KeyStruct struct {
	Secret string `json:"secret"`
}

// EncryptionKey is a shared config struct used for tile configuration.
type EncryptionKey struct {
	Name    string    `json:"name"`
	Key     KeyStruct `json:"key"`
	Primary bool      `json:"primary"`
}

// EncryptionKeyValue is a shared config struct used for tile configuration.
type EncryptionKeyValue struct {
	Value []EncryptionKey `json:"value"`
}

// Resource is a shared config struct used for tile configuration.
type Resource struct {
	RouterNames       []string `json:"elb_names,omitempty"`
	Instances         *int     `json:"instances,omitempty"`
	InternetConnected bool     `json:"internet_connected"`
	VMTypeID          string   `json:"vm_type_id,omitempty"`
}
