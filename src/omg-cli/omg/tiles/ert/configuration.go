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

package ert

import (
	"encoding/json"
	"fmt"
	"omg-cli/config"
	"omg-cli/omg/tiles"
	"omg-cli/ops_manager"
)

type Properties struct {
	// Domains
	AppsDomain                tiles.Value              `json:".cloud_controller.apps_domain"`
	SysDomain                 tiles.Value              `json:".cloud_controller.system_domain"`
	TcpRouting                tiles.Value              `json:".properties.tcp_routing"`
	TcpRoutingReservablePorts tiles.Value              `json:".properties.tcp_routing.enable.reservable_ports"`
	GoRouterSSLCiphers        tiles.Value              `json:".properties.gorouter_ssl_ciphers"`
	HAProxySSLCiphers         tiles.Value              `json:".properties.haproxy_ssl_ciphers"`
	SkipSSLVerification       tiles.BooleanValue       `json:".ha_proxy.skip_cert_verify"`
	HAProxyForwardTLS         tiles.Value              `json:".properties.haproxy_forward_tls"`
	IngressCertificates       tiles.CertificateValue   `json:".properties.networking_poe_ssl_certs"`
	CredhubEncryptionKey      tiles.EncryptionKeyValue `json:".properties.credhub_key_encryption_passwords"`
	// Application Security Groups
	SecurityAcknowledgement tiles.Value `json:".properties.security_acknowledgement"`
	// UAA
	ServiceProviderCredentials tiles.OldCertificateValue `json:".uaa.service_provider_key_credentials"`

	UaaDbChoice tiles.Value `json:".properties.uaa_database"`

	// Databases
	ErtDbChoice tiles.Value `json:".properties.system_database"`

	// MySQL
	MySqlMonitorRecipientEmail tiles.Value `json:".mysql_monitor.recipient_email"`

	// Credhub
	CredhubDbChoice tiles.Value `json:".properties.credhub_database"`
}

type LargeFootprintResources struct {
	TcpRouter                    tiles.Resource `json:"tcp_router"`
	Router                       tiles.Resource `json:"router"`
	DiegoBrain                   tiles.Resource `json:"diego_brain"`
	ConsulServer                 tiles.Resource `json:"consul_server"`
	Nats                         tiles.Resource `json:"nats"`
	NfsServer                    tiles.Resource `json:"nfs_server"`
	MysqlProxy                   tiles.Resource `json:"mysql_proxy"`
	Mysql                        tiles.Resource `json:"mysql"`
	BackupPrepare                tiles.Resource `json:"backup_restore"`
	DiegoDatabase                tiles.Resource `json:"diego_database"`
	Uaa                          tiles.Resource `json:"uaa"`
	CloudController              tiles.Resource `json:"cloud_controller"`
	HaProxy                      tiles.Resource `json:"ha_proxy"`
	MysqlMonitor                 tiles.Resource `json:"mysql_monitor"`
	ClockGlobal                  tiles.Resource `json:"clock_global"`
	CloudControllerWorker        tiles.Resource `json:"cloud_controller_worker"`
	DiegoCell                    tiles.Resource `json:"diego_cell"`
	LoggregatorTrafficcontroller tiles.Resource `json:"loggregator_trafficcontroller"`
	SyslogAdapter                tiles.Resource `json:"syslog_adapter"`
	SyslogScheduler              tiles.Resource `json:"syslog_scheduler"`
	Doppler                      tiles.Resource `json:"doppler"`
}

type SmallFootprintResources struct {
	TcpRouter tiles.Resource `json:"tcp_router"`
	Router    tiles.Resource `json:"router"`

	Database    tiles.Resource `json:"database"`
	Control     tiles.Resource `json:"control"`
	Compute     tiles.Resource `json:"compute"`
	FileStorage tiles.Resource `json:"blobstore"`

	HaProxy       tiles.Resource `json:"ha_proxy"`
	BackupPrepare tiles.Resource `json:"backup_restore"`
	MysqlMonitor  tiles.Resource `json:"mysql_monitor"`
}

