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

package service_broker

import (
	"encoding/json"
	"omg-cli/config"
	"omg-cli/omg/tiles"
	"omg-cli/ops_manager"
)

type Properties struct {
	ServiceAccountKey tiles.Value       `json:".properties.root_service_account_json"`
	DatabaseHost      tiles.Value       `json:".properties.db_host"`
	DatabaseUsername  tiles.Value       `json:".properties.db_username"`
	DatabasePassword  tiles.SecretValue `json:".properties.db_password"`

	CloudSQLPlans CloudSQLPlansValue `json:".properties.cloudsql_custom_plans"`
	BigTablePlans BigTablePlansValue `json:".properties.bigtable_custom_plans"`
}

type CloudSQLPlansValue struct {
	Plans []CloudSQLPlan `json:"value"`
}

type CloudSQLPlan struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Service     string `json:"service"`
	Tier        string `json:"tier"`
	MaxDiskSize string `json:"max_disk_size"`
	PricingPlan string `json:"pricing_plan"`
}

type BigTablePlansValue struct {
	Plans []BigTablePlan `json:"value"`
}

type BigTablePlan struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Service     string `json:"service"`
	StorageType string `json:"storage_type"`
	NodeCount   string `json:"num_nodes"`
}

const (
	cloudSqlService = "4bc59b9a-8520-409f-85da-1c7552315863"
	bigTableService = "b8e19880-ac58-42ef-b033-f7cd9c94d1fe"
)

func (*Tile) Configure(cfg *config.Config, om *ops_manager.Sdk) error {
	if err := om.StageProduct(tile.Product); err != nil {
		return err
	}

	network := tiles.NetworkConfig(cfg.ServicesSubnetName, cfg)

	networkBytes, err := json.Marshal(&network)
	if err != nil {
		return err
	}

	properties := Properties{
		ServiceAccountKey: tiles.Value{cfg.ServiceBrokerServiceAccountKey},
		DatabaseHost:      tiles.Value{cfg.ServiceBrokerDbIp},
		DatabaseUsername:  tiles.Value{cfg.ServiceBrokerDbUsername},
		DatabasePassword:  tiles.SecretValue{tiles.Secret{cfg.ServiceBrokerDbPassword}},
		CloudSQLPlans: CloudSQLPlansValue{[]CloudSQLPlan{
			{
				Name:        "mysql-micro-dev",
				DisplayName: "MySQL Micro Development",
				Description: "Micro instance with shared CPU and 0.6 GB of memory for development",
				Service:     cloudSqlService,
				Tier:        "db-f1-micro",
				MaxDiskSize: "100",
				PricingPlan: "PER_USE",
			},
			{
				Name:        "mysql-n1-standard-2",
				DisplayName: "MySQL 2 CPU, 7.5 GB Memory",
				Description: "Instance with 2 dedicated CPUs and 7.5 GB of Memory",
				Service:     cloudSqlService,
				Tier:        "db-f1-micro",
				MaxDiskSize: "1000",
				PricingPlan: "PACKAGE",
			},
		}},
		BigTablePlans: BigTablePlansValue{[]BigTablePlan{
			{
				Name:        "bigtable-micro-dev",
				DisplayName: "BigTable Micro Development",
				Description: "3 nodes count with HDD storage",
				Service:     bigTableService,
				StorageType: "HDD",
				NodeCount:   "3",
			},
			{
				Name:        "bigtable-medium",
				DisplayName: "BigTable Medium Deployment",
				Description: "10 nodes with SDD storage",
				Service:     bigTableService,
				StorageType: "SDD",
				NodeCount:   "10",
			},
		}},
	}

	propertiesBytes, err := json.Marshal(&properties)
	if err != nil {
		return err
	}

	resoruces := "{}"

	return om.ConfigureProduct(tile.Product.Name, string(networkBytes), string(propertiesBytes), resoruces)
}
