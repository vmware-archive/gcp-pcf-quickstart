package provider_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pcf-releng-ci/tasks/om/generate-ert-configuration/provider"
)

var _ = Describe("AzureMetadata", func() {
	var metadata string

	BeforeEach(func() {
		metadata = `{
			"network_name": "foo-pcf-network",
			"ert_subnet_name": "foo-ert-subnet",
			"tenant_id": "foo-tenant",
			"sys_domain": "foo-sys-domain",
			"apps_domain": "foo-app-domain",
			"web_lb_name": "web-lb-name",
			"mysql_dns": "mysql-dns-name",
			"mysql_lb_name": "mysql-lb-name",
			"tcp_lb_name": "tcp-routing-lb-name"
		}`
	})

	It("parses Azure metadata", func() {
		parsedAzureMetadata, err := provider.ParseAzureMetadata(metadata)
		Expect(err).NotTo(HaveOccurred())

		Expect(parsedAzureMetadata).To(Equal(provider.AzureMetadata{
			Network:     "foo-pcf-network",
			ERTSubnet:   "foo-ert-subnet",
			SysDomain:   "foo-sys-domain",
			AppsDomain:  "foo-app-domain",
			WebLBName:   "web-lb-name",
			MySQLDNS:    "mysql-dns-name",
			MySQLLBName: "mysql-lb-name",
			TCPLBName:   "tcp-routing-lb-name",
		}))
	})
})
