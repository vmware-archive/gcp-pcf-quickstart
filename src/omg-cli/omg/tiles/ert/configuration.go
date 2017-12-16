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
	AppsDomain tiles.Value `json:".cloud_controller.apps_domain"`
	SysDomain  tiles.Value `json:".cloud_controller.system_domain"`
	// Networking
	NetworkingPointOfEntry    tiles.Value            `json:".properties.networking_point_of_entry"`
	TcpRouting                tiles.Value            `json:".properties.tcp_routing"`
	TcpRoutingReservablePorts tiles.Value            `json:".properties.tcp_routing.enable.reservable_ports"`
	GoRouterSSLCiphers        tiles.Value            `json:".properties.gorouter_ssl_ciphers"`
	HAProxySSLCiphers         tiles.Value            `json:".properties.haproxy_ssl_ciphers"`
	SkipSSLVerification       tiles.BooleanValue     `json:".ha_proxy.skip_cert_verify"`
	HAProxyForwardTLS         tiles.Value            `json:".properties.haproxy_forward_tls"`
	IngressCertificates       tiles.CertificateValue `json:".properties.networking_poe_ssl_cert"`
	// Application Containers
	ContainerDNSServers tiles.Value `json:".diego_cell.dns_servers"`
	// Application Security Groups
	SecurityAcknowledgement tiles.Value `json:".properties.security_acknowledgement"`
	// UAA
	ServiceProviderCredentials tiles.CertificateValue `json:".uaa.service_provider_key_credentials"`

	UaaDbChoice   tiles.Value        `json:".properties.uaa_database"`
	UaaDbIp       tiles.Value        `json:".properties.uaa_database.external.host"`
	UaaDbPort     tiles.IntegerValue `json:".properties.uaa_database.external.port"`
	UaaDbUsername tiles.Value        `json:".properties.uaa_database.external.uaa_username"`
	UaaDbPassword tiles.SecretValue  `json:".properties.uaa_database.external.uaa_password"`

	// Databases
	ErtDbChoice tiles.Value        `json:".properties.system_database"`
	ErtDbIp     tiles.Value        `json:".properties.system_database.external.host"`
	ErtDbPort   tiles.IntegerValue `json:".properties.system_database.external.port"`

	ErtDbAppUsageUsername            tiles.Value       `json:".properties.system_database.external.app_usage_service_username"`
	ErtDbAppUsagePassword            tiles.SecretValue `json:".properties.system_database.external.app_usage_service_password"`
	ErtDbAutoscaleUsername           tiles.Value       `json:".properties.system_database.external.autoscale_username"`
	ErtDbAutoscalePassword           tiles.SecretValue `json:".properties.system_database.external.autoscale_password"`
	ErtDbCloudControllerUsername     tiles.Value       `json:".properties.system_database.external.ccdb_username"`
	ErtDbCloudControllerPassword     tiles.SecretValue `json:".properties.system_database.external.ccdb_password"`
	ErtDbDiegoUsername               tiles.Value       `json:".properties.system_database.external.diego_username"`
	ErtDbDiegoPassword               tiles.SecretValue `json:".properties.system_database.external.diego_password"`
	ErtDbLocketUsername              tiles.Value       `json:".properties.system_database.external.locket_username"`
	ErtDbLocketPassword              tiles.SecretValue `json:".properties.system_database.external.locket_password"`
	ErtDbNetworkPolicyServerUsername tiles.Value       `json:".properties.system_database.external.networkpolicyserver_username"`
	ErtDbNetworkPolicyServerPassword tiles.SecretValue `json:".properties.system_database.external.networkpolicyserver_password"`
	ErtDbNfsUsername                 tiles.Value       `json:".properties.system_database.external.nfsvolume_username"`
	ErtDbNfsPassword                 tiles.SecretValue `json:".properties.system_database.external.nfsvolume_password"`
	ErtDbNotificationsUsername       tiles.Value       `json:".properties.system_database.external.notifications_username"`
	ErtDbNotificationsPassword       tiles.SecretValue `json:".properties.system_database.external.notifications_password"`
	ErtDbAccountUsername             tiles.Value       `json:".properties.system_database.external.account_username"`
	ErtDbAccountPassword             tiles.SecretValue `json:".properties.system_database.external.account_password"`
	ErtDbRoutingUsername             tiles.Value       `json:".properties.system_database.external.routing_username"`
	ErtDbRoutingPassword             tiles.SecretValue `json:".properties.system_database.external.routing_password"`
	ErtDbSilkUsername                tiles.Value       `json:".properties.system_database.external.silk_username"`
	ErtDbSilkPassword                tiles.SecretValue `json:".properties.system_database.external.silk_password"`

	// MySQL
	MySqlMonitorRecipientEmail tiles.Value `json:".mysql_monitor.recipient_email"`
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
	BackupPrepare                tiles.Resource `json:"backup-prepare"`
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
	SmokeTests                   tiles.Resource `json:"smoke-tests"`
	PushAppsManager              tiles.Resource `json:"push-apps-manager"`
	Notifications                tiles.Resource `json:"notifications"`
	NotificationsUi              tiles.Resource `json:"notifications-ui"`
	PushPivotalAccount           tiles.Resource `json:"push-pivotal-account"`
	PushUsageService             tiles.Resource `json:"push-usage-service"`
	Autoscaling                  tiles.Resource `json:"autoscaling"`
	AutoscalingRegisterBroker    tiles.Resource `json:"autoscaling-register-broker"`
	Nfsbrokerpush                tiles.Resource `json:"nfsbrokerpush"`
	Bootstrap                    tiles.Resource `json:"bootstrap"`
	MysqlRejoinUnsafe            tiles.Resource `json:"mysql-rejoin-unsafe"`
}

