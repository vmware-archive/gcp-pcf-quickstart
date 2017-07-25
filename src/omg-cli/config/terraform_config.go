package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type TerraformConfigSchema struct {
	OpsManagerIp                       string `json:"ops_manager_private_ip"`
	JumpboxIp                          string `json:"jumpbox_public_ip"`
	NetworkName                        string `json:"network_name"`
	DeploymentTargetTag                string `json:"vm_tag"`
	MgmtSubnetName                     string `json:"management_subnet_name"`
	MgmtSubnetGateway                  string `json:"management_subnet_gateway"`
	MgmtSubnetCIDR                     string `json:"management_subnet_cidrs_0"`
	ServicesSubnetName                 string `json:"services_subnet_name"`
	ServicesSubnetGateway              string `json:"services_subnet_gateway"`
	ServicesSubnetCIDR                 string `json:"services_subnet_cidrs_0"`
	ErtSubnetName                      string `json:"ert_subnet_name"`
	ErtSubnetGateway                   string `json:"ert_subnet_gateway"`
	ErtSubnetCIDR                      string `json:"ert_subnet_cidrs_0"`
	HttpBackendServiceName             string `json:"http_lb_backend_name"`
	SshTargetPoolName                  string `json:"ssh_router_pool"`
	TcpTargetPoolName                  string `json:"tcp_router_pool"`
	TcpPortRange                       string `json:"tcp_port_range"`
	BuildpacksBucket                   string `json:"buildpacks_bucket"`
	DropletsBucket                     string `json:"droplets_bucket"`
	PackagesBucket                     string `json:"packages_bucket"`
	ResourcesBucket                    string `json:"resources_bucket"`
	DirectorBucket                     string `json:"director_blobstore_bucket"`
	DnsSuffix                          string `json:"dns_suffix"`
	SslCertificate                     string `json:"ssl_cert"`
	SslPrivateKey                      string `json:"ssl_cert_private_key"`
	OpsManServiceAccount               string `json:"service_account_email"`
	ServiceBrokerServiceAccountKey     string `json:"servicebroker_service_account_key"`
	StackdriverNozzleServiceAccountKey string `json:"stackdriver_service_account_key"`

	Region      string `json:"region"`
	Zone1       string `json:"azs_0"`
	Zone2       string `json:"azs_1"`
	Zone3       string `json:"azs_2"`
	ProjectName string `json:"project"`
}

func FromTerraform(filename string) (*Config, error) {
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
