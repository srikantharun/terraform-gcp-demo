terraform {
  required_version = ">= 0.13"
}

resource "google_sql_database_instance" "master" {
  name             = var.db_name
  project          = var.project
  region           = var.region
  database_version = var.database_version
  deletion_protection = false

  settings {
    tier                 = var.tier
    activation_policy    = var.activation_policy
    disk_autoresize      = var.disk_autoresize
    disk_size            = var.disk_size
    disk_type            = var.disk_type
    pricing_plan         = var.pricing_plan

    ip_configuration {
      ipv4_enabled = true
      authorized_networks {
        value = "0.0.0.0/0"
        name  = "all"
      }
    }
  }
}

resource "google_sql_database" "default" {
  name      = var.db_name
  project   = var.project
  instance  = google_sql_database_instance.master.name
  charset   = var.db_charset
  collation = var.db_collation
}

resource "google_sql_user" "default" {
  name     = var.user_name
  project  = var.project
  instance = google_sql_database_instance.master.name
  host     = var.user_host
  password = var.user_password
}
