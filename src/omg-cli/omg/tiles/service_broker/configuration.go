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

const (
	postgresService = "cbad6d78-a73c-432d-b8ff-b219a17a803a"
	mySqlService    = "4bc59b9a-8520-409f-85da-1c7552315863"
	bigTableService = "b8e19880-ac58-42ef-b033-f7cd9c94d1fe"
	spannerService  = "51b3e27e-d323-49ce-8c5f-1211e6409e82"
)

type Properties struct {
	ServiceAccountKey tiles.Value       `json:".properties.root_service_account_json"`
	DatabaseHost      tiles.Value       `json:".properties.db_host"`
	DatabaseUsername  tiles.Value       `json:".properties.db_username"`
	DatabasePassword  tiles.SecretValue `json:".properties.db_password"`

	CloudSQLMySQLPlans    CloudSQLPlansValue `json:".properties.cloudsql_mysql_custom_plans"`
	CloudSQLPostgresPlans CloudSQLPlansValue `json:".properties.cloudsql_postgres_custom_plans"`
	BigTablePlans         BigTablePlansValue `json:".properties.bigtable_custom_plans"`
	SpannerPlans          SpannerPlansValue  `json:".properties.spanner_custom_plans"`
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

type SpannerPlansValue struct {
	Plans []SpannerPlan `json:"value"`
}

type SpannerPlan struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Service     string `json:"service"`
	NodeCount   string `json:"num_nodes"`
}

func (*Tile) Configure(envConfig *config.EnvConfig, cfg *config.Config, om *ops_manager.Sdk) error {
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
		CloudSQLPostgresPlans: CloudSQLPlansValue{[]CloudSQLPlan{
			{
				Name:        "postgres-micro-dev",
				DisplayName: "PostgreSQL Micro (Development)",
				Description: "Micro instance with shared CPU, 0.6 GB of memory",
				Service:     postgresService,
				Tier:        "db-f1-micro",
				MaxDiskSize: "100",
				PricingPlan: "PER_USE",
			},
			{
				Name:        "postgres-n1-standard-2",
				DisplayName: "PostgreSQL 2 CPU, 7.5 GB Memory",
				Description: "Instance with 2 dedicated CPUs, 7.5 GB of Memory",
				Service:     postgresService,
				Tier:        "db-n1-standard-2",
				MaxDiskSize: "10000",
				PricingPlan: "PER_USE",
			},
			{
				Name:        "postgres-n1-standard-16",
				DisplayName: "PostgreSQL 16 CPU, 60 GB Memory",
				Description: "Instance with 16 dedicated CPUs, 60 GB of Memory",
				Service:     postgresService,
				Tier:        "db-n1-standard-16",
				MaxDiskSize: "10000",
				PricingPlan: "PER_USE",
			},
			{
				Name:        "postgres-n1-highmem-16",
				DisplayName: "PostgreSQL 16 CPU, 60 GB Memory",
				Description: "Instance with 16 dedicated CPUs, 104 GB of Memory",
				Service:     postgresService,
				Tier:        "db-n1-highmem-16",
				MaxDiskSize: "10000",
				PricingPlan: "PER_USE",
			},
		}},
		CloudSQLMySQLPlans: CloudSQLPlansValue{[]CloudSQLPlan{
			{
				Name:        "mysql-micro-dev",
				DisplayName: "MySQL Micro (Development)",
				Description: "Micro instance with shared CPU, 0.6 GB of memory",
				Service:     mySqlService,
				Tier:        "db-f1-micro",
				MaxDiskSize: "100",
				PricingPlan: "PER_USE",
			},
			{
				Name:        "mysql-n1-standard-2",
				DisplayName: "MySQL 2 CPU, 7.5 GB Memory",
				Description: "Instance with 2 dedicated CPUs, 7.5 GB of Memory",
				Service:     mySqlService,
				Tier: "db-n1-standard-2",
				MaxDiskSize: "10000",
				PricingPlan: "PACKAGE",
			},
			{
				Name:        "mysql-n1-standard-16",
				DisplayName: "MySQL 16 CPU, 60 GB Memory",
				Description: "Instance with 16 dedicated CPUs, 60 GB of Memory",
				Service:     mySqlService,
				Tier:        "db-n1-standard-16",
				MaxDiskSize: "10000",
				PricingPlan: "PACKAGE",
			},
			{
				Name:        "mysql-n1-highmem-16",
				DisplayName: "MySQL 16 CPU, 60 GB Memory",
				Description: "Instance with 16 dedicated CPUs, 104 GB of Memory",
				Service:     mySqlService,
				Tier:        "db-n1-highmem-16",
				MaxDiskSize: "10000",
				PricingPlan: "PACKAGE",
			},
		}},
		BigTablePlans: BigTablePlansValue{[]BigTablePlan{
			{
				Name:        "bigtable-hdd-10",
				DisplayName: "BigTable Micro (Development)",
				Description: "3 nodes count with HDD storage",
				Service:     bigTableService,
				StorageType: "HDD",
				NodeCount:   "3",
			},
			{
				Name:        "bigtable-ssd-10",
				DisplayName: "BigTable 10 Nodes",
				Description: "10 nodes with SDD storage",
				Service:     bigTableService,
				StorageType: "SDD",
				NodeCount:   "10",
			},
			{
				Name:        "bigtable-ssd-20",
				DisplayName: "BigTable 20 Nodes",
				Description: "20 nodes with SDD storage",
				Service:     bigTableService,
				StorageType: "SDD",
				NodeCount:   "20",
			},
			{
				Name:        "bigtable-ssd-30",
				DisplayName: "BigTable 30 Nodes",
				Description: "30 nodes with SDD storage",
				Service:     bigTableService,
				StorageType: "SDD",
				NodeCount:   "30",
			},
		}},
		SpannerPlans: SpannerPlansValue{[]SpannerPlan{
			{
				Name:        "spanner-regional-micro-dev",
				DisplayName: "Spanner Micro (Development)",
				Description: "Spanner Instance with 1 node",
				Service:     spannerService,
				NodeCount:   "1",
			},
			{
				Name:        "spanner-regional-3",
				DisplayName: "Spanner 3 Node",
				Description: "Spanner Instance with 3 nodes, minimum recommendation for production",
				Service:     spannerService,
				NodeCount:   "3",
			},
			{
				Name:        "spanner-regional-10",
				DisplayName: "Spanner 10 Node",
				Description: "Spanner Instance with 10 nodes",
				Service:     spannerService,
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
