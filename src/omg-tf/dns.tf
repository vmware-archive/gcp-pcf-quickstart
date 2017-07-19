resource "google_dns_managed_zone" "env_dns_zone" {
  name        = "${var.env_name}-zone"
  dns_name    = "${var.env_name}.${var.dns_suffix}."
  description = "DNS zone for the ${var.env_name} environment"
}

resource "google_dns_record_set" "ops-manager-dns" {
  name = "pcf.${google_dns_managed_zone.env_dns_zone.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = "${google_dns_managed_zone.env_dns_zone.name}"

  rrdatas = ["${google_compute_instance.ops-manager.network_interface.0.address}"]
}

resource "google_dns_record_set" "wildcard-sys-dns" {
  name = "*.sys.${google_dns_managed_zone.env_dns_zone.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = "${google_dns_managed_zone.env_dns_zone.name}"

  rrdatas = ["${google_compute_global_address.cf.address}"]
}

resource "google_dns_record_set" "wildcard-apps-dns" {
  name = "*.apps.${google_dns_managed_zone.env_dns_zone.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = "${google_dns_managed_zone.env_dns_zone.name}"

  rrdatas = ["${google_compute_global_address.cf.address}"]
}

resource "google_dns_record_set" "app-ssh-dns" {
  name = "ssh.sys.${google_dns_managed_zone.env_dns_zone.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = "${google_dns_managed_zone.env_dns_zone.name}"

  rrdatas = ["${google_compute_address.cf-ssh.address}"]
}

resource "google_dns_record_set" "tcp-dns" {
  name = "tcp.${google_dns_managed_zone.env_dns_zone.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = "${google_dns_managed_zone.env_dns_zone.name}"

  rrdatas = ["${google_compute_address.cf-tcp.address}"]
}

resource "google_dns_record_set" "jumpbox-dns" {
  name = "jumpbox.${google_dns_managed_zone.env_dns_zone.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = "${google_dns_managed_zone.env_dns_zone.name}"

  rrdatas = ["${google_compute_instance.jumpbox.network_interface.0.access_config.0.assigned_nat_ip}"]
}
