terraform {
  required_version = ">= 0.13"
}

resource "google_compute_subnetwork" "subnet" {
  name          = var.name
  project       = var.project
  region        = var.region
  network       = var.network
  ip_cidr_range = var.ip_range
}