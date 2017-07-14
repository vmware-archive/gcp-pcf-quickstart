resource "google_compute_image" "ops-manager-image" {
  name           = "${var.env_name}-ops-manager-image"
  create_timeout = 20

  raw_disk {
    source = "${var.opsman_image_url}"
  }
}

resource "google_compute_instance" "ops-manager" {
  name           = "${var.env_name}-ops-manager"
  machine_type   = "${var.opsman_machine_type}"
  zone           = "${element(var.zones, 1)}"
  create_timeout = 10
  tags           = ["${var.env_name}-ops-manager", "${var.no_ip_instance_tag}"]

  disk {
    image = "${google_compute_image.ops-manager-image.self_link}"
    size  = 50
  }

  network_interface {
    subnetwork = "${google_compute_subnetwork.management-subnet.name}"
    address    = "10.0.0.6"
  }

  service_account {
    email  = "${google_service_account.opsman_service_account.email}"
    scopes = ["cloud-platform"]
  }
}

resource "google_storage_bucket" "director" {
  name          = "${var.env_name}-director"
  force_destroy = true
}