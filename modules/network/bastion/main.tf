terraform {
  required_version = ">= 0.13"
}


resource "google_compute_instance" "bastion" {
  name         = var.name
  project      = var.project
  machine_type = var.instance_type
  zone         = var.zones[0]
  deletion_protection = false


  metadata = {
    "ssh-keys" = "${var.user}:${file(var.ssh_key)}"
  }

  boot_disk {
    initialize_params {
      image = var.image
    }
  }

  network_interface {
    subnetwork = var.subnet_name

    access_config {
      # Ephemeral IP - leaving this block empty will generate a new external IP and assign it to the machine
    }
  }

  tags = ["bastion"]
 labels = {
    environmentname = var.name
    owner = var.name
    
  }
}
