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

package gcp_director

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"text/template"

	"omg-cli/config"

	"omg-cli/ops_manager"
)

const directorTemplateYAML = `---
az-configuration:
- name: "{{.Zone1}}"
- name: "{{.Zone2}}"
- name: "{{.Zone3}}"
director-configuration:
  ntp_servers_string: 169.254.169.254
  retry_bosh_deploys: true
  resurrector_enabled: true
  max_threads: 5
  blobstore_type: gcs
  gcs_blobstore_options:
    bucket_name: "{{.DirectorBucket}}"
    service_account_key: '{{.OpsManagerServiceAccountKey}}'
    storage_class: MULTI_REGIONAL
  database_type: "{{.DatabaseType}}"
{{if eq .DatabaseType "external"}}
  external_database_options:
    host: "{{.ExternalSqlIp}}"
    database: "{{.OpsManagerSqlDbName}}"
    user: "{{.OpsManagerSqlUsername}}"
    password: "{{.OpsManagerSqlPassword}}"
    port: {{.ExternalSqlPort}}
{{end}}
iaas-configuration:
  project: "{{.ProjectName}}"
  auth_json: '{{.OpsManagerServiceAccountKey}}'
  default_deployment_tag: "{{.DeploymentTargetTag}}"
network-assignment:
  singleton_availability_zone:
    name: "{{.Zone1}}"
  network:
    name: "{{.MgmtSubnetName}}"
networks-configuration:
  icmp: false
  networks:
  - name: "{{.MgmtSubnetName}}"
    subnets:
    - iaas_identifier: "{{.NetworkName}}/{{.MgmtSubnetName}}/{{.Region}}"
      cidr: "{{.MgmtSubnetCIDR}}"
      gateway: "{{.MgmtSubnetGateway}}"
      reserved_ip_ranges: "{{reservedIPs .MgmtSubnetCIDR}}"
      dns: 169.254.169.254
      availability_zone_names:
      - "{{.Zone1}}"
      - "{{.Zone2}}"
      - "{{.Zone3}}"
  - name: "{{.ServicesSubnetName}}"
    subnets:
    - iaas_identifier: "{{.NetworkName}}/{{.ServicesSubnetName}}/{{.Region}}"
      cidr: "{{.ServicesSubnetCIDR}}"
      reserved_ip_ranges: "{{reservedIPs .ServicesSubnetCIDR}}"
      gateway: "{{.ServicesSubnetGateway}}"
      dns: 169.254.169.254
      availability_zone_names:
      - "{{.Zone1}}"
      - "{{.Zone2}}"
      - "{{.Zone3}}"
  - name: "{{.ErtSubnetName}}"
    subnets:
    - iaas_identifier: "{{.NetworkName}}/{{.ErtSubnetName}}/{{.Region}}"
      cidr: "{{.ErtSubnetCIDR}}"
      reserved_ip_ranges: "{{reservedIPs .ErtSubnetCIDR}}"
      gateway: "{{.ErtSubnetGateway}}"
      dns: 169.254.169.254
      availability_zone_names:
      - "{{.Zone1}}"
      - "{{.Zone2}}"
      - "{{.Zone3}}"
resource-configuration:
  compilation:
    instances: {{.CompilationInstances}}
    instance_type:
      id: "{{.CompilationInstanceType}}"
  director:
    instances: automatic
    persistent_disk:
      size_mb: automatic
    instance_type:
      id: automatic
    internet_connected: false
`

var funcMap = template.FuncMap{
	// The name "title" is what the function will be called in the template text.
	"reservedIPs": reservedIPs,
}

var tmpl = template.Must(template.New("director").Funcs(funcMap).Parse(directorTemplateYAML))

func reservedIPs(cidr string) string {
	// Reserve .1-.20
	lowerIp, _, err := net.ParseCIDR(cidr)
	lowerIp = lowerIp.To4()
	if err != nil {
		panic(err)
	}
	upperIp := make(net.IP, len(lowerIp))
	copy(upperIp, lowerIp)
	upperIp[3] = 20

	return fmt.Sprintf("%s-%s", lowerIp, upperIp)
}

func (*Tile) Configure(envConfig *config.EnvConfig, cfg *config.Config, om *ops_manager.Sdk) error {
	dc := struct {
		config.Config
		CompilationInstances    int
		CompilationInstanceType string
		DatabaseType            string
	}{
		Config: *cfg,
	}
	if envConfig.SmallFootprint {
		dc.CompilationInstances = 1
		dc.CompilationInstanceType = "medium.mem"
		dc.DatabaseType = "internal"
	} else {
		dc.CompilationInstances = 4
		dc.CompilationInstanceType = "large.cpu"
		dc.DatabaseType = "external"
	}
	// Healthwatch includes a C++ package that requires a large
	// ephemeral disk for compilation.
	if envConfig.IncludeHealthwatch {
		dc.CompilationInstanceType = "large.disk"
	}

	// strip newlines, as we're putting a JSON service account key inside YAML
	dc.OpsManagerServiceAccountKey = strings.Replace(dc.OpsManagerServiceAccountKey, "\n", "", -1)

	b := &bytes.Buffer{}
	err := tmpl.Execute(b, dc)
	if err != nil {
		return fmt.Errorf("cannot generate director YAML: %v", err)
	}

	return om.SetupBosh(b.Bytes())
}
