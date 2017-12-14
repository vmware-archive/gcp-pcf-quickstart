resource "random_id" "db-name" {
  byte_length = 8
}

resource "google_sql_database_instance" "service_broker" {
  region           = "${var.region}"
  database_version = "MYSQL_5_6"
  name             = "${var.env_name}-${replace(lower(random_id.db-name.b64), "_", "-")}"

  settings {
    tier = "${var.service_broker_sql_db_tier}"

    ip_configuration = {
      ipv4_enabled = true

      authorized_networks = [
        {
          name  = "all"
          value = "0.0.0.0/0"
        },
      ]
    }
  }
}

resource "google_sql_database" "uaa" {
  name       = "servicebroker"
  instance   = "${google_sql_database_instance.service_broker.name}"
  depends_on = ["google_sql_user.service_broker"]
}


resource "random_id" "service_broker_username" {
  byte_length = 8
}

resource "random_id" "service_broker_password" {
  byte_length = 32
}

resource "google_sql_user" "service_broker" {
  name     = "${random_id.service_broker_username.b64}"
  password = "${random_id.service_broker_password.b64}"
  instance = "${google_sql_database_instance.service_broker.name}"
  host     = "${var.ert_sql_db_host}"
}

resource "random_id" "service_broker_account" {
  byte_length = 4
}

resource "google_service_account" "service_broker" {
  display_name = "GCP Service Broker"
  account_id   = "sb-${random_id.service_broker_account.hex}"
}

resource "google_service_account_key" "service_broker" {
  service_account_id = "${google_service_account.service_broker.id}"
}

resource "google_project_iam_member" "service_broker" {
  project = "${var.project}"
  role    = "roles/owner"
  member  = "serviceAccount:${google_service_account.service_broker.email}"
}
