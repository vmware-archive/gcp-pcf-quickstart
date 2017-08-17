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

//go:generate counterfeiter ./ ProjectService
type ProjectService interface {
	Quotas() (map[string]Quota, error)
}

type projectService struct {
	logger    *log.Logger
	projectId string
	client    *http.Client
}

func (ps *projectService) Quotas() (map[string]Quota, error) {
	ctx := context.Background()

	computeService, err := compute.New(ps.client)
	if err != nil {
		return nil, err
	}

	project, err := computeService.Projects.Get(ps.projectId).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	quotas := map[string]Quota{}
	for _, quota := range project.Quotas {
		quotas[quota.Metric] = Quota{quota.Metric, quota.Limit}
	}

	return quotas, nil
}

func NewProjectService(logger *log.Logger, projectId string, client *http.Client) (ProjectService, error) {
	if logger == nil {
		return nil, errors.New("missing logger")
	}

	return &projectService{logger, projectId, client}, nil
}
