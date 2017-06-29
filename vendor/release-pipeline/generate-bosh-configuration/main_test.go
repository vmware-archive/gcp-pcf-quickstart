package main_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("generate-bosh-configuration", func() {
	var providerConfiguration string

	Context("when the provider is GCP", func() {
		BeforeEach(func() {
			providerConfiguration = `{
				"azs": ["us-central1-a", "us-central1-b"],
				"network_name": "vine-whale-pcf-network",
				"management_subnet_cidrs": ["10.0.0.0/24"],
				"management_subnet_name": "vine-whale-om-subnet",
				"management_subnet_gateway": "10.0.0.1",
				"ert_subnet_cidrs": ["10.0.1.0/24"],
				"ert_subnet_name": "vine-whale-ert-subnet",
				"ert_subnet_gateway": "10.0.1.1",
				"services_subnet_cidrs": ["10.0.2.0/24"],
				"services_subnet_name": "vine-whale-services-subnet",
				"services_subnet_gateway": "10.0.2.1",
				"project": "cf-release-engineering",
				"region": "us-central1",
				"service_account_key": "foo"
			}`
		})

		It("outputs a configuration JSON for bosh", func() {
			command := exec.Command(pathToMain,
				"--provider", "gcp",
				"--provider-configuration", providerConfiguration,
				"--env-name", "vine-whale",
				"--compilation-vm-type", "xlarge.foo",
			)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))
			Expect(session.Out.Contents()).To(MatchJSON(`{
				"az_configuration": {
					"availability_zones": [
						{ "name": "us-central1-a" },
						{ "name": "us-central1-b" }
					]
			  },
				"iaas_configuration": {
					"project": "cf-release-engineering",
					"default_deployment_tag": "vine-whale-vms",
					"auth_json": "foo"
				},
				"director_configuration": {
				  "ntp_servers_string": "169.254.169.254"
				},
				"networks_configuration": {
					"icmp_checks_enabled": false,
					"networks": [
						{
							"name": "vine-whale-om-subnet",
							"service_network": false,
							"subnets": [
								{
									"iaas_identifier": "vine-whale-pcf-network/vine-whale-om-subnet/us-central1",
									"cidr": "10.0.0.0/24",
									"reserved_ip_ranges": "10.0.0.0-10.0.0.4",
									"dns": "8.8.8.8",
									"gateway": "10.0.0.1",
									"availability_zones": [ "us-central1-a","us-central1-b" ]
								}
							]
						},
						{
							"name": "vine-whale-ert-subnet",
							"service_network": false,
							"subnets": [
								{
									"iaas_identifier": "vine-whale-pcf-network/vine-whale-ert-subnet/us-central1",
									"cidr": "10.0.1.0/24",
									"reserved_ip_ranges": "10.0.1.0-10.0.1.4",
									"dns": "8.8.8.8",
									"gateway": "10.0.1.1",
									"availability_zones": [ "us-central1-a","us-central1-b" ]
								}
							]
						},
						{
							"name": "vine-whale-services-subnet",
							"service_network": true,
							"subnets": [
								{
									"iaas_identifier": "vine-whale-pcf-network/vine-whale-services-subnet/us-central1",
									"cidr": "10.0.2.0/24",
									"reserved_ip_ranges": "10.0.2.0-10.0.2.3",
									"dns": "8.8.8.8",
									"gateway": "10.0.2.1",
									"availability_zones": [ "us-central1-a","us-central1-b" ]
								}
							]
						}
					]
				},
				"network_assignment": {
					"singleton_availability_zone": "us-central1-a",
					"network": "vine-whale-om-subnet"
				},
				"resource_configuration": {
					"compilation": {
						"instance_type": {
							"id": "xlarge.foo"
						}
					}
				}
			}`))
		})
	})

	Context("when the provider is AWS", func() {
		BeforeEach(func() {
			providerConfiguration = `{
					"azs": ["us-west-1a", "us-west-1b"],
					"management_subnet_availability_zones": ["us-west-1a"],
					"management_subnet_cidrs": ["10.0.16.0/20"],
					"management_subnet_ids": ["director-subnet-1"],
					"ert_subnet_availability_zones": ["us-west-1a", "us-west-1b"],
					"ert_subnet_cidrs": ["10.0.32.0/20","10.0.64.0/20"],
					"ert_subnet_ids": ["ert-subnet-1", "ert-subnet-2"],
					"services_subnet_availability_zones": ["us-west-1a"],
					"services_subnet_cidrs": ["10.0.48.0/20"],
					"services_subnet_ids": ["service-subnet-1"],
					"ops_manager_ssh_private_key": "fake-private-key",
					"ops_manager_ssh_public_key_name": "fake-public-key-name",
					"region": "fake-region",
					"vms_security_group_id": "fake-security-group-id",
					"vpc_id": "fake-vpc-id",
					"iam_user_access_key": "some_access_key",
					"iam_user_secret_access_key": "some_secret_access_key"
				}`
		})

		It("outputs a configuration JSON for bosh", func() {
			command := exec.Command(pathToMain,
				"--provider", "aws",
				"--provider-configuration", providerConfiguration,
				"--env-name", "vine-whale",
			)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))
			Expect(session.Out.Contents()).To(MatchJSON(`{
					"iaas_configuration": {
						"access_key_id": "some_access_key",
						"secret_access_key": "some_secret_access_key",
						"vpc_id": "fake-vpc-id",
						"security_group": "fake-security-group-id",
						"key_pair_name": "fake-public-key-name",
						"ssh_private_key": "fake-private-key",
						"region": "fake-region",
						"encrypted": false
					},
					"director_configuration": {
					  "ntp_servers_string": "0.amazon.pool.ntp.org"
					},
					"az_configuration": {
						"availability_zones": [
							{ "name": "us-west-1a" },
							{ "name": "us-west-1b" }
						]
				  },
					"networks_configuration":{
						"icmp_checks_enabled": false,
						"networks": [
							{
								"name": "vine-whale-management-network",
								"service_network": false,
								"subnets": [
									{
										"iaas_identifier": "director-subnet-1",
										"cidr": "10.0.16.0/20",
										"dns": "169.254.169.253",
										"reserved_ip_ranges": "10.0.16.0-10.0.16.4",
										"gateway": "10.0.16.1",
										"availability_zones": [ "us-west-1a" ]
									}
								]
							},
							{
								"name": "vine-whale-ert-network",
								"service_network": false,
								"subnets": [
									{
								    "iaas_identifier": "ert-subnet-1",
										"cidr": "10.0.32.0/20",
										"dns": "169.254.169.253",
										"reserved_ip_ranges": "10.0.32.0-10.0.32.4",
										"gateway": "10.0.32.1",
										"availability_zones": [ "us-west-1a" ]
									},
									{
								    "iaas_identifier": "ert-subnet-2",
										"cidr": "10.0.64.0/20",
										"dns": "169.254.169.253",
										"reserved_ip_ranges": "10.0.64.0-10.0.64.4",
										"gateway": "10.0.64.1",
										"availability_zones": [ "us-west-1b" ]
									}
								]
							},
							{
								"name": "vine-whale-services-network",
								"service_network": true,
								"subnets": [
									{
										"iaas_identifier": "service-subnet-1",
										"cidr": "10.0.48.0/20",
										"dns": "169.254.169.253",
										"reserved_ip_ranges": "10.0.48.0-10.0.48.3",
										"gateway": "10.0.48.1",
										"availability_zones": [ "us-west-1a" ]
									}
								]
							}
						]
					},
					"network_assignment": {
						"singleton_availability_zone": "us-west-1a",
						"network": "vine-whale-management-network"
				  }
				}`))
		})
	})

	Context("when the provider is Azure", func() {
		BeforeEach(func() {
			providerConfiguration = `{
					"bosh_root_storage_account": "bosh-storage-account",
					"wildcard_vm_storage_account": "bosh-vms-account",
					"client_id": "some-client-id",
					"client_secret": "some-client-secret",
					"env_dns_zone_name_servers": [
						"ns1-07.azure-dns.com.",
						"ns3-07.azure-dns.org.",
						"ns2-07.azure-dns.net.",
						"ns4-07.azure-dns.info."
					],
					"ops_manager_dns": "pcf.navy-spur.azure.releng.cf-app.com",
					"ops_manager_public_ip": "111.111.111.111",
					"ops_manager_security_group_name": "navy-spur-ops-manager-security-group",
					"ops_manager_ssh_private_key": "---SOME PRIVATE KEY---",
					"ops_manager_ssh_public_key": "ssh-rsa some-public-key",
					"ops_manager_storage_account": "ops-manager-account",
					"pcf_resource_group_name": "navy-spur-pcf-resource-group",
					"management_subnet_cidrs": ["10.0.0.0/24"],
					"management_subnet_name": "navy-spur-om-subnet",
					"management_subnet_gateway": "om-gateway",
					"ert_subnet_cidrs": ["10.0.1.0/24"],
					"ert_subnet_name": "navy-spur-ert-subnet",
					"ert_subnet_gateway": "ert-gateway",
					"services_subnet_cidrs": ["10.0.2.0/24"],
					"services_subnet_name": "navy-spur-services-subnet",
					"services_subnet_gateway": "services-gateway",
					"network_name": "navy-spur-virtual-network",
					"subscription_id": "some-subscription-id",
					"tenant_id": "some-tenant-id"
			}`
		})

		It("outputs a configuration JSON for bosh", func() {
			command := exec.Command(pathToMain,
				"--provider", "azure",
				"--provider-configuration", providerConfiguration,
				"--env-name", "navy-spur",
			)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))
			Expect(session.Out.Contents()).To(MatchJSON(`{
					"iaas_configuration": {
						"subscription_id": "some-subscription-id",
						"tenant_id": "some-tenant-id",
						"client_id": "some-client-id",
						"client_secret": "some-client-secret",
						"resource_group_name": "navy-spur-pcf-resource-group",
						"bosh_storage_account_name": "bosh-storage-account",
						"default_security_group": "navy-spur-ops-manager-security-group",
						"ssh_public_key": "ssh-rsa some-public-key",
						"ssh_private_key": "---SOME PRIVATE KEY---",
						"deployments_storage_account_name": "bosh-vms-account"
					},
					"director_configuration": {
						"ntp_servers_string": "us.pool.ntp.org"
					},
					"networks_configuration": {
						"icmp_checks_enabled": false,
						"networks": [
							{
								"name": "navy-spur-om-subnet",
								"service_network": false,
								"subnets": [
									{
										"iaas_identifier": "navy-spur-virtual-network/navy-spur-om-subnet",
										"cidr": "10.0.0.0/24",
										"reserved_ip_ranges": "10.0.0.0-10.0.0.5",
										"dns": "8.8.8.8",
										"gateway": "om-gateway"
									}
								]
							},
							{
								"name": "navy-spur-ert-subnet",
								"service_network": false,
								"subnets": [
									{
										"iaas_identifier": "navy-spur-virtual-network/navy-spur-ert-subnet",
										"cidr": "10.0.1.0/24",
										"reserved_ip_ranges": "10.0.1.0-10.0.1.4",
										"dns": "8.8.8.8",
										"gateway": "ert-gateway"
									}
								]
							},
							{
								"name": "navy-spur-services-subnet",
								"service_network": true,
								"subnets": [
									{
										"iaas_identifier": "navy-spur-virtual-network/navy-spur-services-subnet",
										"cidr": "10.0.2.0/24",
										"reserved_ip_ranges": "10.0.2.0-10.0.2.3",
										"dns": "8.8.8.8",
										"gateway": "services-gateway"
									}
								]
							}
						]
					},
					"network_assignment": {
						"network": "navy-spur-om-subnet"
					}
				}`))
		})
	})

	Context("error cases", func() {
		Context("when an invalid provider is specified", func() {
			It("throws an error", func() {
				command := exec.Command(pathToMain,
					"--provider", "foo",
					"--provider-configuration", "{}",
				)
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				Eventually(session).Should(gexec.Exit(1))
				Eventually(session.Err).Should(gbytes.Say("invalid provider: foo"))
			})
		})

		Context("when provided invalid json in the provider configuration", func() {
			It("throws an error", func() {
				command := exec.Command(pathToMain,
					"--provider", "gcp",
					"--provider-configuration", "%%%",
					"--env-name", "vine-whale",
				)
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				Eventually(session).Should(gexec.Exit(1))
				Eventually(session.Err).Should(gbytes.Say("invalid character"))
			})
		})

		Context("when provided an invalid subnet CIDR", func() {
			It("throws an error", func() {
				providerConfiguration = `{
					"azs": ["us-central1-a", "us-central1-b"],
					"network_name": "vine-whale-pcf-network",
					"management_subnet_cidrs": ["om-cidr"],
					"management_subnet_name": "vine-whale-om-subnet",
					"management_subnet_gateway": "10.0.0.1",
					"ert_subnet_cidrs": ["10.0.1.0/24"],
					"ert_subnet_name": "vine-whale-ert-subnet",
					"ert_subnet_gateway": "10.0.1.1",
					"services_subnet_cidrs": ["10.0.2.0/24"],
					"services_subnet_name": "vine-whale-services-subnet",
					"services_subnet_gateway": "10.0.2.1",
					"project": "cf-release-engineering",
					"region": "us-central1",
					"service_account_key": "foo"
				}`
				command := exec.Command(pathToMain,
					"--provider", "gcp",
					"--provider-configuration", providerConfiguration,
					"--env-name", "vine-whale",
				)
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				Eventually(session).Should(gexec.Exit(1))
				Eventually(session.Err).Should(gbytes.Say(`invalid CIDR address: om-cidr`))
			})
		})
	})
})
