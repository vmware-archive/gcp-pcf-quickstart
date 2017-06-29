package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/Jeffail/gabs"
	"github.com/aditya87/hummus"
	"github.com/pivotal-cf/pcf-releng-ci/tasks/om/generate-bosh-configuration/subnet"
)

var (
	provider              string
	providerConfiguration string
	envName               string
	awsAccessKeyID        string
	awsSecretAccessKey    string
	compilationVMType     string
)

type Config struct {
	/**IaaS configuration**/
	//GCP
	Project              string `json:"project" hummus:"iaas_configuration.project,omitempty"`
	DefaultDeploymentTag string `hummus:"iaas_configuration.default_deployment_tag,omitempty"`
	AuthJSON             string `json:"service_account_key" hummus:"iaas_configuration.auth_json,omitempty"`

	//Azure
	SubscriptionID                string `json:"subscription_id" hummus:"iaas_configuration.subscription_id,omitempty"`
	TenantID                      string `json:"tenant_id" hummus:"iaas_configuration.tenant_id,omitempty"`
	ClientID                      string `json:"client_id" hummus:"iaas_configuration.client_id,omitempty"`
	ClientSecret                  string `json:"client_secret" hummus:"iaas_configuration.client_secret,omitempty"`
	ResourceGroupName             string `json:"pcf_resource_group_name" hummus:"iaas_configuration.resource_group_name,omitempty"`
	BoshStorageAccountName        string `json:"bosh_root_storage_account" hummus:"iaas_configuration.bosh_storage_account_name,omitempty"`
	DefaultSecurityGroup          string `json:"ops_manager_security_group_name" hummus:"iaas_configuration.default_security_group,omitempty"`
	SSHPublicKey                  string `json:"ops_manager_ssh_public_key" hummus:"iaas_configuration.ssh_public_key,omitempty"`
	DeploymentsStorageAccountName string `json:"wildcard_vm_storage_account" hummus:"iaas_configuration.deployments_storage_account_name,omitempty"`

	//AWS
	AccessKeyID     string `json:"iam_user_access_key" hummus:"iaas_configuration.access_key_id,omitempty"`
	SecretAccessKey string `json:"iam_user_secret_access_key" hummus:"iaas_configuration.secret_access_key,omitempty"`
	VPCID           string `json:"vpc_id" hummus:"iaas_configuration.vpc_id,omitempty"`
	SecurityGroupID string `json:"vms_security_group_id" hummus:"iaas_configuration.security_group,omitempty"`
	SSHKeyPairName  string `json:"ops_manager_ssh_public_key_name" hummus:"iaas_configuration.key_pair_name,omitempty"`

	// Used for AWS, Azure and GCP
	SSHPrivateKey string `json:"ops_manager_ssh_private_key" hummus:"iaas_configuration.ssh_private_key,omitempty"`

	/**Director configuration**/
	NTPServers string `hummus:"director_configuration.ntp_servers_string,omitempty"`

	/**AZ configuration**/
	AZs           []string `json:"azs"`
	Clusters      []string `json:"clusters"`
	ResourcePools []string `json:"resource_pools"`

	/**Networks configuration**/
	Region  string `json:"region" hummus:"iaas_configuration.region,omitempty"`
	Network string `json:"network_name"`

	ICMPChecksEnabled bool `hummus:"networks_configuration.icmp_checks_enabled"`

	ManagementSubnetIDList   []string `json:"management_subnet_ids"`
	ManagementSubnetCIDRList []string `json:"management_subnet_cidrs"`
	ManagementSubnetAZList   []string `json:"management_subnet_availability_zones"`

	ERTSubnetIDList   []string `json:"ert_subnet_ids"`
	ERTSubnetCIDRList []string `json:"ert_subnet_cidrs"`
	ERTSubnetAZList   []string `json:"ert_subnet_availability_zones"`

	ServicesSubnetIDList   []string `json:"services_subnet_ids"`
	ServicesSubnetCIDRList []string `json:"services_subnet_cidrs"`
	ServicesSubnetAZList   []string `json:"services_subnet_availability_zones"`

	ManagementNetworkName    string   `json:"management_subnet_name" hummus:"networks_configuration.networks[0].name,omitempty"`
	ManagementServiceNetwork bool     `hummus:"networks_configuration.networks[0].service_network"`
	ManagementSubnets        []Subnet `hummus:"networks_configuration.networks[0].subnets"`
	ManagementSubnetGateway  string   `json:"management_subnet_gateway" hummus:",omitempty"`

	ERTNetworkName    string   `json:"ert_subnet_name" hummus:"networks_configuration.networks[1].name,omitempty"`
	ERTServiceNetwork bool     `hummus:"networks_configuration.networks[1].service_network"`
	ERTSubnets        []Subnet `hummus:"networks_configuration.networks[1].subnets"`
	ERTSubnetGateway  string   `json:"ert_subnet_gateway" hummus:",omitempty"`

	ServicesNetworkName    string   `json:"services_subnet_name" hummus:"networks_configuration.networks[2].name,omitempty"`
	ServicesServiceNetwork bool     `hummus:"networks_configuration.networks[2].service_network"`
	ServicesSubnets        []Subnet `hummus:"networks_configuration.networks[2].subnets"`
	ServicesSubnetGateway  string   `json:"services_subnet_gateway" hummus:",omitempty"`

	/**Network/AZ assignment**/
	SingletonAZ     string `hummus:"network_assignment.singleton_availability_zone,omitempty"`
	AssignedNetwork string `hummus:"network_assignment.network,omitempty"`

	/**Resource configuration**/
	ResourceConfiguration string `hummus:"resource_configuration.compilation.instance_type.id,omitempty"`
}

