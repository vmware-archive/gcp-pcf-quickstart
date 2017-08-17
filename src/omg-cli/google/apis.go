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
	"errors"
	"log"
	"net/http"

	"fmt"

	"time"

	servicemanagement "google.golang.org/api/servicemanagement/v1"
)

type API struct {
	Name string
}

//go:generate counterfeiter ./ APIService
type APIService interface {
	Enable([]API) ([]API, error)
}

type apiService struct {
	logger                   *log.Logger
	projectId                string
	servicemanagementService *servicemanagement.APIService
}

func NewAPIService(logger *log.Logger, projectId string, client *http.Client) (APIService, error) {
	if logger == nil {
		return nil, errors.New("logger blank")
	}

	servicemanagementService, err := servicemanagement.New(client)
	if err != nil {
		return nil, err
	}

	return &apiService{logger, projectId, servicemanagementService}, nil
}

func (svc *apiService) Enable(apis []API) ([]API, error) {
	pendingOperations := map[API]string{}
	for _, api := range apis {
		if svc.isEnabled(api) {
			continue
		}

		operation, err := svc.enableOne(api)
		if err != nil {
			return nil, fmt.Errorf("enabling %s: %v", api.Name, err)
		}
		pendingOperations[api] = operation
	}

	apisEnabled := []API{}

	for len(pendingOperations) != 0 {
		svc.logger.Printf("waiting for %d service enable operation(s) to complete", len(pendingOperations))
		time.Sleep(time.Duration(2) * time.Second)

		var err error
		var completed []API
		pendingOperations, completed, err = svc.filterCompleted(pendingOperations)
		if err != nil {
			return apisEnabled, err
		}

		apisEnabled = append(apisEnabled, completed...)
	}

	return apisEnabled, nil
}

func (svc *apiService) enableOne(api API) (operation string, err error) {
	req := &servicemanagement.EnableServiceRequest{
		ConsumerId: fmt.Sprintf("project:%s", svc.projectId),
	}
	oper, err := svc.servicemanagementService.Services.Enable(api.Name, req).Do()
	if err != nil {
		return "", err
	}

	return oper.Name, err
}

func (svc *apiService) isEnabled(api API) bool {
	// TODO(jrjohnson): detect if the API is enabled before trying to enable it
	return false
}

func (svc *apiService) filterCompleted(operations map[API]string) (pending map[API]string, completed []API, err error) {
	pending = map[API]string{}
	completed = []API{}
	for api, operationName := range operations {
		oper, err := svc.servicemanagementService.Operations.Get(operationName).Do()
		if err != nil {
			return operations, completed, fmt.Errorf("pooling for api: %s, operation: %s", api.Name, operationName, err)
		}

		if oper.Error != nil {
			return operations, completed, fmt.Errorf("enabling api: %s, operation %s", api.Name, operationName, err)
		}

		if !oper.Done {
			pending[api] = operationName
		} else {
			completed = append(completed, api)
		}
	}

	return pending, completed, nil
}