func (*Tile) Configure(envConfig *config.EnvConfig, cfg *config.Config, om *ops_manager.Sdk) error {
	if err := om.StageProduct(product); err != nil {
		return err
	}

	network := tiles.NetworkConfig(cfg.ErtSubnetName, cfg)

	networkBytes, err := json.Marshal(&network)
	if err != nil {
		return err
	}

	properties := Properties{
		AppsDomain:          tiles.Value{cfg.AppsDomain},
		SysDomain:           tiles.Value{cfg.SysDomain},
		SkipSSLVerification: tiles.BooleanValue{true},
		HAProxyForwardTLS:   tiles.Value{"disable"},
		IngressCertificates: tiles.CertificateValue{[]tiles.CertificateConstruct{
			{Certificate: tiles.Certificate{cfg.SslCertificate, cfg.SslPrivateKey},
				Name: "Certificate",
			},
		},
		},
		CredhubEncryptionKey: tiles.EncryptionKeyValue{[]tiles.EncryptionKey{
			{
				Name:    cfg.CredhubKey.Name,
				Key:     tiles.KeyStruct{Secret: cfg.CredhubKey.Key},
				Primary: true,
			},
		},
		},
		TcpRouting:                 tiles.Value{"enable"},
		TcpRoutingReservablePorts:  tiles.Value{cfg.TcpPortRange},
		GoRouterSSLCiphers:         tiles.Value{"ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384"},
		HAProxySSLCiphers:          tiles.Value{"DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384"},
		SecurityAcknowledgement:    tiles.Value{"X"},
		ServiceProviderCredentials: tiles.OldCertificateValue{tiles.Certificate{cfg.SslCertificate, cfg.SslPrivateKey}},

		UaaDbChoice:     tiles.Value{"internal_mysql"},
		ErtDbChoice:     tiles.Value{"internal_pxc"},
		CredhubDbChoice: tiles.Value{"internal_mysql"},

		MySqlMonitorRecipientEmail: tiles.Value{"admin@example.org"},
	}

	propertiesBytes, err := json.Marshal(&properties)
	if err != nil {
		return err
	}

	resourcesBytes := []byte{}

	zero := 0
	one := 1
	three := 3
	if envConfig.SmallFootprint {
		resources := SmallFootprintResources{
			TcpRouter: tiles.Resource{
				RouterNames:       []string{fmt.Sprintf("tcp:%s", cfg.TcpTargetPoolName)},
				InternetConnected: false,
				Instances:         &one,
			},
			Router: tiles.Resource{
				RouterNames:       []string{fmt.Sprintf("tcp:%s", cfg.WssTargetPoolName), fmt.Sprintf("http:%s", cfg.HttpBackendServiceName)},
				InternetConnected: false,
			},
			Control: tiles.Resource{
				RouterNames:       []string{fmt.Sprintf("tcp:%s", cfg.SshTargetPoolName)},
				InternetConnected: false,
			},
			HaProxy:      tiles.Resource{Instances: &zero},
			MysqlMonitor: tiles.Resource{Instances: &one},
		}
		// Healthwatch pushes quite a few apps, make sure we have enough compute
		if envConfig.IncludeHealthwatch {
			resources.Compute = tiles.Resource{
				Instances: &three,
			}
		}
		resourcesBytes, err = json.Marshal(&resources)
	} else {
		resources := LargeFootprintResources{
			TcpRouter: tiles.Resource{
				RouterNames:       []string{fmt.Sprintf("tcp:%s", cfg.TcpTargetPoolName)},
				InternetConnected: false,
				Instances:         &three,
			},
			Router: tiles.Resource{
				RouterNames:       []string{fmt.Sprintf("tcp:%s", cfg.WssTargetPoolName), fmt.Sprintf("http:%s", cfg.HttpBackendServiceName)},
				InternetConnected: false,
			},
			DiegoBrain: tiles.Resource{
				RouterNames:       []string{fmt.Sprintf("tcp:%s", cfg.SshTargetPoolName)},
				InternetConnected: false,
			},
			HaProxy:      tiles.Resource{Instances: &zero},
			MysqlProxy:   tiles.Resource{Instances: &one},
			Mysql:        tiles.Resource{Instances: &one},
			MysqlMonitor: tiles.Resource{Instances: &one},
		}
		resourcesBytes, err = json.Marshal(&resources)
	}

	if err != nil {
		return err
	}
	return om.ConfigureProduct(product.Name, string(networkBytes), string(propertiesBytes), string(resourcesBytes))
}