type Subnet struct {
	IaasIdentifier string   `json:"iaas_identifier"`
	CIDR           string   `json:"cidr"`
	ReservedIPs    string   `json:"reserved_ip_ranges"`
	DNS            string   `json:"dns"`
	Gateway        string   `json:"gateway"`
	AZs            []string `json:"availability_zones,omitempty"`
}

type subnetConfig struct {
	subnetIDs        []string
	cidrs            []string
	azs              []string
	reservedRangeMax int
	dns              string
	customGateway    string
	crossZone        bool
}

func main() {
	flag.StringVar(&provider, "provider", "", "provider")
	flag.StringVar(&providerConfiguration, "provider-configuration", "", "provider-specific configuration")
	flag.StringVar(&envName, "env-name", "", "environment name")
	flag.StringVar(&awsAccessKeyID, "aws-access-key-id", "", "AWS access key ID")
	flag.StringVar(&awsSecretAccessKey, "aws-secret-access-key", "", "AWS secret access key")
	flag.StringVar(&compilationVMType, "compilation-vm-type", "", "Optional BOSH compilation VM type")
	flag.Parse()

	var config Config
	err := json.Unmarshal([]byte(providerConfiguration), &config)
	if err != nil {
		log.Fatalln(err)
	}

	var finalJSON []byte
	switch provider {
	case "aws":
		finalJSON, err = configureAWS(config, envName)
	case "azure":
		finalJSON, err = configureAzure(config, envName)
	case "gcp":
		finalJSON, err = configureGCP(config, envName)
	default:
		log.Fatalf("invalid provider: %s", provider)
	}

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(string(finalJSON)))
}

