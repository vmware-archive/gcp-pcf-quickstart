package config

import (
	"encoding/json"
	"net/http"

	runtimeconfig "google.golang.org/api/runtimeconfig/v1beta1"
)

type Config struct {
	OpsManagerIp          string `json:"opsManagerIp"`
	JumpboxName           string `json:"jumpboxName"`
	NetworkName           string `json:"networkName"`
	MgmtSubnetName        string `json:"mgmtSubnetName"`
	MgmtSubnetGateway     string `json:"mgmtSubnetGateway"`
	MgmtSubnetCIDR        string `json:"mgmtSubnetCIDR"`
	ServicesSubnetName    string `json:"servicesSubnetName"`
	ServicesSubnetGateway string `json:"servicesSubnetGateway"`
	ServicesSubnetCIDR    string `json:"servicesSubnetCIDR"`
	ErtSubnetName         string `json:"ertSubnetName"`
	ErtSubnetGateway      string `json:"ertSubnetGateway"`
	ErtSubnetCIDR         string `json:"ertSubnetCIDR"`
	HttpLoadBalancerIP    string `json:"httpLoadBalancerIP"`
	SshTargetPoolName     string `json:"sshTargetPoolName"`
	SshLoadBalancerIP     string `json:"sshLoadBalancerIP"`
	SshTargetTag          string `json:"sshTargetTag"`
	TcpTargetPoolName     string `json:"tcpTargetPoolName"`
	TcpLoadBalancerIP     string `json:"tcpLoadBalancerIP"`
	TcpTargetTag          string `json:"tcpTargetTag"`
	TcpPortRange          string `json:"tcpPortRange"`
	BuildpacksBucket      string `json:"buildpacksBucket"`
	DropletsBucket        string `json:"dropletsBucket"`
	PackagesBucket        string `json:"packagesBucket"`
	ResourcesBucket       string `json:"resourcesBucket"`
	DirectorBucket        string `json:"directorBucket"`
}

func FromEnvironment(client *http.Client, configName string) (*Config, error) {
	cfgMap, err := dumpConfigVariables(client, configName)
	if err != nil {
		return nil, err
	}

	cfg, err := mapToConfig(cfgMap)

	return cfg, err
}

func dumpConfigVariables(client *http.Client, configName string) (map[string]string, error) {
	svc, err := runtimeconfig.New(client)
	if err != nil {
		return nil, err
	}

	list, err := svc.Projects.Configs.Variables.List(configName).Do()
	if err != nil {
		return nil, err
	}

	cfg := map[string]string{}
	trimString := len(configName) + len("/variables/")

	for _, v := range list.Variables {
		v, err := svc.Projects.Configs.Variables.Get(v.Name).Do()
		if err != nil {
			return nil, err
		}
		cfg[v.Name[trimString:len(v.Name)]] = v.Text
	}

	return cfg, nil
}

func mapToConfig(cfgMap map[string]string) (*Config, error) {
	str, err := json.Marshal(cfgMap)

	if err != nil {
		return nil, err
	}

	hydratedCfg := &Config{}
	err = json.Unmarshal(str, hydratedCfg)

	return hydratedCfg, err
}
