resource "random_id" "stackdriver_nozzle_account" {
  byte_length = 4
}

resource "google_service_account" "stackdrvier_nozzle" {
  display_name = "Stackdriver Nozzle"
  account_id   = "sdn-${random_id.stackdriver_nozzle_account.hex}"
}

resource "google_service_account_key" "stackdrvier_nozzle" {
  service_account_id = "${google_service_account.stackdrvier_nozzle.id}"
}

resource "google_project_iam_member" "stackdriver_nozzle" {
  project = "${var.project}"
  role    = "roles/editor"
  member  = "serviceAccount:${google_service_account.stackdrvier_nozzle.email}",
}