func configureAWS(config Config, envName string) ([]byte, error) {
	config.ICMPChecksEnabled = false
	config.NTPServers = "0.amazon.pool.ntp.org"

	config.ManagementNetworkName = fmt.Sprintf("%s-management-network", envName)
	config.ManagementServiceNetwork = false
	managementSubnets, err := generateSubnets(subnetConfig{
		subnetIDs:        config.ManagementSubnetIDList,
		cidrs:            config.ManagementSubnetCIDRList,
		azs:              config.ManagementSubnetAZList,
		reservedRangeMax: 4,
		dns:              "169.254.169.253",
		customGateway:    "",
	})
	if err != nil {
		return []byte{}, err
	}
	config.ManagementSubnets = managementSubnets

	config.ERTNetworkName = fmt.Sprintf("%s-ert-network", envName)
	config.ERTServiceNetwork = false
	ertSubnets, err := generateSubnets(subnetConfig{
		subnetIDs:        config.ERTSubnetIDList,
		cidrs:            config.ERTSubnetCIDRList,
		azs:              config.ERTSubnetAZList,
		reservedRangeMax: 4,
		dns:              "169.254.169.253",
		customGateway:    "",
	})
	if err != nil {
		return []byte{}, err
	}
	config.ERTSubnets = ertSubnets

	config.ServicesNetworkName = fmt.Sprintf("%s-services-network", envName)
	config.ServicesServiceNetwork = true
	servicesSubnets, err := generateSubnets(subnetConfig{
		subnetIDs:        config.ServicesSubnetIDList,
		cidrs:            config.ServicesSubnetCIDRList,
		azs:              config.ServicesSubnetAZList,
		reservedRangeMax: 3,
		dns:              "169.254.169.253",
		customGateway:    "",
	})
	if err != nil {
		return []byte{}, err
	}
	config.ServicesSubnets = servicesSubnets

	config.AssignedNetwork = config.ManagementNetworkName
	config.SingletonAZ = config.AZs[0]
	config.ResourceConfiguration = compilationVMType

	var finalJSON []byte
	finalJSON, err = hummus.Marshal(config)
	if err != nil {
		return []byte{}, err
	}

	parsedJSON, err := gabs.ParseJSON(finalJSON)
	if err != nil {
		return []byte{}, err
	}

	// Have to inject this outside of hummus because we can't use omitempty to prevent it from being included in the config for other IaaS.
	parsedJSON.SetP(false, "iaas_configuration.encrypted")

	parsedJSON.Array("az_configuration", "availability_zones")

	for _, azName := range config.AZs {
		azMap := make(map[string]string)

		azMap["name"] = azName

		parsedJSON.ArrayAppend(azMap, "az_configuration", "availability_zones")
	}

	return parsedJSON.Bytes(), nil
}

func configureAzure(config Config, envName string) ([]byte, error) {
	managementSubnets, err := generateSubnets(subnetConfig{
		subnetIDs:        []string{fmt.Sprintf("%s/%s", config.Network, config.ManagementNetworkName)},
		cidrs:            config.ManagementSubnetCIDRList,
		azs:              config.ManagementSubnetAZList,
		reservedRangeMax: 5,
		dns:              "8.8.8.8",
		customGateway:    config.ManagementSubnetGateway,
	})
	if err != nil {
		return []byte{}, err
	}
	config.ManagementSubnets = managementSubnets

	ertSubnets, err := generateSubnets(subnetConfig{
		subnetIDs:        []string{fmt.Sprintf("%s/%s", config.Network, config.ERTNetworkName)},
		cidrs:            config.ERTSubnetCIDRList,
		azs:              config.ERTSubnetAZList,
		reservedRangeMax: 4,
		dns:              "8.8.8.8",
		customGateway:    config.ERTSubnetGateway,
	})
	if err != nil {
		return []byte{}, err
	}
	config.ERTSubnets = ertSubnets

	servicesSubnets, err := generateSubnets(subnetConfig{
		subnetIDs:        []string{fmt.Sprintf("%s/%s", config.Network, config.ServicesNetworkName)},
		cidrs:            config.ServicesSubnetCIDRList,
		azs:              config.ServicesSubnetAZList,
		reservedRangeMax: 3,
		dns:              "8.8.8.8",
		customGateway:    config.ServicesSubnetGateway,
	})
	if err != nil {
		return []byte{}, err
	}
	config.ServicesSubnets = servicesSubnets

	// Empty these so they don't show up in the final config
	config.ManagementSubnetGateway, config.ERTSubnetGateway, config.ServicesSubnetGateway = "", "", ""
	config.NTPServers = "us.pool.ntp.org"
	config.ICMPChecksEnabled = false
	config.ManagementServiceNetwork = false
	config.ERTServiceNetwork = false
	config.ServicesServiceNetwork = true
	config.AssignedNetwork = config.ManagementNetworkName
	config.ResourceConfiguration = compilationVMType

	var finalJSON []byte
	finalJSON, err = hummus.Marshal(config)
	if err != nil {
		return []byte{}, err
	}

	return finalJSON, nil
}

