package google

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
	"errors"
	"log"

	"net/http"

	compute "google.golang.org/api/compute/v1"
)

type Quota struct {
	Name  string
	Limit float64
}

//go:generate counterfeiter ./ QuotaService
type QuotaService interface {
	Project() (map[string]Quota, error)
	Region(string) (map[string]Quota, error)
}

type quotaService struct {
	logger         *log.Logger
	projectId      string
	computeService *compute.Service
}

func transformQuotas(computeQuotas []*compute.Quota) map[string]Quota {
	quotas := map[string]Quota{}
	for _, quota := range computeQuotas {
		quotas[quota.Metric] = Quota{quota.Metric, quota.Limit}
	}

	return quotas
}

func (ps *quotaService) Project() (map[string]Quota, error) {
	project, err := ps.computeService.Projects.Get(ps.projectId).Context(context.Background()).Do()
	if err != nil {
		return nil, err
	}

	return transformQuotas(project.Quotas), nil
}

func (ps *quotaService) Region(region string) (map[string]Quota, error) {
	regionResponse, err := ps.computeService.Regions.Get(ps.projectId, region).Context(context.Background()).Do()
	if err != nil {
		return nil, err
	}

	return transformQuotas(regionResponse.Quotas), nil
}

func NewQuotaService(logger *log.Logger, projectId string, client *http.Client) (QuotaService, error) {
	if logger == nil {
		return nil, errors.New("missing logger")
	}

	computeService, err := compute.New(client)
	if err != nil {
		return nil, err
	}

	return &quotaService{logger, projectId, computeService}, nil
}
