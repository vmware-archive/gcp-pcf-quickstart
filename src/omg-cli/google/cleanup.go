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

package google

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"omg-cli/version"
	"strings"
	"time"

	compute "google.golang.org/api/compute/v1"
)

type CleanupService interface {
	DeleteVM(...VMFilter) (int, error)
}

type VMFilter func(*vmFilter)

type vmFilter struct {
	tag        string
	namePrefix string
	subnet     string
}

func WithTag(tag string) VMFilter {
	return func(opt *vmFilter) {
		opt.tag = tag
	}
}

func WithNamePrefix(prefix string) VMFilter {
	return func(opt *vmFilter) {
		opt.namePrefix = prefix
	}
}

func WithSubNetwork(subnet string) VMFilter {
	return func(opt *vmFilter) {
		opt.subnet = subnet
	}
}

type cleanupService struct {
	logger         *log.Logger
	projectId      string
	computeService *compute.Service
	dryRun         bool
}

func NewCleanupService(logger *log.Logger, projectId string, client *http.Client, dryRun bool) (CleanupService, error) {
	if logger == nil {
		return nil, errors.New("missing logger")
	}

	computeService, err := compute.New(client)
	if err != nil {
		return nil, err
	}
	computeService.UserAgent = version.UserAgent()

	return &cleanupService{logger, projectId, computeService, dryRun}, nil
}

func buildFilter(filter vmFilter) string {
	if filter.namePrefix != "" {
		return fmt.Sprintf("name eq %s", filter.namePrefix)
	}

	return ""
}

type vm struct {
	name, zone string
}

func contains(needle string, haystack []string) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}
	return false
}

func (cs *cleanupService) findVMs(opts ...VMFilter) ([]vm, error) {
	listCall := cs.computeService.Instances.AggregatedList(cs.projectId)
	filter := vmFilter{}
	for _, opt := range opts {
		opt(&filter)
	}

	listCall.Filter(buildFilter(filter))
	var targets []vm
	err := listCall.Pages(context.Background(), func(page *compute.InstanceAggregatedList) error {
		for _, list := range page.Items {
			for _, instance := range list.Instances {
				if filter.tag != "" && !contains(filter.tag, instance.Tags.Items) {
					continue
				}

				if filter.subnet != "" {
					foundSubnet := false
					for _, nic := range instance.NetworkInterfaces {
						if strings.HasSuffix(nic.Subnetwork, filter.subnet) {
							foundSubnet = true
						}
					}

					if foundSubnet == false {
						continue
					}
				}

				zoneParts := strings.Split(instance.Zone, "/")
				target := vm{instance.Name, zoneParts[len(zoneParts)-1]}
				cs.logger.Printf("found VM for deletion: %#v", target)
				targets = append(targets, target)
			}
		}

		return nil
	})

	return targets, err
}

func (cs *cleanupService) deleteVMs(targets []vm) (operationsMap, error) {
	operations := operationsMap{}
	for _, vm := range targets {
		call := cs.computeService.Instances.Delete(cs.projectId, vm.zone, vm.name)
		oper, err := call.Do()
		if err != nil {
			return operations, err
		}
		if oper != nil {
			operations[vm] = oper.Name
		}
	}

	return operations, nil
}

func (cs *cleanupService) waitOnOperation(operations operationsMap) (completed []vm, errs []error) {
	timeout := time.After(time.Duration(5) * time.Minute)
	for len(operations) != 0 {
		select {
		case <-timeout:
			errs = append(errs, errors.New("timeout waiting for VM delete operation"))
			break
		default:
			cs.logger.Printf("waiting for %d VM deletion operation(s) to complete, encountered %d error(s)", len(operations), len(errs))
			time.Sleep(time.Duration(10) * time.Second)

			var roundErrs []error
			var roundCompleted []vm
			operations, roundCompleted, roundErrs = cs.filterCompleted(operations)

			completed = append(completed, roundCompleted...)
			errs = append(errs, roundErrs...)
		}
	}
	return
}

func (cs *cleanupService) DeleteVM(opts ...VMFilter) (int, error) {
	targets, err := cs.findVMs(opts...)
	if err != nil {
		return 0, fmt.Errorf("finding VMs: %v", err)
	}

	if cs.dryRun {
		cs.logger.Printf("dry-run: exiting without performing delete")
		return 0, nil
	}

	operations, err := cs.deleteVMs(targets)
	if err != nil {
		cs.logger.Printf("error deleting VMs: %v", err)
	}

	completed, errs := cs.waitOnOperation(operations)

	for _, operErr := range errs {
		if err == nil {
			err = operErr
		} else {
			err = fmt.Errorf("%v, %v", err, operErr)
		}
	}

	return len(completed), err
}

type operationsMap map[vm]string

const done = "DONE"

func (cs *cleanupService) filterCompleted(operations operationsMap) (pending operationsMap, completed []vm, errors []error) {
	pending = operationsMap{}
	for vm, operationName := range operations {
		oper, err := cs.computeService.ZoneOperations.Get(cs.projectId, vm.zone, operationName).Do()

		if err != nil {
			errors = append(errors, fmt.Errorf("fetching operation %s for vm %#v: %v", operationName, vm, err))
			pending[vm] = operationName
		} else if oper.Error != nil {
			errors = append(errors, fmt.Errorf("operation %s for vm %#v: %v", operationName, vm, oper.Error))
		} else if oper.Status == done {
			completed = append(completed, vm)
		} else {
			pending[vm] = operationName
		}
	}
	return
}
