package main_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("generate-ert-configuration", func() {
	var (
		providerConfiguration string
		jobInstances          string
	)

	Context("when the provider is AWS", func() {
		BeforeEach(func() {
			providerConfiguration = `{
				"azs": ["my-az1", "my-az2"],
				"network_name": "my-pcf-network",
				"apps_domain": "apps.example.com",
				"sys_domain": "sys.example.com",
				"ert_packages_bucket": "some-packages-bucket",
				"ert_resources_bucket": "some-resources-bucket",
				"ert_droplets_bucket": "some-droplets-bucket",
				"ert_buildpacks_bucket": "some-buildpacks-bucket",
				"smtp_username": "some-smtp-username",
				"smtp_password": "some-smtp-password",
				"ssh_elb_name": "some-ssh-elb-name",
				"tcp_elb_name": "some-tcp-elb-name",
				"web_elb_name": "some-web-elb-name",
				"region": "us-west-1",
				"iam_user_access_key": "some-access-key",
				"iam_user_secret_access_key": "some-access-secret"
			}`

			jobInstances = `{
				"etcd_server": {
					"instances": 3
				},
				"etcd_tls_server": {
					"instances": 3
				},
				"consul_server": {
					"instances": 3
				},
				"diego_cell": {
					"instances": 3
				},
				"diego_database": {
					"instances": 3
				},
				"router": {
					"instances": 3
				},
				"mysql": {
					"instances": 3
				},
				"mysql_proxy": {
					"instances": 2
				}
			}`
		})

		It("outputs a configuration JSON for ERT using ELBs in front of Router, TCP Router, and Diego Brain", func() {
			command := exec.Command(pathToMain,
				"--ops-manager-domain", "pcf.example.com",
				"--provider", "aws",
				"--provider-configuration", providerConfiguration,
				"--ssl-cert", "my-cool-cert",
				"--ssl-private-key", "my-cool-key",
				"--resources", jobInstances,
				"--enable-saml-cert",
			)

			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))

			Expect(session.Out.Contents()).To(MatchJSON(`{
			"product": {
				"ops_manager_domain": "pcf.example.com",
				"properties": {
					".cloud_controller.system_domain": {"value": "sys.example.com"},
					".cloud_controller.apps_domain": {"value": "apps.example.com"},
					".ha_proxy.skip_cert_verify": {"value": true},
					".mysql_monitor.recipient_email": {"value": "cf-release-engineering@pivotal.io"},
					".properties.networking_point_of_entry": {"value": "external_ssl"},
					".properties.networking_point_of_entry.external_ssl.ssl_rsa_certificate": {
						"value": {
							"cert_pem": "my-cool-cert",
							"private_key_pem": "my-cool-key"
						}
					},
					".uaa.service_provider_key_credentials": {
						"value": {
							"cert_pem": "my-cool-cert",
							"private_key_pem": "my-cool-key"
						}
					},
					".properties.security_acknowledgement": {"value":"X"},
					".properties.system_blobstore": {"value":"external"},
					".properties.system_blobstore.external.packages_bucket": {"value":"some-packages-bucket"},
					".properties.system_blobstore.external.droplets_bucket": {"value":"some-droplets-bucket"},
					".properties.system_blobstore.external.resources_bucket": {"value":"some-resources-bucket"},
					".properties.system_blobstore.external.buildpacks_bucket": {"value":"some-buildpacks-bucket"},
					".properties.system_blobstore.external.endpoint": {"value":"https://s3-us-west-1.amazonaws.com"},
					".properties.system_blobstore.external.region": {"value":"us-west-1"},
					".properties.system_blobstore.external.access_key": {"value":"some-access-key"},
					".properties.system_blobstore.external.secret_key": {"value":
					{
						"secret": "some-access-secret"
					}
				},
				".properties.tcp_routing": {"value":"enable"},
				".properties.tcp_routing.enable.reservable_ports": {"value":"1024-1123"},
				".properties.smtp_from":{"value": "identitystaging@mailinator.com"},
				".properties.smtp_address":{"value": "smtp.sendgrid.net"},
				".properties.smtp_port":{"value": "2525"},
				".properties.smtp_credentials":{"value": { "identity": "some-smtp-username", "password": "some-smtp-password" }},
				".properties.smtp_enable_starttls_auto": {"value":true},
				".properties.smtp_auth_mechanism": {"value":"plain"},
				".properties.logger_endpoint_port": {"value": "4443"}
			},
			"network": {
				"singleton_availability_zone": {"name": "my-az1"},
				"other_availability_zones": [{"name": "my-az1"}, {"name": "my-az2"}],
				"network": {"name": "my-pcf-network"}
			},
			"resources": {
				"tcp_router": {
					"elb_names": ["some-tcp-elb-name"]
				},
				"mysql": {
					"instances": 3
				},
				"mysql_proxy": {
					"instances": 2
				},
				"router": {
					"instances": 3,
					"elb_names": ["some-web-elb-name"]
				},
				"consul_server": {
					"instances": 3
				},
				"etcd_server": {
					"instances": 3
				},
				"etcd_tls_server": {
					"instances": 3
				},
				"diego_cell": {
					"instances": 3
				},
				"diego_database": {
					"instances": 3
				},
				"diego_brain": {
					"elb_names": ["some-ssh-elb-name"]
				}
			}
		}}`))
		})
	})

	Context("when the provider is GCP", func() {
		BeforeEach(func() {
			providerConfiguration = `{
			"azs": ["my-az1", "my-az2"],
			"network_name": "my-pcf-network",
			"ert_subnet_name": "my-cf-subnet",
			"http_lb_backend_name": "spearmint",
			"ws_router_pool": "peppermint",
			"ssh_router_pool": "vanilla",
			"tcp_router_pool": "chocolate",
			"apps_domain": "apps.example.com",
			"sys_domain": "sys.example.com",
			"packages_bucket": "some-packages-bucket",
			"resources_bucket": "some-resources-bucket",
			"droplets_bucket": "some-droplets-bucket",
			"buildpacks_bucket": "some-buildpacks-bucket",
			"storage_interop_access_key": "some-access-key",
			"storage_interop_secret_key": "some-secret-key",
			"smtp_username": "some-smtp-username",
			"smtp_password": "some-smtp-password",
			"enable_container_networking": true
		}`

			jobInstances = `{
			"etcd_server": {
				"instances": 1
			},
			"etcd_tls_server": {
				"instances": 3
			},
			"consul_server": {
				"instances": 3
			},
			"diego_cell": {
				"instances": 3
			},
			"diego_database": {
				"instances": 3
			},
			"router": {
				"instances": 3
			},
			"mysql": {
				"instances": 3
			},
			"mysql_proxy": {
				"instances": 2
			}
		}`
		})

		It("outputs a configuration JSON for ERT using ELBs in front of Router, TCP Router, and Diego Brain", func() {
			command := exec.Command(pathToMain,
				"--ops-manager-domain", "pcf.example.com",
				"--provider", "gcp",
				"--provider-configuration", providerConfiguration,
				"--ssl-cert", "my-cool-cert",
				"--ssl-private-key", "my-cool-key",
				"--resources", jobInstances,
			)

			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))

			Expect(session.Out.Contents()).To(MatchJSON(`{
		"product": {
			"ops_manager_domain": "pcf.example.com",
			"properties": {
				".cloud_controller.system_domain": {"value": "sys.example.com"},
				".cloud_controller.apps_domain": {"value": "apps.example.com"},
				".ha_proxy.skip_cert_verify": {"value": true},
				".mysql_monitor.recipient_email": {"value": "cf-release-engineering@pivotal.io"},
				".properties.networking_point_of_entry": {"value": "external_ssl"},
				".properties.networking_point_of_entry.external_ssl.ssl_rsa_certificate": {
					"value": {
						"cert_pem": "my-cool-cert",
						"private_key_pem": "my-cool-key"
					}
				},
				".properties.security_acknowledgement": {"value":"X"},
				".properties.system_blobstore": {"value":"external_gcs"},
				".properties.system_blobstore.external_gcs.packages_bucket": {"value":"some-packages-bucket"},
				".properties.system_blobstore.external_gcs.droplets_bucket": {"value":"some-droplets-bucket"},
				".properties.system_blobstore.external_gcs.resources_bucket": {"value":"some-resources-bucket"},
				".properties.system_blobstore.external_gcs.buildpacks_bucket": {"value":"some-buildpacks-bucket"},
				".properties.system_blobstore.external_gcs.access_key": {"value":"some-access-key"},
				".properties.system_blobstore.external_gcs.secret_key": {"value":
				{
					"secret": "some-secret-key"
				}
			},
			".properties.tcp_routing": {"value":"enable"},
			".properties.tcp_routing.enable.reservable_ports": {"value":"1024-1123"},
			".properties.smtp_from":{"value": "identitystaging@mailinator.com"},
			".properties.smtp_address":{"value": "smtp.sendgrid.net"},
			".properties.smtp_port":{"value": "2525"},
			".properties.smtp_credentials":{"value": { "identity": "some-smtp-username", "password": "some-smtp-password" }},
			".properties.smtp_enable_starttls_auto": {"value":true},
			".properties.smtp_auth_mechanism": {"value":"plain"},
			".properties.container_networking": {"value":"enable"}
		},
		"network": {
			"singleton_availability_zone": {"name": "my-az1"},
			"other_availability_zones": [{"name": "my-az1"}, {"name": "my-az2"}],
			"network": {"name": "my-cf-subnet"}
		},
		"resources": {
			"tcp_router": {
				"elb_names": ["tcp:chocolate"]
			},
			"mysql": {
				"instances": 3
			},
			"mysql_proxy": {
				"instances": 2
			},
			"router": {
				"instances": 3,
				"elb_names": ["http:spearmint", "tcp:peppermint"]
			},
			"consul_server": {
				"instances": 3
			},
			"etcd_server": {
				"instances": 1
			},
			"etcd_tls_server": {
				"instances": 3
			},
			"diego_cell": {
				"instances": 3
			},
			"diego_database": {
				"instances": 3
			},
			"diego_brain": {
				"elb_names": ["tcp:vanilla"]
			}
		}
	}}`))
		})

		Context("when there is an external DB configured", func() {
			BeforeEach(func() {
				providerConfiguration = `{
			"sql_db_ip": "1.2.3.4",
			"ert_sql_username": "sqlusername",
			"ert_sql_password": "sqlpassword",
			"azs": ["my-az1", "my-az2"],
			"network_name": "my-pcf-network",
			"ert_subnet_name": "my-cf-subnet",
			"http_lb_backend_name": "spearmint",
			"ws_router_pool": "peppermint",
			"ssh_router_pool": "vanilla",
			"tcp_router_pool": "chocolate",
			"apps_domain": "apps.example.com",
			"sys_domain": "sys.example.com",
			"packages_bucket": "some-packages-bucket",
			"resources_bucket": "some-resources-bucket",
			"droplets_bucket": "some-droplets-bucket",
			"buildpacks_bucket": "some-buildpacks-bucket",
			"storage_interop_access_key": "some-access-key",
			"storage_interop_secret_key": "some-secret-key",
			"smtp_username": "some-smtp-username",
			"smtp_password": "some-smtp-password"
		}`
			})

			It("outputs the correct ERT configuration", func() {
				command := exec.Command(pathToMain,
					"--ops-manager-domain", "pcf.example.com",
					"--provider", "gcp",
					"--provider-configuration", providerConfiguration,
					"--ssl-cert", "my-cool-cert",
					"--ssl-private-key", "my-cool-key",
				)

				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				Eventually(session).Should(gexec.Exit(0))

				Expect(session.Out.Contents()).To(MatchJSON(`{
		"product": {
			"ops_manager_domain": "pcf.example.com",
			"properties": {
				".properties.system_database": { "value": "external" },
				".properties.system_database.external.port": { "value": 3306 },
				".properties.system_database.external.host": { "value": "1.2.3.4" },
				".properties.system_database.external.app_usage_service_username": { "value": "sqlusername" },
				".properties.system_database.external.app_usage_service_password": { "value":
				{
					"secret": "sqlpassword"
				}
			},
			".properties.system_database.external.autoscale_username": { "value": "sqlusername" },
			".properties.system_database.external.autoscale_password": { "value":
			{
				"secret": "sqlpassword"
			}
		},
		".properties.system_database.external.ccdb_username": { "value": "sqlusername" },
		".properties.system_database.external.ccdb_password": { "value":
		{
			"secret": "sqlpassword"
		}
	},
	".properties.system_database.external.diego_username": { "value": "sqlusername" },
	".properties.system_database.external.diego_password": { "value":
	{
		"secret": "sqlpassword"
	}
},
".properties.system_database.external.networkpolicyserver_username": { "value": "sqlusername" },
".properties.system_database.external.networkpolicyserver_password": { "value":
{
	"secret": "sqlpassword"
}
						},
						".properties.system_database.external.nfsvolume_username": { "value": "sqlusername" },
						".properties.system_database.external.nfsvolume_password": { "value":
						{
							"secret": "sqlpassword"
						}
					},
					".properties.system_database.external.notifications_username": { "value": "sqlusername" },
					".properties.system_database.external.notifications_password": { "value":
					{
						"secret": "sqlpassword"
					}
				},
				".properties.system_database.external.account_username": { "value": "sqlusername" },
				".properties.system_database.external.account_password": { "value":
				{
					"secret": "sqlpassword"
				}
			},
			".properties.system_database.external.routing_username": { "value": "sqlusername" },
			".properties.system_database.external.routing_password": { "value":
			{
				"secret": "sqlpassword"
			}
		},
		".properties.system_database.external.uaa_username": { "value": "sqlusername" },
		".properties.system_database.external.uaa_password": { "value":
		{
			"secret": "sqlpassword"
		}
	},
	".cloud_controller.system_domain": {"value": "sys.example.com"},
	".cloud_controller.apps_domain": {"value": "apps.example.com"},
	".ha_proxy.skip_cert_verify": {"value": true},
	".mysql_monitor.recipient_email": {"value": "cf-release-engineering@pivotal.io"},
	".properties.networking_point_of_entry": {"value": "external_ssl"},
	".properties.networking_point_of_entry.external_ssl.ssl_rsa_certificate": {		
		"value": {		
			"cert_pem": "my-cool-cert",		
			"private_key_pem": "my-cool-key"
		}		
	},
	".properties.security_acknowledgement": {"value":"X"},
	".properties.system_blobstore": {"value":"external_gcs"},
	".properties.system_blobstore.external_gcs.packages_bucket": {"value":"some-packages-bucket"},
	".properties.system_blobstore.external_gcs.droplets_bucket": {"value":"some-droplets-bucket"},
	".properties.system_blobstore.external_gcs.resources_bucket": {"value":"some-resources-bucket"},
	".properties.system_blobstore.external_gcs.buildpacks_bucket": {"value":"some-buildpacks-bucket"},
	".properties.system_blobstore.external_gcs.access_key": {"value":"some-access-key"},
	".properties.system_blobstore.external_gcs.secret_key": {"value":
	{
		"secret": "some-secret-key"
	}
},
".properties.tcp_routing": {"value":"enable"},
".properties.tcp_routing.enable.reservable_ports": {"value":"1024-1123"},
".properties.smtp_from":{"value": "identitystaging@mailinator.com"},
".properties.smtp_address":{"value": "smtp.sendgrid.net"},
".properties.smtp_port":{"value": "2525"},
".properties.smtp_credentials":{"value": { "identity": "some-smtp-username", "password": "some-smtp-password" }},
".properties.smtp_enable_starttls_auto": {"value":true},
".properties.smtp_auth_mechanism": {"value":"plain"}
					},
					"network": {
						"singleton_availability_zone": {"name": "my-az1"},
						"other_availability_zones": [{"name": "my-az1"}, {"name": "my-az2"}],
						"network": {"name": "my-cf-subnet"}
					},
					"resources": {
						"tcp_router": {
							"elb_names": ["tcp:chocolate"]
						},
						"router": {
							"elb_names": ["http:spearmint", "tcp:peppermint"]
						},
						"diego_brain": {
							"elb_names": ["tcp:vanilla"]
						}
					}
				}}`))
			})
		})

		Context("failure cases", func() {
			Context("when given an invalid JSON as provider config", func() {
				It("returns the error from ParseGCPMetadata", func() {
					command := exec.Command(pathToMain,
						"--ops-manager-domain", "pcf.example.com",
						"--provider", "gcp",
						"--provider-configuration", "%%%",
						"--ssl-cert", "my-cool-cert",
						"--ssl-private-key", "my-cool-key",
					)

					session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())
					Eventually(session).Should(gexec.Exit(1))
					Eventually(session.Err).Should(gbytes.Say("could not parse provider configuration"))
				})
			})

			Context("when given an invalid JSON as resource config", func() {
				It("returns the error from PrepareGCPforOm", func() {
					command := exec.Command(pathToMain,
						"--ops-manager-domain", "pcf.example.com",
						"--provider", "gcp",
						"--provider-configuration", providerConfiguration,
						"--resources", "%%%%%%",
						"--ssl-cert", "my-cool-cert",
						"--ssl-private-key", "my-cool-key",
					)

					session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())
					Eventually(session).Should(gexec.Exit(1))
					Eventually(session.Err).Should(gbytes.Say("could not parse resource configuration"))
				})
			})
		})
	})

	Context("when the provider is Azure", func() {
		BeforeEach(func() {
			providerConfiguration = `{
				"apps_domain": "apps.example.com",
				"sys_domain": "sys.example.com",
				"web_lb_name": "some-web-lb",
				"mysql_lb_name": "some-mysql-lb",
				"tcp_lb_name": "some-tcp-lb",
				"ert_subnet_name": "my-ert-subnet",
				"smtp_username": "some-smtp-username",
				"smtp_password": "some-smtp-password",
				"cf_storage_account_name": "some-account-name",
				"cf_storage_account_access_key": "some-key",
				"cf_buildpacks_storage_container": "some-buildpacks",
				"cf_droplets_storage_container": "some-droplets",
				"cf_packages_storage_container": "some-packages",
				"cf_resources_storage_container": "some-resources"
			}`

			jobInstances = `{
				"etcd_server": {
					"instances": 1
				},
				"etcd_tls_server": {
					"instances": 3
				},
				"consul_server": {
					"instances": 3
				},
				"diego_cell": {
					"instances": 3
				},
				"diego_database": {
					"instances": 3
				},
				"ha_proxy": {
					"instances": 3
				},
				"router": {
					"instances": 3
				},
				"mysql": {
					"instances": 3
				},
				"mysql_proxy": {
					"instances": 2
				}
			}`
		})

		It("outputs a configuration JSON for ERT using ELBs in front of HA Proxy, TCP Router, and Mysql Proxy", func() {
			command := exec.Command(pathToMain,
				"--ops-manager-domain", "pcf.example.com",
				"--provider", "azure",
				"--provider-configuration", providerConfiguration,
				"--ssl-cert", "my-cool-cert",
				"--ssl-private-key", "my-cool-key",
				"--resources", jobInstances,
			)

			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))

			Expect(session.Out.Contents()).To(MatchJSON(`{
			"product": {
				"ops_manager_domain": "pcf.example.com",
				"properties": {
					".cloud_controller.system_domain": {"value": "sys.example.com"},
					".cloud_controller.apps_domain": {"value": "apps.example.com"},
					".ha_proxy.skip_cert_verify": {"value": true},
					".mysql_monitor.recipient_email": {"value": "cf-release-engineering@pivotal.io"},
					".properties.networking_point_of_entry": {"value": "haproxy"},
					".properties.networking_point_of_entry.haproxy.ssl_rsa_certificate": {
						"value": {
							"cert_pem": "my-cool-cert",
							"private_key_pem": "my-cool-key"
						}
					},
					".properties.security_acknowledgement": {"value":"X"},
					".properties.tcp_routing": {"value":"enable"},
					".properties.tcp_routing.enable.reservable_ports": {"value":"1024-1173"},
					".properties.smtp_from":{"value": "identitystaging@mailinator.com"},
					".properties.smtp_address":{"value": "smtp.sendgrid.net"},
					".properties.smtp_port":{"value": "2525"},
					".properties.smtp_credentials":{"value": { "identity": "some-smtp-username", "password": "some-smtp-password" } },
					".properties.smtp_enable_starttls_auto": {"value":true},
					".properties.smtp_auth_mechanism": {"value":"plain"},
					".properties.system_blobstore": {"value":"external_azure"},
					".properties.system_blobstore.external_azure.account_name": {"value":"some-account-name"},
					".properties.system_blobstore.external_azure.access_key": {"value":
					{
						"secret": "some-key"
					}
				},
				".properties.system_blobstore.external_azure.buildpacks_container": {"value": "some-buildpacks"},
				".properties.system_blobstore.external_azure.droplets_container": {"value": "some-droplets"},
				".properties.system_blobstore.external_azure.packages_container": {"value": "some-packages"},
				".properties.system_blobstore.external_azure.resources_container": {"value": "some-resources"}
			},
			"network": {
				"singleton_availability_zone": {"name": "null"},
				"other_availability_zones": [{"name": "null"}],
				"network": {"name": "my-ert-subnet"}
			},
			"resources": {
				"tcp_router": {
					"elb_names": ["some-tcp-lb"]
				},
				"ha_proxy": {
					"instances": 3,
					"elb_names": ["some-web-lb"]
				},
				"router": {
					"instances": 3
				},
				"mysql": {
					"instances": 3
				},
				"mysql_proxy": {
					"instances": 2,
					"elb_names": ["some-mysql-lb"]
				},
				"consul_server": {
					"instances": 3
				},
				"etcd_server": {
					"instances": 1
				},
				"etcd_tls_server": {
					"instances": 3
				},
				"diego_cell": {
					"instances": 3
				},
				"diego_database": {
					"instances": 3
				}
			}
		}
	}`))
		})

		Context("when the storage account name is not specified in the provider configuration", func() {
			BeforeEach(func() {
				providerConfiguration = `{
			"apps_domain": "apps.example.com",
			"sys_domain": "sys.example.com",
			"web_lb_name": "some-web-lb",
			"mysql_lb_name": "some-mysql-lb",
			"tcp_lb_name": "some-tcp-lb",
			"ert_subnet_name": "my-ert-subnet",
			"smtp_username": "some-smtp-username",
			"smtp_password": "some-smtp-password"
		}`
			})

			It("does not attempt to configure the external storage option", func() {
				command := exec.Command(pathToMain,
					"--ops-manager-domain", "pcf.example.com",
					"--provider", "azure",
					"--provider-configuration", providerConfiguration,
					"--ssl-cert", "my-cool-cert",
					"--ssl-private-key", "my-cool-key",
					"--resources", jobInstances,
				)

				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				Eventually(session).Should(gexec.Exit(0))

				Expect(session.Out.Contents()).To(MatchJSON(`{
		"product": {
			"ops_manager_domain": "pcf.example.com",
			"properties": {
				".cloud_controller.system_domain": {"value": "sys.example.com"},
				".cloud_controller.apps_domain": {"value": "apps.example.com"},
				".ha_proxy.skip_cert_verify": {"value": true},
				".mysql_monitor.recipient_email": {"value": "cf-release-engineering@pivotal.io"},
				".properties.networking_point_of_entry": {"value": "haproxy"},
				".properties.networking_point_of_entry.haproxy.ssl_rsa_certificate": {
					"value": {
						"cert_pem": "my-cool-cert",
						"private_key_pem": "my-cool-key"
					}
				},
				".properties.security_acknowledgement": {"value":"X"},
				".properties.tcp_routing": {"value":"enable"},
				".properties.tcp_routing.enable.reservable_ports": {"value":"1024-1173"},
				".properties.smtp_from":{"value": "identitystaging@mailinator.com"},
				".properties.smtp_address":{"value": "smtp.sendgrid.net"},
				".properties.smtp_port":{"value": "2525"},
				".properties.smtp_credentials":{"value": { "identity": "some-smtp-username", "password": "some-smtp-password" } },
				".properties.smtp_enable_starttls_auto": {"value":true},
				".properties.smtp_auth_mechanism": {"value":"plain"}
			},
			"network": {
				"singleton_availability_zone": {"name": "null"},
				"other_availability_zones": [{"name": "null"}],
				"network": {"name": "my-ert-subnet"}
			},
			"resources": {
				"tcp_router": {
					"elb_names": ["some-tcp-lb"]
				},
				"ha_proxy": {
					"instances": 3,
					"elb_names": ["some-web-lb"]
				},
				"router": {
					"instances": 3
				},
				"mysql": {
					"instances": 3
				},
				"mysql_proxy": {
					"instances": 2,
					"elb_names": ["some-mysql-lb"]
				},
				"consul_server": {
					"instances": 3
				},
				"etcd_server": {
					"instances": 1
				},
				"etcd_tls_server": {
					"instances": 3
				},
				"diego_cell": {
					"instances": 3
				},
				"diego_database": {
					"instances": 3
				}
			}
		}
	}`))
			})
		})

		Context("failure cases", func() {
			Context("when given an invalid JSON as provider config", func() {
				It("returns the error from ParseAzureMetadata", func() {
					command := exec.Command(pathToMain,
						"--ops-manager-domain", "pcf.example.com",
						"--provider", "azure",
						"--provider-configuration", "%%%",
						"--ssl-cert", "my-cool-cert",
						"--ssl-private-key", "my-cool-key",
					)

					session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())
					Eventually(session).Should(gexec.Exit(1))
					Eventually(session.Err).Should(gbytes.Say("could not parse provider configuration"))
				})
			})

			Context("when given an invalid JSON as resource config", func() {
				It("returns the error from PrepareAzureForOm", func() {
					command := exec.Command(pathToMain,
						"--ops-manager-domain", "pcf.example.com",
						"--provider", "azure",
						"--provider-configuration", providerConfiguration,
						"--resources", "%%%",
						"--ssl-cert", "my-cool-cert",
						"--ssl-private-key", "my-cool-key",
					)

					session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())
					Eventually(session).Should(gexec.Exit(1))
					Eventually(session.Err).Should(gbytes.Say("could not parse resource configuration"))
				})
			})

			Context("when an invalid provider is provided", func() {
				It("returns an error", func() {
					command := exec.Command(pathToMain,
						"--provider", "foobar",
						"--provider-configuration", providerConfiguration,
						"--resources", "%%%",
						"--ssl-cert", "my-cool-cert",
						"--ssl-private-key", "my-cool-key",
					)

					session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())
					Eventually(session).Should(gexec.Exit(1))
					Eventually(session.Err).Should(gbytes.Say("unsupported provider"))
				})
			})
		})
	})

	Context("when the provider is vSphere", func() {
		BeforeEach(func() {
			providerConfiguration = `{
				"azs": ["my-az1", "my-az2"],
				"network_name": "my-pcf-network",
				"apps_domain": "apps.example.com",
				"sys_domain": "sys.example.com",
				"smtp_username": "some-smtp-username",
				"smtp_password": "some-smtp-password"
			}`

			jobInstances = `{
				"etcd_server": {
					"instances": 1
				},
				"consul_server": {
					"instances": 1
				},
				"diego_cell": {
					"instances": 3
				},
				"diego_database": {
					"instances": 1
				},
				"router": {
					"instances": 1
				},
				"mysql": {
					"instances": 1
				},
				"mysql_proxy": {
					"instances": 1
				}
			}`
		})

		It("outputs a configuration JSON for ERT with SSL terminating at HAProxy", func() {
			command := exec.Command(pathToMain,
				"--ops-manager-domain", "pcf.example.com",
				"--provider", "vsphere",
				"--provider-configuration", providerConfiguration,
				"--ssl-cert", "my-cool-cert",
				"--ssl-private-key", "my-cool-key",
				"--resources", jobInstances,
				"--enable-saml-cert",
			)

			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))

			Expect(session.Out.Contents()).To(MatchJSON(`{
			"product": {
				"ops_manager_domain": "pcf.example.com",
				"properties": {
					".cloud_controller.system_domain": {"value": "sys.example.com"},
					".cloud_controller.apps_domain": {"value": "apps.example.com"},
					".ha_proxy.skip_cert_verify": {"value": true},
					".mysql_monitor.recipient_email": {"value": "cf-release-engineering@pivotal.io"},
					".properties.networking_point_of_entry": {"value": "haproxy"},
					".properties.networking_point_of_entry.haproxy.ssl_rsa_certificate": {
						"value": {
							"cert_pem": "my-cool-cert",
							"private_key_pem": "my-cool-key"
						}
					},
					".uaa.service_provider_key_credentials": {
						"value": {
							"cert_pem": "my-cool-cert",
							"private_key_pem": "my-cool-key"
						}
					},
					".properties.security_acknowledgement": {"value":"X"},
					".properties.system_blobstore": {"value":"internal"},
					".properties.tcp_routing": {"value":"disable"},
					".properties.smtp_from":{"value": "identitystaging@mailinator.com"},
					".properties.smtp_address":{"value": "smtp.sendgrid.net"},
					".properties.smtp_port":{"value": "2525"},
					".properties.smtp_credentials":{"value": { "identity": "some-smtp-username", "password": "some-smtp-password" }},
					".properties.smtp_enable_starttls_auto": {"value":true},
					".properties.smtp_auth_mechanism": {"value":"plain"}
				},
				"network": {
					"singleton_availability_zone": {"name": "my-az1"},
					"other_availability_zones": [{"name": "my-az1"}, {"name": "my-az2"}],
					"network": {"name": "my-pcf-network"}
				},
				"resources": {
					"mysql": {
						"instances": 1
					},
					"mysql_proxy": {
						"instances": 1
					},
					"router": {
						"instances": 1
					},
					"consul_server": {
						"instances": 1
					},
					"etcd_server": {
						"instances": 1
					},
					"diego_cell": {
						"instances": 3
					},
					"diego_database": {
						"instances": 1
					}
				}
			}}`))
		})

		Context("failure cases", func() {
			Context("when given an invalid JSON as provider config", func() {
				It("returns the error from ParseVsphereMetadata", func() {
					command := exec.Command(pathToMain,
						"--ops-manager-domain", "pcf.example.com",
						"--provider", "vsphere",
						"--provider-configuration", "%%%",
						"--ssl-cert", "my-cool-cert",
						"--ssl-private-key", "my-cool-key",
					)

					session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())
					Eventually(session).Should(gexec.Exit(1))
					Eventually(session.Err).Should(gbytes.Say("could not parse provider configuration"))
				})
			})
		})
	})
})
