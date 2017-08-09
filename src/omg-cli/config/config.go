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

type Config struct {
	OpsManagerIp                       string
	JumpboxIp                          string
	NetworkName                        string
	DeploymentTargetTag                string
	MgmtSubnetName                     string
	MgmtSubnetGateway                  string
	MgmtSubnetCIDR                     string
	ServicesSubnetName                 string
	ServicesSubnetGateway              string
	ServicesSubnetCIDR                 string
	ErtSubnetName                      string
	ErtSubnetGateway                   string
	ErtSubnetCIDR                      string
	HttpBackendServiceName             string
	SshTargetPoolName                  string
	TcpTargetPoolName                  string
	TcpPortRange                       string
	BuildpacksBucket                   string
	DropletsBucket                     string
	PackagesBucket                     string
	ResourcesBucket                    string
	DirectorBucket                     string
	DnsSuffix                          string
	SslCertificate                     string
	SslPrivateKey                      string
	OpsManServiceAccount               string
	StackdriverNozzleServiceAccountKey string

	ServiceBrokerServiceAccountKey string
	ServiceBrokerDbIp              string
	ServiceBrokerDbUsername        string
	ServiceBrokerDbPassword        string

	Region      string
	Zone1       string
	Zone2       string
	Zone3       string
	ProjectName string
}

type OpsManagerCredentials struct {
	Username            string
	Password            string
	DecryptionPhrase    string
	SkipSSLVerification bool
}
