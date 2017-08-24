output "service_account_email" {
  value = "${google_service_account.opsman_service_account.email}"
}

output "ops_manager_dns" {
  value = "pcf.${var.dns_suffix}"
}

output "ops_manager_instance_name" {
  value = "${google_compute_instance.ops-manager.name}"
}

output "ops_manager_instance_zone" {
  value = "${google_compute_instance.ops-manager.zone}"
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

output "dns_suffix" {
  value = "${var.dns_suffix}"
}

output "ops_manager_private_ip" {
  value = "${google_compute_instance.ops-manager.network_interface.0.address}"
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

output "http_lb_backend_name" {
  value = "${google_compute_backend_service.http_lb_backend_service.name}"
}

output "isoseg_lb_backend_name" {
  value = "${module.isolation_segment.load_balancer_name}"
}

output "ssh_router_pool" {
  value = "${google_compute_target_pool.cf-ssh.name}"
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
  value = "${module.external_database.ert_sql_password}"
}

output "opsman_sql_username" {
  value = "${module.external_database.opsman_sql_username}"
}

output "opsman_sql_password" {
  value = "${module.external_database.opsman_sql_password}"
}

output "jumpbox_public_ip" {
  value = "${google_compute_instance.jumpbox.network_interface.0.access_config.0.assigned_nat_ip}"
}

output "ssl_cert" {
  value = "${var.ssl_cert}"
}

output "ssl_cert_private_key" {
  value = "${var.ssl_cert_private_key}"
}

output "tcp_port_range" {
  value = "${var.tcp_port_range}"
}

output "stackdriver_service_account_key" {
  value = "${var.stackdriver_service_account_key}"
}

output "service_broker_service_account_key" {
  value = "${var.service_broker_service_account_key}"
}

output "service_account_key" {
  value = "${var.service_account_key}"
}

output "service_broker_db_ip" {
  value = "${google_sql_database_instance.service_broker.ip_address.0.ip_address}"
}

output "service_broker_db_username" {
  value = "${random_id.service_broker_username.b64}"
}

output "service_broker_db_password" {
  value = "${random_id.service_broker_password.b64}"
}

output "pivnet_api_token" {
  value = "${var.pivnet_api_token}"
}

output "pivnet_accept_eula" {
  value = "${var.pivnet_accept_eula}"
}