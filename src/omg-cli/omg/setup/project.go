package setup

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
	"errors"
	"log"
	"omg-cli/config"
	"omg-cli/google"
)

type ProjectValidator struct {
	logger              *log.Logger
	quotaService        google.QuotaService
	apiService          google.APIService
	projectRequirements []google.Quota
	regionRequirements  map[string][]google.Quota
	apiRequirements     []google.API
}

type QuotaError struct {
	google.Quota
	Actual float64
	Region string
}

var UnsatisfiedQuotaErr = errors.New("Compute Engine quota is unsatisfied, request an increase at: https://console.cloud.google.com/iam-admin/quotas")

func NewProjectValidator(logger *log.Logger, quotaService google.QuotaService, apiService google.APIService, projectRequirements []google.Quota, regionRequirements map[string][]google.Quota, apiRequirements []google.API) (*ProjectValidator, error) {
	if logger == nil {
		return nil, errors.New("missing logger")
	}
	return &ProjectValidator{logger, quotaService, apiService, projectRequirements, regionRequirements, apiRequirements}, nil
}

func (pv *ProjectValidator) ValidateQuotas() (errors []QuotaError, satisfied []google.Quota, err error) {
	quotas, err := pv.quotaService.Project()
	if err != nil {
		return nil, nil, err
	}

	errors, satisfied = validateQuotas(pv.projectRequirements, quotas, "global")

	for region, requirements := range pv.regionRequirements {
		quotas, err = pv.quotaService.Region(region)
		regionErrors, regionSatisfied := validateQuotas(requirements, quotas, region)
		errors = append(errors, regionErrors...)
		satisfied = append(satisfied, regionSatisfied...)
	}

	if len(errors) != 0 {
		err = UnsatisfiedQuotaErr
	}
	return errors, satisfied, err
}

func validateQuotas(requirements []google.Quota, quotas map[string]google.Quota, region string) (errors []QuotaError, satisfied []google.Quota) {
	errors = []QuotaError{}
	satisfied = []google.Quota{}

	for _, requirement := range requirements {
		quota, ok := quotas[requirement.Name]
		if !ok {
			errors = append(errors, QuotaError{requirement, 0, region})
		} else {
			if quota.Limit < requirement.Limit {
				errors = append(errors, QuotaError{requirement, quota.Limit, region})
			} else {
				satisfied = append(satisfied, requirement)
			}
		}
	}

	return
}

func (pv *ProjectValidator) EnableAPIs() ([]google.API, error) {
	enabled, err := pv.apiService.Enable(pv.apiRequirements)
	if err != nil {
		return nil, err
	}
	return enabled, nil
}

func ProjectQuotaRequirements() []google.Quota {
	return []google.Quota{
		{"NETWORKS", 2.0},
		{"FIREWALLS", 7.0},
		{"IMAGES", 15.0},
		{"STATIC_ADDRESSES", 1.0},
		{"ROUTES", 20.0},
		{"FORWARDING_RULES", 4.0},
		{"TARGET_POOLS", 2.0},
		{"HEALTH_CHECKS", 3.0},
		{"IN_USE_ADDRESSES", 4.0},
		{"TARGET_HTTP_PROXIES", 1.0},
		{"URL_MAPS", 2.0},
		{"BACKEND_SERVICES", 4.0},
		{"TARGET_HTTPS_PROXIES", 1.0},
		{"SSL_CERTIFICATES", 1.0},
		{"SUBNETWORKS", 15.0},
	}
}

func RegionalQuotaRequirements(cfg *config.EnvConfig) map[string][]google.Quota {
	quotas := []google.Quota{
		{"DISKS_TOTAL_GB", 2000.0},
		{"STATIC_ADDRESSES", 5.0},
		{"IN_USE_ADDRESSES", 6.0},
		{"INSTANCE_GROUPS", 10.0},
		{"INSTANCES", 100.0},
	}

	if cfg.SmallFootprint {
		quotas = append(quotas, google.Quota{"CPUS", 24.0})
	} else {
		quotas = append(quotas, google.Quota{"CPUS", 200.0})
	}

	return map[string][]google.Quota{cfg.Region: quotas}
}

func RequiredAPIs() []google.API {
	return []google.API{
		{"bigquery-json.googleapis.com"},
		{"cloudbuild.googleapis.com"},
		{"clouddebugger.googleapis.com"},
		{"cloudresourcemanager.googleapis.com"},
		{"datastore.googleapis.com"},
		{"storage-component.googleapis.com"},
		{"pubsub.googleapis.com"},
		{"vision.googleapis.com"},
		{"storage-api.googleapis.com"},
		{"logging.googleapis.com"},
		{"resourceviews.googleapis.com"},
		{"replicapool.googleapis.com"},
		{"cloudapis.googleapis.com"},
		{"deploymentmanager.googleapis.com"},
		{"containerregistry.googleapis.com"},
		{"sqladmin.googleapis.com"},
		{"monitoring.googleapis.com"},
		{"dns.googleapis.com"},
		{"runtimeconfig.googleapis.com"},
		{"compute.googleapis.com"},
		{"sql-component.googleapis.com"},
		{"iam.googleapis.com"},
		{"cloudtrace.googleapis.com"},
		{"servicemanagement.googleapis.com"},
		{"replicapoolupdater.googleapis.com"},
	}
}
