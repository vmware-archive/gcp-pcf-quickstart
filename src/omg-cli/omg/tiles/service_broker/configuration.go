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
	ServiceAccountKey Value       `json:".properties.root_service_account_json"`
	DatabaseHost      Value       `json:".properties.db_host"`
	DatabaseUsername  Value       `json:".properties.db_username"`
	DatabasePassword  SecretValue `json:".properties.db_password"`

	CloudSQLPlans CloudSQLPlanValue `json:".properties.cloudsql_custom_plans"`
}

type CloudSQLPlanValue struct {
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

const cloudSqlService = "4bc59b9a-8520-409f-85da-1c7552315863"

type Value struct {
	Value string `json:"value"`
}

type Secret struct {
	Value string `json:"secret"`
}

type SecretValue struct {
	Sec Secret `json:"value"`
}

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
		ServiceAccountKey: Value{cfg.ServiceBrokerServiceAccountKey},
		DatabaseHost:      Value{cfg.ServiceBrokerDbIp},
		DatabaseUsername:  Value{cfg.ServiceBrokerDbUsername},
		DatabasePassword:  SecretValue{Secret{cfg.ServiceBrokerDbPassword}},
		CloudSQLPlans: CloudSQLPlanValue{[]CloudSQLPlan{
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
	}

	propertiesBytes, err := json.Marshal(&properties)
	if err != nil {
		return err
	}

	resoruces := "{}"

	return om.ConfigureProduct(tile.Product.Name, string(networkBytes), string(propertiesBytes), resoruces)
}
