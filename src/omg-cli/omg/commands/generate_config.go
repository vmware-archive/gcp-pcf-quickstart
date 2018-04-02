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

package commands

import (
	"context"
	"log"

	"omg-cli/config"
	"omg-cli/google"
	"omg-cli/version"

	"encoding/json"

	"fmt"
	"io/ioutil"

	"github.com/alecthomas/kingpin"
	googleauth "golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

type GenerateConfigCommand struct {
	logger         *log.Logger
	envDir         string
	dnsZone        string
	pivnetApiToken string
	baseZone       string
	projectId      string
}

const GenerateConfigCommandName = "generate-config"

func (cmd *GenerateConfigCommand) register(app *kingpin.Application) {
	c := app.Command(GenerateConfigCommandName, "Generate default environment configuration").Action(cmd.run)
	registerEnvConfigFlag(c, &cmd.envDir)
	registerPivnetApiTokenFlag(c, &cmd.pivnetApiToken)
	c.Flag("gcp-project", "Google Cloud Project ID for deployment").Required().StringVar(&cmd.projectId)
	c.Flag("zone", "Base Zone used for deployment location. Other zones in the region will be used for the deployment.").Default("us-east1-b").StringVar(&cmd.baseZone)
	c.Flag("dns-zone", "Existing Cloud DNS Zone used to create DNS records for deployment").Default("pcf-zone").StringVar(&cmd.dnsZone)
}

func (cmd *GenerateConfigCommand) run(c *kingpin.ParseContext) error {
	cfg, err := config.DefaultEnvConfig()
	if err != nil {
		return err
	}

	gcpClient, err := googleauth.DefaultClient(context.Background(), compute.CloudPlatformScope)
	if err != nil {
		cmd.logger.Fatalf("loading application default credentials: %v.\nHave you ran `gcloud auth application-default login`?", err)
	}
	computeService, err := compute.New(gcpClient)
	if err != nil {
		cmd.logger.Fatalf("creating new compute client: %v", err)
	}
	computeService.UserAgent = version.UserAgent()

	zoneResult, err := google.ParseZone(cmd.projectId, cmd.baseZone, computeService)
	if err != nil {
		return fmt.Errorf("parsing zone %s: %v", cmd.baseZone, err)
	}
	cfg.Region = zoneResult.Region
	cfg.Zone1 = zoneResult.Zone1
	cfg.Zone2 = zoneResult.Zone2
	cfg.Zone3 = zoneResult.Zone3
	cfg.DnsZoneName = cmd.dnsZone
	cfg.PivnetApiToken = cmd.pivnetApiToken
	cfg.ProjectID = cmd.projectId

	cfgStr, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	cfgPath := fmt.Sprintf("%s/config.json", cmd.envDir)
	if err := ioutil.WriteFile(cfgPath, cfgStr, 0644); err != nil {
		return err
	}

	cmd.logger.Printf("successfully wrote default config to: %s", cfgPath)

	return nil
}
