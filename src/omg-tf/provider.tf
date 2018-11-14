provider "google" {
  version     = "v1.19.1"

  project     = "${var.project}"
  region      = "${var.region}"
  credentials = "${var.service_account_key}"
}

provider "random" {
  version = "~> 1.0"
}
