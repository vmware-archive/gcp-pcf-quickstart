package provider_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pcf-releng-ci/tasks/om/generate-ert-configuration/provider"
)

var _ = Describe("GCPMetadata", func() {
	var metadata string

	BeforeEach(func() {
		metadata = `{
			"azs": ["us-central1-a","us-central1-b"],
			"network_name": "vine-whale-pcf-network",
			"ert_subnet_name": "vine-whale-cf-subnet",
			"project": "cf-release-engineering",
			"region": "us-central1",
			"service_account_key": "{\"key\": \"value\"}"
		}`
	})

	Describe("ParseMetadata", func() {
		It("parses GCP metadata", func() {
			parsedGCPMetadata, err := provider.ParseGCPMetadata(metadata)

			Expect(err).NotTo(HaveOccurred())
			Expect(parsedGCPMetadata).To(Equal(provider.GCPMetadata{
				AZs:        []string{"us-central1-a", "us-central1-b"},
				Network:    "vine-whale-pcf-network",
				ERTSubnet:  "vine-whale-cf-subnet",
				Project:    "cf-release-engineering",
				Region:     "us-central1",
				ServiceKey: "{\"key\": \"value\"}",
			}))
		})

		Context("failure cases", func() {
			Context("when no AZs are provided", func() {
				It("returns an error", func() {
					metadata = `{
					"azs": [],
					"network_name": "vine-whale-pcf-network",
					"ert_subnet_name": "vine-whale-cf-subnet",
					"project": "cf-release-engineering",
					"region": "us-central1",
					"service_account_key": "{\"key\": \"value\"}"
				}`
					_, err := provider.ParseGCPMetadata(metadata)
					Expect(err).To(MatchError(ContainSubstring("expected at least one AZ")))
				})
			})
		})
	})
})
