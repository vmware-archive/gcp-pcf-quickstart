resource "random_id" "credhub_key_generator" {
  byte_length = 20
}

resource "random_id" "credhub_key_name_generator" {
  byte_length = 4
}
