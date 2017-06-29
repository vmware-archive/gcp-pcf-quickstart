package subnet_test

import (
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pcf-releng-ci/tasks/om/generate-bosh-configuration/subnet"
)

var _ = Describe("Subnet", func() {
	Describe("ParseSubnet", func() {
		It("parses the subnet", func() {
			s, err := subnet.ParseSubnet("10.0.0.0/24")
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(Equal(subnet.Subnet{
				IP: []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xff\xff\n\x00\x00\x00"),
				IPNet: &net.IPNet{
					IP:   []byte("\n\x00\x00\x00"),
					Mask: []byte("\xff\xff\xff\x00"),
				},
			}))
		})

		Context("failure cases", func() {
			Context("when the subnet does not have a mask", func() {
				It("returns an error", func() {
					_, err := subnet.ParseSubnet("bread")
					Expect(err.Error()).To(Equal("invalid CIDR address: bread"))
				})
			})
		})
	})

	Describe("Range", func() {
		var s subnet.Subnet

		BeforeEach(func() {
			var err error
			s, err = subnet.ParseSubnet("10.0.16.16/26")
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a range of the subnet", func() {
			r, err := s.Range(1, 10)
			Expect(err).NotTo(HaveOccurred())
			Expect(r).To(Equal("10.0.16.17-10.0.16.26"))
		})

		Context("failure cases", func() {
			Context("when the start value is negative", func() {
				It("returns an error", func() {
					_, err := s.Range(-1, 10)
					Expect(err).To(MatchError("subnet range start \"-1\" cannot be negative"))
				})
			})

			Context("when the end value is greater than 255", func() {
				It("returns an error", func() {
					_, err := s.Range(1, 258)
					Expect(err).To(MatchError("subnet range end \"258\" cannot exceed 255"))
				})
			})

			Context("when the start value is greater than the end value", func() {
				It("returns an error", func() {
					_, err := s.Range(250, 245)
					Expect(err).To(MatchError("subnet range start \"250\" cannot exceed subnet range end \"245\""))
				})
			})
		})
	})

	Describe("IPAddress", func() {
		var s subnet.Subnet

		BeforeEach(func() {
			var err error
			s, err = subnet.ParseSubnet("10.0.0.0/24")
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns the address in the provided subnet", func() {
			r, err := s.IPAddress(2)
			Expect(err).NotTo(HaveOccurred())
			Expect(r).To(Equal("10.0.0.2"))
		})

		Context("failure cases", func() {
			Context("when the value is negative", func() {
				It("returns an error", func() {
					_, err := s.IPAddress(-1)
					Expect(err).To(MatchError("IP address containing \"-1\" is invalid, cannot be negative"))
				})
			})

			Context("when the value is greater than 255", func() {
				It("returns an error", func() {
					_, err := s.IPAddress(258)
					Expect(err).To(MatchError("IP address containing \"258\" is invalid, cannot exceed 255"))
				})
			})
		})
	})
})
