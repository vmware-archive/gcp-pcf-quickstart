# Allow HTTP/S access to Ops Manager from the outside world if exposed
resource "google_compute_firewall" "ops-manager-external" {
  name        = "${var.env_name}-ops-manager-external"
  network     = "${google_compute_network.pcf-network.name}"
  target_tags = ["${var.env_name}-ops-manager-external"]

  allow {
    protocol = "tcp"
    ports    = ["443", "80"]
  }

  source_ranges = ["0.0.0.0/0"]
}

resource "google_compute_image" "ops-manager-image" {
  count          = "${var.opsman_image_selflink != "" ? 0 : 1}"
  name           = "${var.env_name}-ops-manager-image"
  create_timeout = 20

  raw_disk {
    source = "${var.opsman_image_url}"
  }
}

resource "google_compute_instance" "ops-manager-internal" {
  count = "${var.opsman_external_ip != "" ? 0 : 1}"

  name           = "${var.env_name}-ops-manager"
  machine_type   = "${var.opsman_machine_type}"
  zone           = "${element(var.zones, 1)}"
  create_timeout = 10
  tags           = ["${var.env_name}-ops-manager", "${var.no_ip_instance_tag}"]

  boot_disk {
    initialize_params {
      image = "${var.opsman_image_selflink != "" ? var.opsman_image_selflink : google_compute_image.ops-manager-image.self_link}"
      size  = 250
      type  = "pd-ssd"
    }
  }

  network_interface {
    subnetwork = "${google_compute_subnetwork.management-subnet.name}"
    address    = "10.0.0.6"
  }

  metadata = {
    ssh-keys               = "${format("ubuntu:%s", var.ssh_public_key)}"
    block-project-ssh-keys = "TRUE"
  }
}

resource "google_compute_address" "ops-manager-external" {
  count = "${var.opsman_external_ip != "" ? 1 : 0}"
  name  = "${var.env_name}-ops-manager"
}

resource "google_compute_instance" "ops-manager-external" {
  count = "${var.opsman_external_ip != "" ? 1 : 0}"

  name           = "${var.env_name}-ops-manager"
  machine_type   = "${var.opsman_machine_type}"
  zone           = "${element(var.zones, 1)}"
  create_timeout = 10
  tags           = ["${var.env_name}-ops-manager", "${var.env_name}-ops-manager-external"]

  boot_disk {
    initialize_params {
      image = "${var.opsman_image_selflink != "" ? var.opsman_image_selflink : google_compute_image.ops-manager-image.self_link}"
      size  = 250
      type  = "pd-ssd"
    }
  }

  network_interface {
    subnetwork = "${google_compute_subnetwork.management-subnet.name}"
    address    = "10.0.0.6"

    access_config {
      nat_ip = "${google_compute_address.ops-manager-external.address}"
    }
  }

  metadata = {
    ssh-keys               = "${format("ubuntu:%s", var.ssh_public_key)}"
    block-project-ssh-keys = "TRUE"
  }
}

resource "random_id" "ops_manager_password_generator" {
  byte_length = 16
}

resource "random_id" "ops_manager_decryption_phrase_generator" {
  byte_length = 16
}

resource "random_id" "ops_manager_account" {
  byte_length = 4
}

resource "google_service_account" "ops_manager" {
  display_name = "Ops Manager"
  account_id   = "ops-${random_id.ops_manager_account.hex}"
}

resource "google_service_account_key" "ops_manager" {
  service_account_id = "${google_service_account.ops_manager.id}"
}

resource "google_project_iam_member" "ops_manager" {
  project = "${var.project}"
  role    = "roles/owner"
  member  = "serviceAccount:${google_service_account.ops_manager.email}"
}
