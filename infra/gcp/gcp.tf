locals {
  region              = "europe-west1" # Equivalent GCP region
  zone_a              = "europe-west1-b"
  machine_type        = "e2-micro" # Smallest machine type
  ssh_user            = "ubuntu"
}
provider "google" {
  project = "mdl-default-test-project"
  region  = local.region
}
# SSH Key
resource "tls_private_key" "ssh_key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}
resource "local_file" "private_key" {
  content  = tls_private_key.ssh_key.private_key_pem
  filename = "${path.module}/id_rsa_todos.pem"
}

# Firewall Rules (Equivalent to Security Groups)
resource "google_compute_firewall" "allow_web_ssh" {
  name    = "allow-web-ssh"
  network = google_compute_network.main_vpc.id

  allow {
    protocol = "tcp"
    ports    = ["22", "8000"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["web-server"]
}

resource "google_compute_firewall" "allow_postgres" {
  name    = "allow-postgres"
  network = google_compute_network.main_vpc.id

  allow {
    protocol = "tcp"
    ports    = ["5432"]
  }

  source_ranges = ["0.0.0.0/0"] # Replace this with Bastion IP for security
  target_tags   = ["data-server"]
}

# Web Server Instances (Public Subnet)
resource "google_compute_instance" "web_vm_1" {
  name         = "web-vm-1"
  machine_type = local.machine_type
  zone         = local.zone_a

  tags = ["web-server"]

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-2204-lts" # Ubuntu Image
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.public_subnet.id
    access_config {}
  }

  metadata_startup_script = file("init-webapp-vm.yaml")

  metadata = {
    ssh-keys = "${local.ssh_user}:${tls_private_key.ssh_key.public_key_openssh}"
  }
}

resource "google_compute_instance" "data_vm_1" {
  name         = "data-vm-1"
  machine_type = local.machine_type
  zone         = local.zone_a

  tags = ["data-server"]

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-2204-lts" # Ubuntu Image
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.private_subnet.id
  }

  metadata_startup_script = file("init-data-vm.yaml")

  metadata = {
    ssh-keys = "${local.ssh_user}:${tls_private_key.ssh_key.public_key_openssh}"
  }
}

