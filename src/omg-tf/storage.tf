resource "random_id" "suffix" {
  byte_length = 8
}

resource "google_storage_bucket" "buildpacks" {
  name          = "${var.env_name}-buildpacks-${random_id.suffix.b64}"
  force_destroy = true
}

resource "google_storage_bucket" "droplets" {
  name          = "${var.env_name}-droplets-${random_id.suffix.b64}"
  force_destroy = true
}

resource "google_storage_bucket" "packages" {
  name          = "${var.env_name}-packages-${random_id.suffix.b64}"
  force_destroy = true
}

resource "google_storage_bucket" "resources" {
  name          = "${var.env_name}-resources-${random_id.suffix.b64}"
  force_destroy = true
}
