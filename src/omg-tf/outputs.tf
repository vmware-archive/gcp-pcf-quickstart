output "ops_manager_dns" {
  value = "opsman.${var.dns_suffix}"
}

output "ops_manager_instance_name" {
  value = "${var.opsman_external_ip != "" ? element(concat(google_compute_instance.ops-manager-external.*.name, list("")), 0) : element(concat(google_compute_instance.ops-manager-internal.*.name, list("")), 0)}"
}

output "ops_manager_instance_zone" {
  value = "${var.opsman_external_ip != "" ? element(concat(google_compute_instance.ops-manager-external.*.zone, list("")), 0) : element(concat(google_compute_instance.ops-manager-internal.*.zone, list("")), 0)}"
}

output "sys_domain" {
  value = "sys.${var.dns_suffix}"
}

output "apps_domain" {
  value = "apps.${var.dns_suffix}"
}

output "tcp_domain" {
  value = "tcp.${var.dns_suffix}"
}

output "doppler_domain" {
  value = "doppler.sys.${var.dns_suffix}"
}

output "loggregator_domain" {
  value = "loggregator.sys.${var.dns_suffix}"
}

output "dns_suffix" {
  value = "${var.dns_suffix}"
}

output "project" {
  value = "${var.project}"
}

output "region" {
  value = "${var.region}"
}

output "azs" {
  value = "${var.zones}"
}

output "vm_tag" {
  value = "${var.no_ip_instance_tag}"
}

output "network_name" {
  value = "${google_compute_network.pcf-network.name}"
}

output "sql_db_ip" {
  value = "${module.external_database.ip}"
}

output "sql_db_port" {
  value = "${module.external_database.sql_db_port}"
}

output "management_subnet_gateway" {
  value = "${google_compute_subnetwork.management-subnet.gateway_address}"
}

output "management_subnet_cidrs" {
  value = ["${google_compute_subnetwork.management-subnet.ip_cidr_range}"]
}

output "management_subnet_name" {
  value = "${google_compute_subnetwork.management-subnet.name}"
}

output "opsman_sql_db_name" {
  value = "${module.external_database.opsman_sql_db_name}"
}

output "ert_subnet_gateway" {
  value = "${google_compute_subnetwork.ert-subnet.gateway_address}"
}

output "ert_subnet_cidrs" {
  value = ["${google_compute_subnetwork.ert-subnet.ip_cidr_range}"]
}

output "ert_subnet_name" {
  value = "${google_compute_subnetwork.ert-subnet.name}"
}

output "services_subnet_gateway" {
  value = "${google_compute_subnetwork.services-subnet.gateway_address}"
}

output "services_subnet_cidrs" {
  value = ["${google_compute_subnetwork.services-subnet.ip_cidr_range}"]
}

output "services_subnet_name" {
  value = "${google_compute_subnetwork.services-subnet.name}"
}

output "dynamic_services_subnet_gateway" {
  value = "${google_compute_subnetwork.dynamic-services-subnet.gateway_address}"
}

output "dynamic_services_subnet_cidrs" {
  value = ["${google_compute_subnetwork.dynamic-services-subnet.ip_cidr_range}"]
}

output "dynamic_services_subnet_name" {
  value = "${google_compute_subnetwork.dynamic-services-subnet.name}"
}

output "http_lb_backend_name" {
  value = "${google_compute_backend_service.http_lb_backend_service.name}"
}

output "isoseg_lb_backend_name" {
  value = "${module.isolation_segment.load_balancer_name}"
}

output "ssh_router_pool" {
  value = "${google_compute_target_pool.cf-ssh.name}"
}

output "wss_router_pool" {
  value = "${google_compute_target_pool.cf-wss.name}"
}

output "tcp_router_pool" {
  value = "${google_compute_target_pool.cf-tcp.name}"
}

output "buildpacks_bucket" {
  value = "${google_storage_bucket.buildpacks.name}"
}

output "droplets_bucket" {
  value = "${google_storage_bucket.droplets.name}"
}

output "packages_bucket" {
  value = "${google_storage_bucket.packages.name}"
}

output "resources_bucket" {
  value = "${google_storage_bucket.resources.name}"
}

output "director_blobstore_bucket" {
  value = "${google_storage_bucket.director.name}"
}

output "ert_sql_username" {
  value = "${module.external_database.ert_sql_username}"
}

output "ert_sql_password" {
  value     = "${module.external_database.ert_sql_password}"
  sensitive = true
}

output "opsman_sql_username" {
  value = "${module.external_database.opsman_sql_username}"
}

output "opsman_sql_password" {
  value     = "${module.external_database.opsman_sql_password}"
  sensitive = true
}

output "jumpbox_public_ip" {
  value = "${google_compute_address.jumpbox.address}"
}

output "ssl_cert" {
  value = "${var.ssl_cert}"
}

output "ssl_cert_private_key" {
  value     = "${var.ssl_cert_private_key}"
  sensitive = true
}

output "tcp_port_range" {
  value = "${var.tcp_port_range}"
}

output "stackdriver_service_account_key_base64" {
  value     = "${google_service_account_key.stackdrvier_nozzle.private_key}"
  sensitive = true
}

output "service_broker_service_account_key_base64" {
  value     = "${google_service_account_key.service_broker.private_key}"
  sensitive = true
}

output "ops_manager_service_account_key_base64" {
  value     = "${google_service_account_key.ops_manager.private_key}"
  sensitive = true
}

output "service_broker_db_ip" {
  value = "${google_sql_database_instance.service_broker.ip_address.0.ip_address}"
}

output "service_broker_db_username" {
  value = "${random_id.service_broker_username.b64}"
}

output "service_broker_db_password" {
  value     = "${random_id.service_broker_password.b64}"
  sensitive = true
}

output "ops_manager_username" {
  value = "${var.ops_manager_username}"
}

output "ops_manager_password" {
  value     = "${var.ops_manager_password == "" ? random_id.ops_manager_password_generator.b64 : var.ops_manager_password}"
  sensitive = true
}

output "ops_manager_decryption_phrase" {
  value     = "${var.ops_manager_decryption_phrase == "" ? random_id.ops_manager_decryption_phrase_generator.b64 : var.ops_manager_decryption_phrase}"
  sensitive = true
}

output "ops_manager_skip_ssl_verify" {
  value = "${var.ops_manager_skip_ssl_verify}"
}

output "credhub_key_name" {
  value = "${var.credhub_key_name == "" ? random_id.credhub_key_name_generator.b64 : var.credhub_key_name}"
}

output "credhub_key" {
  value     = "${var.credhub_key == "" ? random_id.credhub_key_generator.b64 : var.credhub_key}"
  sensitive = true
}
