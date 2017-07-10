package config

import (
	"context"
	"encoding/json"
	"fmt"
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

	// Not from the environment:
	OpsManUsername         string
	OpsManPassword         string
	OpsManDecryptionPhrase string
	ProjectName            string
}

func FromEnvironment(ctx context.Context, client *http.Client, projectName string) (*Config, error) {
	cfgMap, err := dumpConfigVariables(ctx, client, fmt.Sprintf("projects/%s/configs/omgConfig", projectName))
	if err != nil {
		return nil, err
	}

	cfg, err := mapToConfig(cfgMap)

	fillInDefaults(cfg, projectName)

	return cfg, err
}

func dumpConfigVariables(ctx context.Context, client *http.Client, configName string) (map[string]string, error) {
	svc, err := runtimeconfig.New(client)

	if err != nil {
		return nil, err
	}

	trimString := len(configName) + len("/variables/")

	cfg := map[string]string{}
	call := svc.Projects.Configs.Variables.List(configName).ReturnValues(true)
	err = call.Pages(ctx, func(res *runtimeconfig.ListVariablesResponse) error {
		for _, v := range res.Variables {
			cfg[v.Name[trimString:len(v.Name)]] = v.Text
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return cfg, err
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

func fillInDefaults(cfg *Config, projectName string) {
	cfg.OpsManUsername = "foo"
	cfg.OpsManPassword = "foobar"
	cfg.OpsManDecryptionPhrase = "foobar"
	cfg.ProjectName = projectName
}
