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
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
)

type TerraformConfigSchema struct {
	OpsManagerHostname  string `json:"ops_manager_dns"`
	OpsManagerIp        string `json:"ops_manager_private_ip"`
	JumpboxIp           string `json:"jumpbox_public_ip"`
	NetworkName         string `json:"network_name"`
	DeploymentTargetTag string `json:"vm_tag"`

	OpsManagerServiceAccountKey string `json:"ops_manager_service_account_key"`

	ExternalSqlIp         string `json:"sql_db_ip"`
	ExternalSqlPort       int
	OpsManagerSqlDbName   string `json:"opsman_sql_db_name"`
	OpsManagerSqlUsername string `json:"opsman_sql_username"`
	OpsManagerSqlPassword string `json:"opsman_sql_password"`
	ERTSqlDbName          string `json:"ert_sql_db_name"`
	ERTSqlUsername        string `json:"ert_sql_username"`
	ERTSqlPassword        string `json:"ert_sql_password"`

	MgmtSubnetName    string `json:"management_subnet_name"`
	MgmtSubnetGateway string `json:"management_subnet_gateway"`
	MgmtSubnetCIDR    string `json:"management_subnet_cidrs_0"`

	ServicesSubnetName    string `json:"services_subnet_name"`
	ServicesSubnetGateway string `json:"services_subnet_gateway"`
	ServicesSubnetCIDR    string `json:"services_subnet_cidrs_0"`

	DynamicServicesSubnetName    string `json:"dynamic_services_subnet_name"`
	DynamicServicesSubnetGateway string `json:"dynamic_services_subnet_gateway"`
	DynamicServicesSubnetCIDR    string `json:"dynamic_services_subnet_cidrs_0"`

	ErtSubnetName    string `json:"ert_subnet_name"`
	ErtSubnetGateway string `json:"ert_subnet_gateway"`
	ErtSubnetCIDR    string `json:"ert_subnet_cidrs_0"`

	HttpBackendServiceName string `json:"http_lb_backend_name"`
	SshTargetPoolName      string `json:"ssh_router_pool"`
	WssTargetPoolName      string `json:"wss_router_pool"`
	TcpTargetPoolName      string `json:"tcp_router_pool"`
	TcpPortRange           string `json:"tcp_port_range"`

	BuildpacksBucket string `json:"buildpacks_bucket"`
	DropletsBucket   string `json:"droplets_bucket"`
	PackagesBucket   string `json:"packages_bucket"`
	ResourcesBucket  string `json:"resources_bucket"`
	DirectorBucket   string `json:"director_blobstore_bucket"`

	DnsSuffix         string `json:"dns_suffix"`
	AppsDomain        string `json:"apps_domain"`
	SysDomain         string `json:"sys_domain"`
	DopplerDomain     string `json:"doppler_domain"`
	LoggregatorDomain string `json:"loggregator_domain"`

	SslCertificate string `json:"ssl_cert"`
	SslPrivateKey  string `json:"ssl_cert_private_key"`

	StackdriverNozzleServiceAccountKey string `json:"stackdriver_service_account_key"`

	ServiceBrokerServiceAccountKey string `json:"service_broker_service_account_key"`
	ServiceBrokerDbIp              string `json:"service_broker_db_ip"`
	ServiceBrokerDbUsername        string `json:"service_broker_db_username"`
	ServiceBrokerDbPassword        string `json:"service_broker_db_password"`

	Region      string `json:"region"`
	Zone1       string `json:"azs_0"`
	Zone2       string `json:"azs_1"`
	Zone3       string `json:"azs_2"`
	ProjectName string `json:"project"`

	OpsManager OpsManagerCredentials
}

func TerraformFromEnvDirectory(path string) (*Config, error) {
	return fromTerraform(filepath.Join(path, TerraformOutputFile))
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

	if flattened["ops_manager_skip_ssl_verify"] == "true" {
		hydratedCfg.OpsManager.SkipSSLVerification = true
	}

	if val := flattened["sql_db_port"]; val != "" {
		parsed, err := strconv.ParseInt(val, 10, 0)
		if err != nil {
			return nil, err
		}
		hydratedCfg.ExternalSqlPort = int(parsed)
	}

	hydratedCfg.OpsManager.Username = flattened["ops_manager_username"]
	hydratedCfg.OpsManager.Password = flattened["ops_manager_password"]
	hydratedCfg.OpsManager.DecryptionPhrase = flattened["ops_manager_decryption_phrase"]

	cfg := Config(hydratedCfg)

	return &cfg, nil
}

/*
 * translate:
 * { "foo": {"value": "bar"}, "baz": {"value": ["pizza", "thebest"]}}
 * to:
 * {"foo": "bar", "baz_0": "pizza", "baz_1": "thebest"}
 */

type TerraformValue struct {
	Value interface{} `json:"value"`
}

func flattenTerraform(contents []byte) (map[string]string, error) {
	res := map[string]string{}

	tf := map[string]TerraformValue{}

	err := json.Unmarshal(contents, &tf)
	if err != nil {
		return nil, err
	}

	for k, v := range tf {
		if str, ok := v.Value.(string); ok {
			res[k] = str
		} else if arr, ok := v.Value.([]interface{}); ok {
			for i, entry := range arr {
				res[fmt.Sprintf("%s_%d", k, i)] = entry.(string)
			}
		} else {
			return nil, fmt.Errorf("encountered unknown type in terraform config: %v", v.Value)
		}
	}

	return res, nil
}
