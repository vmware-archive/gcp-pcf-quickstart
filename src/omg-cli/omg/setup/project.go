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

// ProjectValidator validates a Google Cloud project is capable of launching the quickstart.
type ProjectValidator struct {
	logger              *log.Logger
	quotaService        google.QuotaService
	apiService          google.APIService
	projectRequirements []google.Quota
	regionRequirements  map[string][]google.Quota
	apiRequirements     []google.API
}

// QuotaError represents when a quota cannot be met.
type QuotaError struct {
	google.Quota
	Actual float64
	Region string
}

// ErrUnsatisfiedQuota is thrown when a quota is not met.
var ErrUnsatisfiedQuota = errors.New("unsatisfied quota, request an increase at: https://console.cloud.google.com/iam-admin/quotas")

// NewProjectValidator creates a new project validator.
func NewProjectValidator(logger *log.Logger, quotaService google.QuotaService, apiService google.APIService, projectRequirements []google.Quota, regionRequirements map[string][]google.Quota, apiRequirements []google.API) (*ProjectValidator, error) {
	if logger == nil {
		return nil, errors.New("missing logger")
	}
	return &ProjectValidator{logger, quotaService, apiService, projectRequirements, regionRequirements, apiRequirements}, nil
}

// ValidateQuotas ensures all required quotas are high enough for the quickstart.
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
		err = ErrUnsatisfiedQuota
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

// EnableAPIs enables the required Google Cloud APIs.
func (pv *ProjectValidator) EnableAPIs() ([]google.API, error) {
	enabled, err := pv.apiService.Enable(pv.apiRequirements)
	if err != nil {
		return nil, err
	}
	return enabled, nil
}

// ProjectQuotaRequirements are quota requirements for a Google Cloud project.
func ProjectQuotaRequirements() []google.Quota {
	return []google.Quota{
		{Name: "NETWORKS", Limit: 2.0},
		{Name: "FIREWALLS", Limit: 7.0},
		{Name: "IMAGES", Limit: 15.0},
		{Name: "STATIC_ADDRESSES", Limit: 1.0},
		{Name: "ROUTES", Limit: 20.0},
		{Name: "FORWARDING_RULES", Limit: 4.0},
		{Name: "TARGET_POOLS", Limit: 2.0},
		{Name: "HEALTH_CHECKS", Limit: 3.0},
		{Name: "IN_USE_ADDRESSES", Limit: 4.0},
		{Name: "TARGET_HTTP_PROXIES", Limit: 1.0},
		{Name: "URL_MAPS", Limit: 2.0},
		{Name: "BACKEND_SERVICES", Limit: 4.0},
		{Name: "TARGET_HTTPS_PROXIES", Limit: 1.0},
		{Name: "SSL_CERTIFICATES", Limit: 1.0},
		{Name: "SUBNETWORKS", Limit: 15.0},
	}
}

// RegionalQuotaRequirements are quotas requirements for the deployment Google Cloud region.
func RegionalQuotaRequirements(cfg *config.EnvConfig) map[string][]google.Quota {
	quotas := []google.Quota{
		{Name: "DISKS_TOTAL_GB", Limit: 2000.0},
		{Name: "STATIC_ADDRESSES", Limit: 5.0},
		{Name: "IN_USE_ADDRESSES", Limit: 6.0},
		{Name: "INSTANCE_GROUPS", Limit: 10.0},
		{Name: "INSTANCES", Limit: 100.0},
	}

	if cfg.SmallFootprint {
		quotas = append(quotas, google.Quota{Name: "CPUS", Limit: 24.0})
	} else {
		quotas = append(quotas, google.Quota{Name: "CPUS", Limit: 200.0})
	}

	return map[string][]google.Quota{cfg.Region: quotas}
}

// RequiredAPIs is the set of Google Cloud APIs which must be enabled.
func RequiredAPIs() []google.API {
	return []google.API{
		{Name: "bigquery-json.googleapis.com"},
		{Name: "cloudbuild.googleapis.com"},
		{Name: "clouddebugger.googleapis.com"},
		{Name: "cloudresourcemanager.googleapis.com"},
		{Name: "datastore.googleapis.com"},
		{Name: "storage-component.googleapis.com"},
		{Name: "pubsub.googleapis.com"},
		{Name: "vision.googleapis.com"},
		{Name: "storage-api.googleapis.com"},
		{Name: "logging.googleapis.com"},
		{Name: "resourceviews.googleapis.com"},
		{Name: "replicapool.googleapis.com"},
		{Name: "cloudapis.googleapis.com"},
		{Name: "deploymentmanager.googleapis.com"},
		{Name: "containerregistry.googleapis.com"},
		{Name: "sqladmin.googleapis.com"},
		{Name: "monitoring.googleapis.com"},
		{Name: "dns.googleapis.com"},
		{Name: "runtimeconfig.googleapis.com"},
		{Name: "compute.googleapis.com"},
		{Name: "sql-component.googleapis.com"},
		{Name: "iam.googleapis.com"},
		{Name: "cloudtrace.googleapis.com"},
		{Name: "servicemanagement.googleapis.com"},
		{Name: "replicapoolupdater.googleapis.com"},
	}
}
