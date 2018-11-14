output "sql_db_port" {
  value = "3306"
}

output "sql_db_ip" {
  value = "${ join(" ", google_sql_database_instance.master.*.ip_address.0.ip_address) }"
}

output "opsman_sql_db_name" {
  value = "${ join(" ", google_sql_database.opsman.*.name) }"
}

output "opsman_sql_username" {
  value = "${ join(" ", random_id.opsman_db_username.*.b64) }"
}

output "opsman_sql_password" {
  sensitive = true
  value = "${ join(" ", random_id.opsman_db_password.*.b64) }"
}

output "ert_sql_username" {
  value = "${ join(" ", random_id.ert_db_username.*.b64) }"
}

output "ert_sql_password" {
  sensitive = true
  value = "${ join(" ", random_id.ert_db_password.*.b64) }"
}

output "ip" {
  value = "${ join(" ", google_sql_database_instance.master.*.ip_address.0.ip_address) }"
}
