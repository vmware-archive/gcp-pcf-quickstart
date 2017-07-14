resource "google_compute_instance" "nat" {
  name           = "${var.env_name}-nat-${count.index}"
  machine_type   = "${var.nat_machine_type}"
  zone           = "${element(var.zones, 1)}"
  create_timeout = 10
  tags           = ["${var.env_name}-nat-external"]
  count          = "${var.nat_instance_count}"

  disk {
    image = "debian-8-jessie-v20170619"
    size  = 10
  }

  network_interface {
    subnetwork = "${google_compute_subnetwork.management-subnet.name}"
    access_config {
      # ephemeral IP
    }
  }

  can_ip_forward = true

  metadata_startup_script = <<EOF
#!/bin/bash
sh -c "echo 1 > /proc/sys/net/ipv4/ip_forward"
iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
EOF
}

resource "google_compute_route" "nat-primary" {
  // Explicitly declare dependency on all of the .nat instances.
  // This is needed because we only explicitly reference the first nat instance.
  // We implicitly depend on the others because of the name interpolation.
  depends_on             = ["google_compute_instance.nat"]

  name                   = "${var.env_name}-nat-primary-${count.index}"
  dest_range             = "0.0.0.0/0"
  network                = "${google_compute_network.pcf-network.name}"
  next_hop_instance      = "${var.env_name}-nat-${count.index}"
  next_hop_instance_zone = "${google_compute_instance.nat.0.zone}"
  priority               = 800
  tags                   = ["${var.instance_tag}"]
  count                  = "${var.nat_instance_count}"
}