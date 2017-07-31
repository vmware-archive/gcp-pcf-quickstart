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
}

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
	}

	propertiesBytes, err := json.Marshal(&properties)
	if err != nil {
		return err
	}

	resoruces := "{}"

	return om.ConfigureProduct(tile.Product.Name, string(networkBytes), string(propertiesBytes), resoruces)
}
