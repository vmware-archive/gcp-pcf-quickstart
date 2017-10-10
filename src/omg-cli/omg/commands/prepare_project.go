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
	"context"
	"log"
	"net/http"

	"omg-cli/config"
	"omg-cli/google"
	"omg-cli/omg/setup"

	"github.com/alecthomas/kingpin"
	googleauth "golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

type PrepareProjectCommand struct {
	logger *log.Logger
	envDir string
}

const PrepareProjectName = "prepare-project"

func (cmd *PrepareProjectCommand) register(app *kingpin.Application) {
	c := app.Command(PrepareProjectName, "Prepare the GCP Project").Action(cmd.run)
	registerEnvConfigFlag(c, &cmd.envDir)
}

func (cmd *PrepareProjectCommand) parseArgs() (cfg *config.EnvConfig, gcpClient *http.Client) {
	var err error
	cfg, err = config.FromEnvironmentDirectory(cmd.envDir)
	if err != nil {
		cmd.logger.Fatalf("loading environment config: %v", err)
	}

	gcpClient, err = googleauth.DefaultClient(context.Background(), compute.CloudPlatformScope)
	if err != nil {
		cmd.logger.Fatalf("loading application default credentials: %v.\nHave you ran `gcloud auth application-default login`?", err)
	}

	return
}

func (cmd *PrepareProjectCommand) run(c *kingpin.ParseContext) error {
	cfg, gcpClient := cmd.parseArgs()
	validator := cmd.createValidator(cfg, gcpClient)

	validateApis(cmd.logger, validator)
	validateQuotas(cmd.logger, validator)

	return nil
}

func (cmd *PrepareProjectCommand) createValidator(cfg *config.EnvConfig, gcpClient *http.Client) *setup.ProjectValidator {
	quotaService, err := google.NewQuotaService(cmd.logger, cfg.ProjectID, gcpClient)
	if err != nil {
		cmd.logger.Fatalf("creating QuotaService: %v", err)
	}

	apiService, err := google.NewAPIService(cmd.logger, cfg.ProjectID, gcpClient)
	if err != nil {
		cmd.logger.Fatalf("creating ApiService: %v", err)
	}

	validator, err := setup.NewProjectValidator(cmd.logger, quotaService, apiService, setup.ProjectQuotaRequirements(), setup.RegionalQuotaRequirements(cfg), setup.RequiredAPIs())
	if err != nil {
		cmd.logger.Fatalf("creating ProjectValidator: %v", err)
	}

	return validator
}

func validateQuotas(logger *log.Logger, validator *setup.ProjectValidator) {
	logger.Printf("validating Google Cloud Compute Engine quotas")
	errors, satisfied, err := validator.ValidateQuotas()

	for _, quotaError := range errors {
		logger.Printf("Compute Engine quota requirement not satisfied: Name: %s, Region: %s, Minimum Required: %v (Current Quota: %v)", quotaError.Name, quotaError.Region, quotaError.Limit, quotaError.Actual)
	}

	if err != nil {
		logger.Fatal(err)
	}
	logger.Printf("Compute Engine quota is adequate, satisfied %v rules", len(satisfied))
}

func validateApis(logger *log.Logger, validator *setup.ProjectValidator) {
	logger.Printf("enabling Google Cloud APIs")
	enabledApis, err := validator.EnableAPIs()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Printf("enusred %d API(s) are enabled", len(enabledApis))
}