func configureGCP(config Config, envName string) ([]byte, error) {
	mgmtIaasIdentifier := fmt.Sprintf("%s/%s/%s", config.Network, config.ManagementNetworkName, config.Region)
	managementSubnets, err := generateSubnets(subnetConfig{
		subnetIDs:        []string{mgmtIaasIdentifier},
		cidrs:            config.ManagementSubnetCIDRList,
		azs:              config.AZs,
		reservedRangeMax: 4,
		dns:              "8.8.8.8",
		customGateway:    config.ManagementSubnetGateway,
		crossZone:        true,
	})
	if err != nil {
		return []byte{}, err
	}
	config.ManagementSubnets = managementSubnets

	ertSubnets, err := generateSubnets(subnetConfig{
		subnetIDs:        []string{fmt.Sprintf("%s/%s/%s", config.Network, config.ERTNetworkName, config.Region)},
		cidrs:            config.ERTSubnetCIDRList,
		azs:              config.AZs,
		reservedRangeMax: 4,
		dns:              "8.8.8.8",
		customGateway:    config.ERTSubnetGateway,
		crossZone:        true,
	})
	if err != nil {
		return []byte{}, err
	}
	config.ERTSubnets = ertSubnets

	servicesSubnets, err := generateSubnets(subnetConfig{
		subnetIDs:        []string{fmt.Sprintf("%s/%s/%s", config.Network, config.ServicesNetworkName, config.Region)},
		cidrs:            config.ServicesSubnetCIDRList,
		azs:              config.AZs,
		reservedRangeMax: 3,
		dns:              "8.8.8.8",
		customGateway:    config.ServicesSubnetGateway,
		crossZone:        true,
	})
	if err != nil {
		return []byte{}, err
	}
	config.ServicesSubnets = servicesSubnets

	// Empty these so they don't show up in the final config
	config.ManagementSubnetGateway, config.ERTSubnetGateway, config.ServicesSubnetGateway = "", "", ""
	config.Region = ""

	config.DefaultDeploymentTag = fmt.Sprintf("%s-vms", envName)
	config.NTPServers = "169.254.169.254"
	config.ICMPChecksEnabled = false
	config.ManagementServiceNetwork = false
	config.ERTServiceNetwork = false
	config.ServicesServiceNetwork = true
	config.AssignedNetwork = config.ManagementNetworkName
	config.SingletonAZ = config.AZs[0]

	config.ResourceConfiguration = compilationVMType

	var finalJSON []byte
	finalJSON, err = hummus.Marshal(config)
	if err != nil {
		return []byte{}, err
	}

	parsedJSON, err := gabs.ParseJSON(finalJSON)
	if err != nil {
		return []byte{}, err
	}

	parsedJSON.Array("az_configuration", "availability_zones")

	for _, azName := range config.AZs {
		azMap := make(map[string]string)
		azMap["name"] = azName
		parsedJSON.ArrayAppend(azMap, "az_configuration", "availability_zones")
	}

	return parsedJSON.Bytes(), nil
}

func generateSubnets(config subnetConfig) ([]Subnet, error) {
	var subnets []Subnet

	for index, cidr := range config.cidrs {
		subnetCIDR, err := subnet.ParseSubnet(cidr)
		if err != nil {
			return []Subnet{}, err
		}

		reservedIPs, err := subnetCIDR.Range(0, config.reservedRangeMax)
		if err != nil {
			return []Subnet{}, err
		}

		var gateway string
		if config.customGateway != "" {
			gateway = config.customGateway
		} else {
			gateway, err = subnetCIDR.IPAddress(1)
			if err != nil {
				return []Subnet{}, err
			}
		}

		var subnetAZs []string
		if config.crossZone {
			subnetAZs = config.azs
		} else {
			if index <= len(config.azs)-1 {
				subnetAZs = append(subnetAZs, config.azs[index])
			}
		}

		subnets = append(subnets, Subnet{
			IaasIdentifier: config.subnetIDs[index],
			CIDR:           cidr,
			DNS:            config.dns,
			AZs:            subnetAZs,
			Gateway:        gateway,
			ReservedIPs:    reservedIPs,
		})
	}

	return subnets, nil
}
