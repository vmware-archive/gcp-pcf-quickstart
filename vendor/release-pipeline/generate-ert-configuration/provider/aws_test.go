package provider_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pcf-releng-ci/tasks/om/generate-ert-configuration/provider"
)

var _ = Describe("AWSMetadata", func() {
	var metadata string

	BeforeEach(func() {
		metadata = `{
			"azs": ["us-east-1a","us-east-1b"],
			"network_name": "vine-whale-pcf-network",
			"ssh_elb_name": "some-ssh-elb-name",
			"tcp_elb_name": "some-tcp-elb-name",
			"web_elb_name": "some-web-elb-name",
			"region": "us-west-1",
			"iam_user_access_key": "some-access-key",
			"iam_user_secret_access_key": "some-access-secret"
		}`
	})

	Describe("ParseMetadata", func() {
		It("parses AWS metadata", func() {
			parsedAWSMetadata, err := provider.ParseAWSMetadata(metadata)

			Expect(err).NotTo(HaveOccurred())
			Expect(parsedAWSMetadata).To(Equal(provider.AWSMetadata{
				AZs:       []string{"us-east-1a", "us-east-1b"},
				Network:   "vine-whale-pcf-network",
				SSHELB:    "some-ssh-elb-name",
				TCPELB:    "some-tcp-elb-name",
				WebELB:    "some-web-elb-name",
				AccessKey: "some-access-key",
				SecretKey: "some-access-secret",
				Region:    "us-west-1",
			}))
		})

		Context("failure cases", func() {
			Context("when no AZs are provided", func() {
				It("returns an error", func() {
					metadata = `{
					"azs": [],
					"network_name": "vine-whale-pcf-network"
				}`
					_, err := provider.ParseAWSMetadata(metadata)
					Expect(err).To(MatchError(ContainSubstring("expected at least one AZ")))
				})
			})
		})
	})
})
