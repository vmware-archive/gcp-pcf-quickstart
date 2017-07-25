resource "google_compute_instance" "jumpbox" {
  name           = "${var.env_name}-jumpbox"
  machine_type   = "${var.jumpbox_machine_type}"
  zone           = "${element(var.zones, 1)}"
  create_timeout = 10
  tags           = ["${var.env_name}-jumpbox-external"]

  disk {
    image = "ubuntu-1404-trusty-v20170703"
    size  = 50
  }

  network_interface {
    subnetwork = "${google_compute_subnetwork.unmanaged-subnet.name}"
    access_config {
      # ephemeral IP
    }
  }

  metadata = {
    ssh-keys               = "${format("omg:%s", var.ssh_public_key)}"
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

  source_ranges = ["0.0.0.0/0"]
}