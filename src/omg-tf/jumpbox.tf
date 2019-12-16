resource "google_compute_instance" "jumpbox" {
  name           = "${var.env_name}-jumpbox"
  machine_type   = "${var.jumpbox_machine_type}"
  zone           = "${element(var.zones, 1)}"

  timeouts {
    create = "10m"
  }

  tags           = ["${var.env_name}-jumpbox-external"]

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-1804-lts"
      size  = 250
      type  = "pd-ssd"
    }
  }

  network_interface {
    subnetwork = "${google_compute_subnetwork.unmanaged-subnet.name}"
    access_config {
      nat_ip = "${google_compute_address.jumpbox.address}"
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

resource "google_compute_address" "jumpbox" {
  name = "${var.env_name}-jumpbox"
}
