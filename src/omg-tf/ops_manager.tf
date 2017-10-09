variable "ops_manager_image_name" {
  default = "ops-manager-image"
}

resource "google_compute_image" "ops-manager-image" {
  count = "${var.opsman_image_selflink != "" ? 0 : 1}"
  name           = "${var.ops_manager_image_name}"
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

  boot_disk {
    image = "${var.opsman_image_selflink != "" ? var.opsman_image_selflink : var.ops_manager_image_name}"
    size  = 250
    type  = "pd-ssd"
  }

  network_interface {
    subnetwork = "${google_compute_subnetwork.management-subnet.name}"
    address    = "10.0.0.6"
  }

  service_account {
    email  = "${google_service_account.opsman_service_account.email}"
    scopes = ["cloud-platform"]
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


