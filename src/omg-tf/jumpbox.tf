resource "google_compute_instance" "jumpbox" {
  name           = "${var.env_name}-jumpbox"
  machine_type   = "${var.jumpbox_machine_type}"
  zone           = "${element(var.zones, 1)}"
  create_timeout = 10
  tags           = ["${var.env_name}-jumpbox-external", "${var.instance_tag}"]

  disk {
    image = "projects/ubuntu-os-cloud/global/images/ubuntu-1404-trusty-v20170505"
    size  = 50
  }

  network_interface {
    subnetwork = "${google_compute_subnetwork.management-subnet.name}"
    access_config {
      # ephemeral IP
    }
  }
}

resource "google_compute_firewall" "jumpbox-external" {
  name        = "${var.env_name}-jumpbox-external"
  network     = "${google_compute_network.pcf-network.name}"
  target_tags = ["${var.env_name}-jumpbox-external"]

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }
}

