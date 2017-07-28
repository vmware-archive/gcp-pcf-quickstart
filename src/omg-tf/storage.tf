resource "random_id" "suffix" {
  byte_length = 4
}

resource "google_storage_bucket" "director" {
  name          = "${var.env_name}-director-${random_id.suffix.hex}"
  force_destroy = true
}

resource "google_storage_bucket" "buildpacks" {
  name          = "${var.env_name}-buildpacks-${random_id.suffix.hex}"
  force_destroy = true
}

resource "google_storage_bucket" "droplets" {
  name          = "${var.env_name}-droplets-${random_id.suffix.hex}"
  force_destroy = true
}

resource "google_storage_bucket" "packages" {
  name          = "${var.env_name}-packages-${random_id.suffix.hex}"
  force_destroy = true
}

resource "google_storage_bucket" "resources" {
  name          = "${var.env_name}-resources-${random_id.suffix.hex}"
  force_destroy = true
}
