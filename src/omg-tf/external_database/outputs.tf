output "sql_db_port" {
  value = "3306"
}

output "sql_db_ip" {
  value = "${google_sql_database_instance.master.0.ip_address.0.ip_address}"
}

output "opsman_sql_db_name" {
  value = "${google_sql_database.opsman.0.name}"
}

output "opsman_sql_username" {
  value = "${random_id.opsman_db_username.0.b64}"
}

output "opsman_sql_password" {
  sensitive = true
  value = "${random_id.opsman_db_password.0.b64}"
}

output "ert_sql_username" {
  value = "${random_id.ert_db_username.0.b64}"
}

output "ert_sql_password" {
  sensitive = true
  value = "${random_id.ert_db_password.0.b64}"
}

output "ip" {
  value = "${google_sql_database_instance.master.0.ip_address.0.ip_address}"
}