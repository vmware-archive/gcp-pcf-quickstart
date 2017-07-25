resource "google_dns_record_set" "ops-manager-dns" {
  name = "pcf.${var.dns_suffix}."
  type = "A"
  ttl  = 300

  managed_zone = "${var.dns_zone_name}"

  rrdatas = ["${google_compute_instance.ops-manager.network_interface.0.address}"]
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
