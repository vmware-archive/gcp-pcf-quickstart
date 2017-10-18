resource "google_dns_record_set" "ops-manager-dns-internal" {
  count          = "${var.opsman_external_ip != "" ? 0 : 1}"

  name = "pcf.${var.dns_suffix}."
  type = "A"
  ttl  = 300

  managed_zone = "${var.dns_zone_name}"

  rrdatas = ["${google_compute_instance.ops-manager-internal.network_interface.0.address}"]
}

resource "google_dns_record_set" "ops-manager-dns-external" {
  count          = "${var.opsman_external_ip != "" ? 1 : 0}"

  name = "pcf.${var.dns_suffix}."
  type = "A"
  ttl  = 30

  managed_zone = "${var.dns_zone_name}"

  rrdatas = ["${google_compute_instance.ops-manager-external.network_interface.0.access_config.0.assigned_nat_ip}"]
}

resource "google_dns_record_set" "wildcard-sys-dns" {
  name = "*.sys.${var.dns_suffix}."
  type = "A"
  ttl  = 300

  managed_zone = "${var.dns_zone_name}"

  rrdatas = ["${google_compute_global_address.cf.address}"]
}

resource "google_dns_record_set" "wildcard-apps-dns" {
  name = "*.apps.${var.dns_suffix}."
  type = "A"
  ttl  = 300

  managed_zone = "${var.dns_zone_name}"

  rrdatas = ["${google_compute_global_address.cf.address}"]
}

resource "google_dns_record_set" "app-ssh-dns" {
  name = "ssh.sys.${var.dns_suffix}."
  type = "A"
  ttl  = 300

  managed_zone = "${var.dns_zone_name}"

  rrdatas = ["${google_compute_address.cf-ssh.address}"]
}

resource "google_dns_record_set" "tcp-dns" {
  name = "tcp.${var.dns_suffix}."
  type = "A"
  ttl  = 300

  managed_zone = "${var.dns_zone_name}"

  rrdatas = ["${google_compute_address.cf-tcp.address}"]
}

resource "google_dns_record_set" "jumpbox-dns" {
  name = "jumpbox.${var.dns_suffix}."
  type = "A"
  ttl  = 300

  managed_zone = "${var.dns_zone_name}"

  rrdatas = ["${google_compute_instance.jumpbox.network_interface.0.access_config.0.assigned_nat_ip}"]
}
