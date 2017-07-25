variable "project" {
  type = "string"
}

variable "env_name" {
  type    = "string"
  default = "omg"
}

variable "region" {
  type    = "string"
  default = "us-east1"
}

variable "zones" {
  type    = "list"
  default = ["us-east1-b", "us-east1-c", "us-east1-d"]
}

variable "opsman_image_url" {
  type        = "string"
  description = "location of ops manager image on google cloud storage"
  default     = "https://storage.cloud.google.com/ops-manager-us/pcf-gcp-1.11.4.tar.gz"
}

variable "opsman_image_selflink" {
  type        = "string"
  description = "location of ops manager image hosted by a google project"
  default     = ""
}

variable "opsman_machine_type" {
  type    = "string"
  default = "n1-standard-2"
}

variable "jumpbox_machine_type" {
  type    = "string"
  default = "n1-standard-2"
}

variable "nat_machine_type" {
  type    = "string"
  default = "n1-standard-1"
}

variable "nat_instance_count" {
  type    = "string"
  default = 3
}

variable "service_account_key" {
  type = "string"
}

variable "dns_suffix" {
  type = "string"
}

variable "ssl_cert" {
  type        = "string"
  description = "ssl certificate content for *.{env_name}.{dns_suffix}"
}

variable "ssl_cert_private_key" {
  type        = "string"
  description = "ssl certificate private key content *.{env_name}.{dns_suffix}"
}

variable "external_database" {
  description = "standups up a cloud sql database instance for the ops manager and ERT"
  default     = false
}

variable "no_ip_instance_tag" {
  description = "Instance tag for unnetworked instances and NAT routes"
  default     = "omg-no-ip"
}

variable "tcp_port_range" {
  description = "Port range for TCP router"
  default     = "1024-65535"
}

variable "ssh_public_key" {
  description = "Public SSH key to access Jumpbox/Ops Manager VMs"
}

/*******************
 * OpsMan Options  *
 *******************/

variable "ert_sql_db_host" {
  type    = "string"
  default = ""
}

variable "opsman_sql_db_host" {
  type    = "string"
  default = ""
}

/*****************************
 * Isolation Segment Options *
 *****************************/

variable "isolation_segment" {
  description = "create the required infrastructure to deploy isolation segment"
  default     = false
}

variable "iso_seg_ssl_cert" {
  type        = "string"
  description = "ssl certificate content"
  default     = ""
}

variable "iso_seg_ssl_cert_private_key" {
  type        = "string"
  description = "ssl certificate private key content"
  default     = ""
}
