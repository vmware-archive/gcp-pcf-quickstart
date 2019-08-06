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

resource "google_storage_bucket" "backup" {
  name          = "${var.env_name}-backup-${random_id.suffix.hex}"
  force_destroy = true
}

resource "google_service_account" "blobstore" {
  count = "${var.create_blobstore_service_account_key}"

  account_id   = "${var.env_name}-blobstore"
  display_name = "${var.env_name} Blobstore Service Account"
}

resource "google_service_account_key" "blobstore" {
  count = "${var.create_blobstore_service_account_key}"

  service_account_id = "${google_service_account.blobstore.id}"
}

resource "google_project_iam_member" "blobstore_cloud_storage_admin" {
  count = "${var.create_blobstore_service_account_key}"

  project = "${var.project}"
  role    = "roles/storage.admin"
  member  = "serviceAccount:${google_service_account.blobstore.email}"
}
