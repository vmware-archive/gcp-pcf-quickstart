package commands

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

import (
	"fmt"
	"log"

	"omg-cli/config"
	"omg-cli/google"
	"omg-cli/omg/setup"

	"github.com/alecthomas/kingpin"
	"golang.org/x/oauth2"
	googleauth "golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

type PrepareProjectCommand struct {
	logger              *log.Logger
	terraformConfigPath string
}

const PrepareProjectName = "prepare-project"

func (ppc *PrepareProjectCommand) register(app *kingpin.Application) {
	c := app.Command(PrepareProjectName, "Prepare the GCP Project").Action(ppc.run)
	registerTerraformConfigFlag(c, &ppc.terraformConfigPath)
}

func (ppc *PrepareProjectCommand) run(c *kingpin.ParseContext) error {
	cfg, err := config.FromTerraform(ppc.terraformConfigPath)
	if err != nil {
		return err
	}

	creds, err := googleauth.JWTConfigFromJSON([]byte(cfg.ServiceAccountKey), compute.CloudPlatformScope)
	if err != nil {
		return fmt.Errorf("loading ServiceAccountKey: %v ", err)
	}
	gcpClient := creds.Client(oauth2.NoContext)

	quotaService, err := google.NewQuotaService(ppc.logger, cfg.ProjectName, gcpClient)
	if err != nil {
		return fmt.Errorf("creating QuotaService: %v", err)
	}

	validator, err := setup.NewProjectValiadtor(ppc.logger, quotaService, setup.ProjectQuotaRequirements(), setup.RegionalQuotaRequirements(cfg))
	if err != nil {
		return fmt.Errorf("creating ProjectValidator: %v", err)
	}

	errors, satisfied, err := validator.EnsureQuota()
	if err == nil {
		ppc.logger.Printf("quotaService quota is adequate, satisfied %v rules", len(satisfied))
		return nil
	}

	if err != setup.UnsatisfiedQuotaErr {
		return fmt.Errorf("error validating quota: %v", err)
	}

	for _, quotaError := range errors {
		ppc.logger.Printf("Compute Engine quota requirement not satisfied: Name %s, Region: %s, Minimum Required: %v (Current Quota: %v)", quotaError.Name, quotaError.Region, quotaError.Limit, quotaError.Actual)
	}

	if err != nil {
		ppc.logger.Fatal(err)
	}

	return nil
}
