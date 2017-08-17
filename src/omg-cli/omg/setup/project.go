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

func (pv *ProjectValidator) EnsureQuota() ([]QuotaError, error) {
	quotas, err := pv.project.Quotas()
	if err != nil {
		return nil, err
	}

	errors := []QuotaError{}

	for _, requirement := range pv.requirements {
		quota, ok := quotas[requirement.Name]
		if !ok {
			errors = append(errors, QuotaError{requirement, 0})
		} else {
			if quota.Limit < requirement.Limit {
				errors = append(errors, QuotaError{requirement, quota.Limit})
			}
		}
	}

	if len(errors) != 0 {
		return errors, UnsatisfiedQuotaErr
	}

	return nil, nil
}
