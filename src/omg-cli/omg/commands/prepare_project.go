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

	"os"

	"errors"

	"github.com/alecthomas/kingpin"
	"golang.org/x/oauth2"
	googleauth "golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

type PrepareProjectCommand struct {
	logger              *log.Logger
	terraformConfigPath string
	projectId           string
	region              string
}

const PrepareProjectName = "prepare-project"

func (ppc *PrepareProjectCommand) register(app *kingpin.Application) {
	c := app.Command(PrepareProjectName, "Prepare the GCP Project").Action(ppc.run)
	registerTerraformConfigFlag(c, &ppc.terraformConfigPath)
	c.Flag("project-id", "Project ID (if not using terraform-config)").StringVar(&ppc.projectId)
	c.Flag("region", "Region (if not using terraform-config)").StringVar(&ppc.region)
}

func (ppc *PrepareProjectCommand) parseArgs() (cfg *config.Config, gcpClient *http.Client, err error) {
	if _, err := os.Stat(ppc.terraformConfigPath); err == nil {
		cfg, err = config.FromTerraform(ppc.terraformConfigPath)
		if err != nil {
			ppc.logger.Fatalf("loading terraform config: %v", err)
		}

		creds, err := googleauth.JWTConfigFromJSON([]byte(cfg.ServiceAccountKey), compute.CloudPlatformScope)
		if err != nil {
			ppc.logger.Fatalf("loading ServiceAccountKey: %v ", err)
		}
		gcpClient = creds.Client(oauth2.NoContext)

	} else {
		cfg = &config.Config{ProjectName: ppc.projectId, Region: ppc.region}
		gcpClient, err = googleauth.DefaultClient(context.Background(), compute.CloudPlatformScope)
		if err != nil {
			ppc.logger.Fatalf("loading application default credentials: %v.\nHave you ran `gcloud auth application-default login`?", err)
		}

		if ppc.projectId == "" {
			err = errors.New("specify --project-id")
		}

		if ppc.region == "" {
			err = errors.New("specify --region")
		}
	}

	return
}

func (ppc *PrepareProjectCommand) run(c *kingpin.ParseContext) error {
	cfg, gcpClient, err := ppc.parseArgs()
	if err != nil {
		return err
	}

	validator, err := createValidator(ppc.logger, cfg, gcpClient)
	if err != nil {
		ppc.logger.Fatalf("creating ProjectValidator: %v", err)
	}

	validateApis(ppc.logger, validator)
	validateQuotas(ppc.logger, validator)

	return nil
}

func createValidator(logger *log.Logger, cfg *config.Config, gcpClient *http.Client) (*setup.ProjectValidator, error) {
	quotaService, err := google.NewQuotaService(logger, cfg.ProjectName, gcpClient)
	if err != nil {
		logger.Fatalf("creating QuotaService: %v", err)
	}

	apiService, err := google.NewAPIService(logger, cfg.ProjectName, gcpClient)
	if err != nil {
		logger.Fatalf("creating ApiService: %v", err)
	}

	return setup.NewProjectValiadtor(logger, quotaService, apiService, setup.ProjectQuotaRequirements(), setup.RegionalQuotaRequirements(cfg), setup.RequiredAPIs())
}

func validateQuotas(logger *log.Logger, validator *setup.ProjectValidator) {
	logger.Printf("validating Google Cloud Compute Engine quotas")
	errors, satisfied, err := validator.ValidateQuotas()

	for _, quotaError := range errors {
		logger.Printf("Compute Engine quota requirement not satisfied: Name %s, Region: %s, Minimum Required: %v (Current Quota: %v)", quotaError.Name, quotaError.Region, quotaError.Limit, quotaError.Actual)
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
