package setup_test

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
	. "omg-cli/omg/setup"

	"omg-cli/google"
	"omg-cli/google/googlefakes"

	"log"

	"os"

	"errors"

	"omg-cli/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GcpProject", func() {
	var (
		logger *log.Logger
	)
	BeforeEach(func() {
		logger = log.New(os.Stdout, "", 0)
	})
	Describe("AdequateQuota", func() {
		It("handles errors", func() {
			service := &googlefakes.FakeQuotaService{ProjectStub: func() (map[string]google.Quota, error) {
				return nil, errors.New("my error")
			}}
			project, err := NewProjectValidator(logger, service, nil, []google.Quota{}, map[string][]google.Quota{}, []google.API{})
			Expect(err).NotTo(HaveOccurred())

			quotaErrors, _, err := project.ValidateQuotas()
			Expect(err).To(HaveOccurred())
			Expect(quotaErrors).To(BeNil())
		})

		It("detects inadequate global quota", func() {
			quotaRequirement := google.Quota{Limit: 2.0, Name: "NETWORKS"}
			quotaActual := google.Quota{Limit: 0, Name: "NETWORKS"}
			service := &googlefakes.FakeQuotaService{ProjectStub: func() (map[string]google.Quota, error) {
				return map[string]google.Quota{
					"NETWORKS": quotaActual,
				}, nil
			}}

			project, err := NewProjectValidator(logger, service, nil, []google.Quota{quotaRequirement}, map[string][]google.Quota{}, []google.API{})
			Expect(err).NotTo(HaveOccurred())

			quotaErrors, satisfied, err := project.ValidateQuotas()
			Expect(err).To(Equal(ErrUnsatisfiedQuota))
			Expect(quotaErrors).ToNot(BeNil())
			Expect(quotaErrors).To(ContainElement(QuotaError{Quota: quotaRequirement, Region: "global"}))
			Expect(satisfied).To(BeEmpty())
		})

		It("detects adequate global quota", func() {
			quotaRequirement := google.Quota{Limit: 2.0, Name: "NETWORKS"}
			quotaActual := google.Quota{Limit: 5.0, Name: "NETWORKS"}
			service := &googlefakes.FakeQuotaService{ProjectStub: func() (map[string]google.Quota, error) {
				return map[string]google.Quota{
					"NETWORKS": quotaActual,
				}, nil
			}}

			project, err := NewProjectValidator(logger, service, nil, []google.Quota{quotaRequirement}, map[string][]google.Quota{}, []google.API{})
			Expect(err).NotTo(HaveOccurred())

			errors, satisfied, err := project.ValidateQuotas()
			Expect(err).NotTo(HaveOccurred())
			Expect(errors).To(BeEmpty())
			Expect(satisfied).To(ContainElement(quotaRequirement))
		})

		It("detects inadequate regional quota", func() {
			quotaRequirement := google.Quota{Name: "CPUS", Limit: 100.0}
			quotaActual := google.Quota{Name: "CPUS", Limit: 10.0}

			service := &googlefakes.FakeQuotaService{RegionStub: func(region string) (map[string]google.Quota, error) {
				Expect(region).To(Equal("us-east1"))
				return map[string]google.Quota{
					"CPUS": quotaActual,
				}, nil
			}}

			project, err := NewProjectValidator(logger, service, nil, []google.Quota{}, map[string][]google.Quota{"us-east1": {quotaRequirement}}, []google.API{})
			Expect(err).NotTo(HaveOccurred())

			errors, satisfied, err := project.ValidateQuotas()
			Expect(err).To(Equal(ErrUnsatisfiedQuota))
			Expect(errors).To(ContainElement(QuotaError{Quota: quotaRequirement, Actual: 10.0, Region: "us-east1"}))
			Expect(satisfied).To(BeEmpty())
		})
	})
	Describe("ProjectQuotaRequirements", func() {
		It("generates requirements", func() {
			Expect(ProjectQuotaRequirements()).NotTo(BeEmpty())
		})
	})
	Describe("RegionalQuotaRequirements", func() {
		It("generates requirements", func() {
			cfg := &config.EnvConfig{Region: "us-west1"}
			req := RegionalQuotaRequirements(cfg)
			Expect(req).To(HaveKey("us-west1"))
			Expect(req["us-west1"]).NotTo(BeEmpty())
		})
	})

	Describe("EnableAPIs", func() {
		It("enables", func() {
			service := &googlefakes.FakeAPIService{}
			requiredApis := []google.API{{"foo"}, {"bar"}}

			project, err := NewProjectValidator(logger, nil, service, []google.Quota{}, map[string][]google.Quota{}, requiredApis)
			Expect(err).NotTo(HaveOccurred())

			project.EnableAPIs()

			Expect(service.EnableCallCount()).To(Equal(1))
			Expect(service.EnableArgsForCall(0)).To(Equal(requiredApis))
		})
	})

	Describe("RequiredAPIs", func() {
		It("generates required APIs", func() {
			Expect(RequiredAPIs()).NotTo(BeEmpty())
		})
	})
})
