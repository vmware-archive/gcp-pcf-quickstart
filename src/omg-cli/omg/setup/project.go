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
	"omg-cli/google"
)

type ProjectValidator struct {
	logger       *log.Logger
	project      google.ProjectService
	requirements []google.Quota
}

type QuotaError struct {
	google.Quota
	Actual float64
}

var UnsatisfiedQuotaErr = errors.New("quota unsatisfied")

func NewProjectValiadtor(logger *log.Logger, projectService google.ProjectService, requirements []google.Quota) (*ProjectValidator, error) {
	if logger == nil {
		return nil, errors.New("missing logger")
	}
	return &ProjectValidator{logger, projectService, requirements}, nil
}

func (pv *ProjectValidator) EnsureQuota() (errors []QuotaError, satisfied []google.Quota, err error) {
	quotas, err := pv.project.Quotas()
	if err != nil {
		return nil, nil, err
	}

	errors = []QuotaError{}
	satisfied = []google.Quota{}

	for _, requirement := range pv.requirements {
		quota, ok := quotas[requirement.Name]
		if !ok {
			errors = append(errors, QuotaError{requirement, 0})
		} else {
			if quota.Limit < requirement.Limit {
				errors = append(errors, QuotaError{requirement, quota.Limit})
			} else {
				satisfied = append(satisfied, requirement)
			}
		}
	}

	if len(errors) != 0 {
		err = UnsatisfiedQuotaErr
	}
	return
}

func QuotaRequirements() []google.Quota {
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
