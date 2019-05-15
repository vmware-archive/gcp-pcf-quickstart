provider "google" {
  version     = "~> 2.6.0"

  project     = "${var.project}"
  region      = "${var.region}"
  credentials = "${var.service_account_key}"
}

provider "random" {
  version = "~> 1.0"
}
