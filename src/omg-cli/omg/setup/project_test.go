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
			service := &googlefakes.FakeProjectService{QuotasStub: func() (map[string]google.Quota, error) {
				return nil, errors.New("My Error")
			}}
			project, err := NewProjectValiadtor(logger, service, []google.Quota{})
			Expect(err).NotTo(HaveOccurred())

			quotaErrors, _, err := project.EnsureQuota()
			Expect(err).To(HaveOccurred())
			Expect(quotaErrors).To(BeNil())
		})

		It("detects inadequate quota", func() {
			quotaRequirement := google.Quota{Limit: 2.0, Name: "NETWORKS"}
			quotaActual := google.Quota{Limit: 0, Name: "NETWORKS"}
			service := &googlefakes.FakeProjectService{QuotasStub: func() (map[string]google.Quota, error) {
				return map[string]google.Quota{
					"NETWORKS": quotaActual,
				}, nil
			}}

			project, err := NewProjectValiadtor(logger, service, []google.Quota{quotaRequirement})
			Expect(err).NotTo(HaveOccurred())

			quotaErrors, satisfied, err := project.EnsureQuota()
			Expect(err).To(Equal(UnsatisfiedQuotaErr))
			Expect(quotaErrors).ToNot(BeNil())
			Expect(quotaErrors).To(ContainElement(QuotaError{quotaRequirement, 0.0}))
			Expect(satisfied).To(BeEmpty())
		})

		It("detects adequate quota", func() {
			quotaRequirement := google.Quota{Limit: 2.0, Name: "NETWORKS"}
			quotaActual := google.Quota{Limit: 5.0, Name: "NETWORKS"}
			service := &googlefakes.FakeProjectService{QuotasStub: func() (map[string]google.Quota, error) {
				return map[string]google.Quota{
					"NETWORKS": quotaActual,
				}, nil
			}}

			project, err := NewProjectValiadtor(logger, service, []google.Quota{quotaRequirement})
			Expect(err).NotTo(HaveOccurred())

			errors, satisfied, err := project.EnsureQuota()
			Expect(err).NotTo(HaveOccurred())
			Expect(errors).To(BeEmpty())
			Expect(satisfied).To(ContainElement(quotaRequirement))
		})
	})
})