type SmallFootprintResources struct {
	TcpRouter tiles.Resource `json:"tcp_router"`
	Router    tiles.Resource `json:"router"`

	Database    tiles.Resource `json:"database"`
	Control     tiles.Resource `json:"control"`
	Compute     tiles.Resource `json:"compute"`
	FileStorage tiles.Resource `json:"blobstore"`

	HaProxy                   tiles.Resource `json:"ha_proxy"`
	BackupPrepare             tiles.Resource `json:"backup-prepare"`
	MysqlMonitor              tiles.Resource `json:"mysql_monitor"`
	SmokeTests                tiles.Resource `json:"smoke-tests"`
	PushAppsManager           tiles.Resource `json:"push-apps-manager"`
	Notifications             tiles.Resource `json:"notifications"`
	NotificationsUi           tiles.Resource `json:"notifications-ui"`
	PushPivotalAccount        tiles.Resource `json:"push-pivotal-account"`
	PushUsageService          tiles.Resource `json:"push-usage-service"`
	Autoscaling               tiles.Resource `json:"autoscaling"`
	AutoscalingRegisterBroker tiles.Resource `json:"autoscaling-register-broker"`
	Nfsbrokerpush             tiles.Resource `json:"nfsbrokerpush"`
	Bootstrap                 tiles.Resource `json:"bootstrap"`
	MysqlRejoinUnsafe         tiles.Resource `json:"mysql-rejoin-unsafe"`
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
		AppsDomain:                 tiles.Value{cfg.AppsDomain},
		SysDomain:                  tiles.Value{cfg.SysDomain},
		NetworkingPointOfEntry:     tiles.Value{"external_non_ssl"},
		ContainerDNSServers:        tiles.Value{"8.8.8.8,8.8.4.4"},
		SkipSSLVerification:        tiles.BooleanValue{true},
		HAProxyForwardTLS:          tiles.Value{"disable"},
		IngressCertificates:        tiles.CertificateValue{tiles.Certificate{cfg.SslCertificate, cfg.SslPrivateKey}},
		TcpRouting:                 tiles.Value{"enable"},
		TcpRoutingReservablePorts:  tiles.Value{cfg.TcpPortRange},
		GoRouterSSLCiphers:         tiles.Value{"ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384"},
		HAProxySSLCiphers:          tiles.Value{"DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384"},
		SecurityAcknowledgement:    tiles.Value{"X"},
		ServiceProviderCredentials: tiles.CertificateValue{tiles.Certificate{cfg.SslCertificate, cfg.SslPrivateKey}},

		UaaDbChoice:   tiles.Value{"external"},
		UaaDbIp:       tiles.Value{cfg.ExternalSqlIp},
		UaaDbPort:     tiles.IntegerValue{cfg.ExternalSqlPort},
		UaaDbUsername: tiles.Value{cfg.ERTSqlUsername},
		UaaDbPassword: tiles.SecretValue{tiles.Secret{cfg.ERTSqlPassword}},

		ErtDbChoice:                      tiles.Value{"external"},
		ErtDbIp:                          tiles.Value{cfg.ExternalSqlIp},
		ErtDbPort:                        tiles.IntegerValue{cfg.ExternalSqlPort},
		ErtDbAppUsageUsername:            tiles.Value{cfg.ERTSqlUsername},
		ErtDbAppUsagePassword:            tiles.SecretValue{tiles.Secret{cfg.ERTSqlPassword}},
		ErtDbAutoscaleUsername:           tiles.Value{cfg.ERTSqlUsername},
		ErtDbAutoscalePassword:           tiles.SecretValue{tiles.Secret{cfg.ERTSqlPassword}},
		ErtDbCloudControllerUsername:     tiles.Value{cfg.ERTSqlUsername},
		ErtDbCloudControllerPassword:     tiles.SecretValue{tiles.Secret{cfg.ERTSqlPassword}},
		ErtDbDiegoUsername:               tiles.Value{cfg.ERTSqlUsername},
		ErtDbDiegoPassword:               tiles.SecretValue{tiles.Secret{cfg.ERTSqlPassword}},
		ErtDbLocketUsername:              tiles.Value{cfg.ERTSqlUsername},
		ErtDbLocketPassword:              tiles.SecretValue{tiles.Secret{cfg.ERTSqlPassword}},
		ErtDbNetworkPolicyServerUsername: tiles.Value{cfg.ERTSqlUsername},
		ErtDbNetworkPolicyServerPassword: tiles.SecretValue{tiles.Secret{cfg.ERTSqlPassword}},
		ErtDbNfsUsername:                 tiles.Value{cfg.ERTSqlUsername},
		ErtDbNfsPassword:                 tiles.SecretValue{tiles.Secret{cfg.ERTSqlPassword}},
		ErtDbNotificationsUsername:       tiles.Value{cfg.ERTSqlUsername},
		ErtDbNotificationsPassword:       tiles.SecretValue{tiles.Secret{cfg.ERTSqlPassword}},
		ErtDbAccountUsername:             tiles.Value{cfg.ERTSqlUsername},
		ErtDbAccountPassword:             tiles.SecretValue{tiles.Secret{cfg.ERTSqlPassword}},
		ErtDbRoutingUsername:             tiles.Value{cfg.ERTSqlUsername},
		ErtDbRoutingPassword:             tiles.SecretValue{tiles.Secret{cfg.ERTSqlPassword}},
		ErtDbSilkUsername:                tiles.Value{cfg.ERTSqlUsername},
		ErtDbSilkPassword:                tiles.SecretValue{tiles.Secret{cfg.ERTSqlPassword}},

		MySqlMonitorRecipientEmail: tiles.Value{"admin@example.org"},
	}

	propertiesBytes, err := json.Marshal(&properties)
	if err != nil {
		return err
	}

	resorucesBytes := []byte{}

	if envConfig.SmallFootprint {
		resoruces := SmallFootprintResources{
			TcpRouter: tiles.Resource{
				RouterNames:       []string{fmt.Sprintf("tcp:%s", cfg.TcpTargetPoolName)},
				InternetConnected: false,
				Instances:         1,
			},
			Router: tiles.Resource{
				RouterNames:       []string{fmt.Sprintf("tcp:%s", cfg.WssTargetPoolName), fmt.Sprintf("http:%s", cfg.HttpBackendServiceName)},
				InternetConnected: false,
			},
			Control: tiles.Resource{
				RouterNames:       []string{fmt.Sprintf("tcp:%s", cfg.SshTargetPoolName)},
				InternetConnected: false,
			},
			HaProxy: tiles.Resource{
				Instances: 0,
			},
			MysqlMonitor: tiles.Resource{
				Instances: 0,
			},
		}
		resorucesBytes, err = json.Marshal(&resoruces)
	} else {
		resoruces := LargeFootprintResources{
			TcpRouter: tiles.Resource{
				RouterNames:       []string{fmt.Sprintf("tcp:%s", cfg.TcpTargetPoolName)},
				InternetConnected: false,
				Instances:         1,
			},
			Router: tiles.Resource{
				RouterNames:       []string{fmt.Sprintf("tcp:%s", cfg.WssTargetPoolName), fmt.Sprintf("http:%s", cfg.HttpBackendServiceName)},
				InternetConnected: false,
			},
			DiegoBrain: tiles.Resource{
				RouterNames:       []string{fmt.Sprintf("tcp:%s", cfg.SshTargetPoolName)},
				InternetConnected: false,
			},
			HaProxy: tiles.Resource{
				Instances: 0,
			},
			MysqlMonitor: tiles.Resource{
				Instances: 0,
			},
		}
		resorucesBytes, err = json.Marshal(&resoruces)
	}

	if err != nil {
		return err
	}
	return om.ConfigureProduct(product.Name, string(networkBytes), string(propertiesBytes), string(resorucesBytes))
}
