package provider_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pcf-releng-ci/tasks/om/generate-ert-configuration/provider"
)

var _ = Describe("VSphereMetadata", func() {
	var metadata string

	BeforeEach(func() {
		metadata = `{
			"azs": ["the-first-az","the-second-az"],
			"network_name": "the-first-network",
			"sys_domain": "sys.example.com",
			"apps_domain": "apps.example.com",
			"smtp_username": "postmaster",
			"smtp_password": "the-smtp-password"

		}`
	})

	Describe("ParseMetadata", func() {
		It("parses VSphere metadata", func() {
			parsedVSphereMetadata, err := provider.ParseVSphereMetadata(metadata)

			Expect(err).NotTo(HaveOccurred())
			Expect(parsedVSphereMetadata).To(Equal(provider.VSphereMetadata{
				AZs:          []string{"the-first-az", "the-second-az"},
				Network:      "the-first-network",
				SysDomain:    "sys.example.com",
				AppsDomain:   "apps.example.com",
				SMTPUsername: "postmaster",
				SMTPPassword: "the-smtp-password",
			}))
		})
	})
})
